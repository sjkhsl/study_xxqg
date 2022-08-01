// Package utils
// @Description:
package utils

import (
	"crypto/md5"
	"encoding/hex"
	"os"
)

// FileIsExist
/* @Description:
 * @param path
 * @return bool
 */
func FileIsExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// StrMd5
/* @Description:
 * @param str
 * @return retMd5
 */
func StrMd5(str string) (retMd5 string) {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
