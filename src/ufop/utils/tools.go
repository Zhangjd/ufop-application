package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func Md5Hex(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Download(remoteUrl, localPath string) (contentType string, err error) {
	resp, respErr := http.Get(remoteUrl)
	if respErr != nil || resp.StatusCode != http.StatusOK {
		if respErr != nil {
			err = errors.New(fmt.Sprintf("get resource by url '%s' failed, %s", remoteUrl, respErr.Error()))
		} else {
			err = errors.New(fmt.Sprintf("get resource by url '%s' faild, %s", remoteUrl, resp.Status))
			if resp.Body != nil {
				resp.Body.Close()
			}
		}
		return
	}

	defer resp.Body.Close()

	contentType = resp.Header.Get("Content-Type")
	localFp, openErr := os.Create(localPath)
	if openErr != nil {
		err = errors.New(fmt.Sprintf("open file by local path failed, %s", openErr.Error()))
		return
	}

	defer localFp.Close()

	_, cpErr := io.Copy(localFp, resp.Body)

	if cpErr != nil {
		err = errors.New(fmt.Sprintf("save remote file to local failed, %s", cpErr.Error()))
		return
	}

	return
}

func MaxInt(array ...int) int {
	max := array[0]
	for _, val := range array {
		if val > max {
			max = val
		}
	}
	return max
}

func MinInt(array ...int) int {
	min := array[0]
	for _, val := range array {
		if val <= min {
			min = val
		}
	}
	return min
}
