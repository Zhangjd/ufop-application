package utils

import (
	"golang.org/x/text/encoding/simplifiedchinese"
)

var gbkDecoder = simplifiedchinese.GBK.NewDecoder()
var gbkEncoder = simplifiedchinese.GBK.NewEncoder()

func Gbk2Utf8(text string) (string, error) {
	utf8Dst := make([]byte, len(text)*3)
	_, _, err := gbkDecoder.Transform(utf8Dst, []byte(text), true)
	if err != nil {
		return "", err
	}
	gbkDecoder.Reset()
	utf8Bytes := make([]byte, 0)
	for _, b := range utf8Dst {
		if b != 0 {
			utf8Bytes = append(utf8Bytes, b)
		}
	}
	return string(utf8Bytes), nil
}

func Utf82Gbk(text string) (string, error) {
	gbkDst := make([]byte, len(text)*2)
	_, _, err := gbkEncoder.Transform(gbkDst, []byte(text), true)
	if err != nil {
		return "", err
	}
	gbkEncoder.Reset()
	gbkBytes := make([]byte, 0)
	for _, b := range gbkDst {
		if b != 0 {
			gbkBytes = append(gbkBytes, b)
		}
	}
	return string(gbkBytes), nil
}
