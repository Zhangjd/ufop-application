package utils

import (
	"encoding/base64"
	"regexp"
	"strings"
)

func GetParam(fromStr, pattern, key string) (value string) {
	keyRegx := regexp.MustCompile(pattern)
	matchStr := keyRegx.FindString(fromStr)
	value = strings.Replace(matchStr, key+"/", "", -1)
	return
}

func GetParamDecoded(fromStr, pattern, key string) (value string, err error) {
	strToDecode := GetParam(fromStr, pattern, key)
	decodedBytes, decodeErr := base64.URLEncoding.DecodeString(strToDecode)
	if decodeErr != nil {
		err = decodeErr
		return
	}
	value = string(decodedBytes)
	return
}
