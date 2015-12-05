/**
 * Author: Zhangjd
 * Date: December 5th, 2015
 * Reference: http://developer.qiniu.com/docs/v6/api/reference/fop/pfop/pfop.html
 * Description: 模拟调用七牛触发持久化处理（pfop）接口
 */

package videomerge

import (
    "encoding/json"
    "errors"
    "fmt"
    "github.com/qiniu/api.v6/auth/digest"
    "github.com/qiniu/log"
    "os"
    "ufop"
)

const (
    AUDIO_MERGE_MAX_FIRST_FILE_LENGTH  = 100 * 1024 * 1024
    AUDIO_MERGE_MAX_SECOND_FILE_LENGTH = 100 * 1024 * 1024
)

type ReqArgs struct {
    Cmd string `json:"cmd"`
    Src struct {
        Url      string `json:"url"`
        Mimetype string `json:"mimetype"`
        Fsize    int32  `json:"fsize"`
        Bucket   string `json:"bucket"`
        Key      string `json:"key"`
    } `json: "src"`
}

type VideoMerger struct {
    mac                 *digest.Mac
    maxFirstFileLength  uint64
    maxSecondFileLength uint64
}

type VideoMergerConfig struct {
    //ak & sk
    AccessKey string `json:"access_key"`
    SecretKey string `json:"secret_key"`

    AmergeMaxFirstFileLength  uint64 `json:"amerge_max_first_file_length,omitempty"`
    AmergeMaxSecondFileLength uint64 `json:"amerge_max_second_file_length,omitempty"`
}

func (this *VideoMerger) Name() string {
    return "videomerge"
}

func (this *VideoMerger) InitConfig(jobConf string) (err error) {
    confFp, openErr := os.Open(jobConf)
    if openErr != nil {
        err = errors.New(fmt.Sprintf("Open amerge config failed, %s", openErr.Error()))
        return
    }

    config := VideoMergerConfig{}
    decoder := json.NewDecoder(confFp)
    decodeErr := decoder.Decode(&config)
    if decodeErr != nil {
        err = errors.New(fmt.Sprintf("Parse amerge config failed, %s", decodeErr.Error()))
        return
    }

    if config.AmergeMaxFirstFileLength <= 0 {
        this.maxFirstFileLength = AUDIO_MERGE_MAX_FIRST_FILE_LENGTH
    } else {
        this.maxFirstFileLength = config.AmergeMaxFirstFileLength
    }

    if config.AmergeMaxSecondFileLength <= 0 {
        this.maxSecondFileLength = AUDIO_MERGE_MAX_SECOND_FILE_LENGTH
    } else {
        this.maxSecondFileLength = config.AmergeMaxSecondFileLength
    }

    this.mac = &digest.Mac{config.AccessKey, []byte(config.SecretKey)}

    return
}


func (this *VideoMerger) parse(cmd string) (format string, mime string, bucket string, url string, duration string, err error) {
    return
}


/**
 * [Do 此ufop指令对应的处理函数入口]
 * @param  {ufop.UfopRequest} req []
 * @return {interface{}}      result [返回结果]
 * @return {int}              resultType [返回结果类型]
 * @return {string}           contentType [http头部的Content-Type值]
 * @return {error}            err [错误对象]
 */
func (this *VideoMerger) Do(req ufop.UfopRequest) (result interface{}, resultType int, contentType string, err error) {
    // http.HandleFunc("/uop", fooHandler)
    // httpErr := http.ListenAndServe(":9100", nil)
    // if httpErr != nil {
    //     log.Fatal("Demo server failed to start:", err)
    // }

    log.Info("haha1")
    return
}




