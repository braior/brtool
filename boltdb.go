package brtool

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/boltdb/bolt"
)

var (
	db *bolt.DB
)

// BoltDB DB 操作
type BoltDB struct {
	DBPath    string // 数据库路径
	TableName string // 表名
}

// NewBoltDB 初始化数据库对象
func NewBoltDB(dbPath, tableName string) (*BoltDB, error) {
	if dbPath == "" {
		return nil, fmt.Errorf("database path required")
	}

	dirName := path.Dir(dbPath)
	if !IsExist(dirName){
		err := os.MkdirAll(dirName, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("create dir(%s) error: %s", dirName, err)
		}
	}
	return &BoltDB{DBPath: dbPath, TableName: tableName}, nil
}

// talbe 获取表，表不存在则创建
func (btb *BoltDB) table() error {
	var err error
	db, err = bolt.Open(btb.DBPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return fmt.Errorf("open db error: %s", err)
	}

	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(btb.TableName))
		if err != nil {
			return fmt.Errorf("create table error: %s", err)
		}
		return nil
	})
}

// Set 设置键值对
func (btb *BoltDB) Set(kv map[string][]byte) error {
	err := btb.table()
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(btb.TableName))
		var err error
		for k, v := range kv {
			err = b.Put([]byte(k), v)
			if err != nil {
				return err
			}
		}
		return err
	})
}

// Get 根据键名数组获取各自的值
// keys 键名数组
func (btb *BoltDB) Get(keys []string) (map[string][]byte, error) {
	err := btb.table()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	values := make(map[string][]byte)
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(btb.TableName))
		for _, k := range keys {
			result := b.Get([]byte(k))
			if result == nil {
				continue
			}
			tmp := make([]byte, len(result))
			copy(tmp, result)
			values[k] = tmp
		}
		return nil
	})
	return values, err
}

// GetAll 获取全部键值
func (btb *BoltDB) GetAll() (map[string][]byte, error) {
	err := btb.table()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	values := make(map[string][]byte)
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(btb.TableName))
		b.ForEach(func(k, v []byte) error {
			tmpV := make([]byte, len(v))
			copy(tmpV, v)

			tmpK := make([]byte, len(k))
			copy(tmpK, k)
			values[string(tmpK)] = tmpV
			return nil
		})
		return nil
	})
	return values, nil
}

// Delete 删除键值
func (btb *BoltDB) Delete(keys []string) error {
	err := btb.table()
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(btb.TableName))
		var err error
		for _, k := range keys {
			err = b.Delete([]byte(k))
			if err != nil {
				return err
			}
		}
		return err
	})
}

// Backup 备份数据库文件
func (btb *BoltDB) Backup(filePath string) error {
	db, err := bolt.Open(btb.DBPath, 0600, nil)
	if err != nil {
		return fmt.Errorf("open db error: %s", err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		err := tx.CopyFile(filePath, 0644)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return db.Close()
}
