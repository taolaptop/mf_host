package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s)) // 需要加密的字符串为
	cipherStr := h.Sum(nil)
	fmt.Println(cipherStr)
	fmt.Printf("%s\n", hex.EncodeToString(cipherStr)) // 输出加密结果
	return hex.EncodeToString(cipherStr)
}
