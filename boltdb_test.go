package brtool

import (
	"testing"
)

func TestNewBoltDB(t *testing.T) {
	if ans, _ := NewBoltDB("test.db", "token"); ans.TableName != "token" {
		t.Errorf("return mytoken, but %s got", ans.TableName)
	}

}

func Test_BoltDB_Set(t *testing.T) {

	ans, err := NewBoltDB("test.db", "token")
	if ans.TableName != "token" {
		t.Errorf("new bolt db false, got err: %s", err)
	}

	kv := make(map[string][]byte)
	kv["root"] = []byte("123")

	err = ans.Set(kv)
	if err != nil {
		t.Errorf("set key false, got err: %s", err)
	}

	result, _ := ans.Get([]string{"root"})
	if string(result["root"]) != "123" {
		t.Errorf("get key will be 123, but got %s", string(result["root"]))
	}
}

// func TestNewJWToken(t *testing.T) {
// 	if ans := NewJWToken("myToken"); ans.SignString != "yToken" {
// 		t.Errorf("return mytoken, but %s got", ans.SignString)
// 	}
// }
