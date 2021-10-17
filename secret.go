package brtool

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

// GetMD5 生成32位MD5
func GetMD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

// GenRandomString 生成随机字符串
// length 生成长度
// specialChar 是否生成特殊字符
func GenRandomString(length int, charSet string) string {
	rand.Seed(time.Now().UnixNano())
	numStr := "0123456789"
	charStr := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specStr := "+=-@#~,.[]()!%^*$"

	var passwd []byte = make([]byte, length)
	var sourceStr string

	switch charSet {
	case "num":
		sourceStr = numStr
	case "char":
		sourceStr = charStr
	case "mix":
		sourceStr = fmt.Sprintf("%s%s", numStr, charStr)
	case "advance":
		sourceStr = fmt.Sprintf("%s%s%s", numStr, charStr, specStr)
	default:
		sourceStr = fmt.Sprintf("%s%s", numStr, charStr)
	}

	for i := 0; i < length; i++ {
		index := rand.Intn(len(sourceStr))
		passwd[i] = sourceStr[index]
	}

	return string(passwd)
}
