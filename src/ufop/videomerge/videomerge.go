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
    "github.com/qiniu/api.v6/rs"
    "github.com/qiniu/log"
    "io"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"
    "os/exec"
    "regexp"
    "strings"
    "ufop"
    "ufop/utils"
)

const (
    AUDIO_MERGE_MAX_FIRST_FILE_LENGTH  = 100 * 1024 * 1024
    AUDIO_MERGE_MAX_SECOND_FILE_LENGTH = 100 * 1024 * 1024
)

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

/**
 * [InitConfig 读取配置文件信息(只在程序启动时执行一次)]
 */
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

/**
 * [parse 解析指令内容]
 * @param  {string} cmd [命令字符串(已经过trim prefix)]
 * @return {string} format  [输出格式]
 * @return {string} mime    [输出Mime-Type]
 * @return {string} bucket  [输出格式]
 * @return {string} url     [第二个视频的Url]
 * @return {string} duration  [输出格式]
 * @return {error}  err  [输出格式]
 */
func (this *VideoMerger) parse(cmd string) (format string, mime string, bucket string, url string, duration string, err error) {
    // 正则匹配
    // pattern := "^videomerge/format/[a-zA-Z0-9]+/mime/[0-9a-zA-Z-_=]+/bucket/[0-9a-zA-Z-_=]+/url/[0-9a-zA-Z-_=]+(/duration/(first|shortest|longest)){0,1}$"
    pattern := "^videomerge/format/[a-zA-Z0-9]+/mime/[0-9a-zA-Z-_=]+/bucket/[0-9a-zA-Z-_=]+/url/[0-9a-zA-Z-_=]+$"
    matched, _ := regexp.MatchString(pattern, cmd)
    if !matched {
        err = errors.New("invalid videomerge command format")
        return
    }
    var decodeErr error
    // 获取格式
    format = utils.GetParam(cmd, "format/[a-zA-Z0-9]+", "format")
    // 获取Mime-Type
    mime, decodeErr = utils.GetParamDecoded(cmd, "mime/[0-9a-zA-Z-_=]+", "mime")
    if decodeErr != nil {
        err = errors.New("invalid amerge parameter 'mime'")
        return
    }
    // 获取bucket
    bucket, decodeErr = utils.GetParamDecoded(cmd, "bucket/[0-9a-zA-Z-_=]+", "bucket")
    if decodeErr != nil {
        err = errors.New("invalid amerge parameter 'bucket'")
        return
    }
    // 获取文件url
    url, decodeErr = utils.GetParamDecoded(cmd, "url/[0-9a-zA-Z-_=]+", "url")
    if decodeErr != nil {
        err = errors.New("invalid amerge parameter 'url'")
        return
    }
    // 获取duration
    duration = utils.GetParam(cmd, "duration/(first|shortest|longest)", "duration")
    if duration == "" {
        duration = "longest"
    }
    return
}


/**
 * [此ufop指令对应的处理函数入口]
 * @param  {ufop.UfopRequest} req         []
 * @return {interface{}}      result      [返回结果]
 * @return {int}              resultType  [返回结果类型]
 * @return {string}           contentType [http头部的Content-Type值]
 * @return {error}            err         [错误对象]
 */
func (this *VideoMerger) Do(req ufop.UfopRequest) (result interface{}, resultType int, contentType string, err error) {
    // parse command
    dstFormat, dstMime, secondFileBucket, secondFileUrl, dstDuration, pErr := this.parse(req.Cmd)
    if pErr != nil {
        log.Error(pErr)
        err = pErr
        return
    }
    log.Info(dstFormat, dstMime, secondFileBucket, secondFileUrl, dstDuration)

    // check first file
    if req.Src.Fsize > this.maxFirstFileLength {
        err = errors.New("first file length exceeds the limit")
        return
    }
    if !strings.HasPrefix(req.Src.MimeType, "video/") {
        err = errors.New("first file mime-type not supported")
        return
    }

    // check second file
    secondFileUri, pErr := url.Parse(secondFileUrl)
    if pErr != nil {
        err = errors.New("second file resource url not valid")
        return
    }
    secondFileKey := strings.TrimPrefix(secondFileUri.Path, "/")
    client := rs.New(this.mac)
    sEntry, sErr := client.Stat(nil, secondFileBucket, secondFileKey)
    if sErr != nil || sEntry.Hash == "" {
        err = errors.New("second file not in the specified bucket")
        return
    }
    if uint64(sEntry.Fsize) > this.maxSecondFileLength {
        err = errors.New("second file length exceeds the limit")
        return
    }
    if !strings.HasPrefix(sEntry.MimeType, "video/") {
        err = errors.New("second file mimetype not supported")
        return
    }

    // retrieve the first file
    fResp, fRespErr := http.Get(req.Src.Url)
    if fRespErr != nil || fResp.StatusCode != 200 {
        if fRespErr != nil {
            err = errors.New(fmt.Sprintf("retrieve first file resource data failed, %s", fRespErr.Error()))
        } else {
            err = errors.New(fmt.Sprintf("retrieve first file resource data failed, %s", fResp.Status))
            if fResp.Body != nil {
                fResp.Body.Close()
            }
        }
        return
    }
    fTmpFp, fErr := ioutil.TempFile("", "first")
    if fErr != nil {
        err = errors.New(fmt.Sprintf("open first file temp file failed, %s", fErr.Error()))
        return
    }
    _, fCpErr := io.Copy(fTmpFp, fResp.Body)
    if fCpErr != nil {
        err = errors.New(fmt.Sprintf("save first temp file failed, %s", fCpErr.Error()))
        return
    }

    // close the first one
    fTmpFname := fTmpFp.Name()
    fTmpFp.Close()
    fResp.Body.Close()

    // retrieve the second file
    sResp, sRespErr := http.Get(secondFileUrl)
    if sRespErr != nil || sResp.StatusCode != 200 {
        if sRespErr != nil {
            err = errors.New(fmt.Sprintf("retrieve second file resource data failed, %s", sRespErr.Error()))
        } else {
            err = errors.New(fmt.Sprintf("retrieve second file resource data failed, %s", sResp.Status))
            if sResp.Body != nil {
                sResp.Body.Close()
            }
        }
        return
    }
    sTmpFp, sErr := ioutil.TempFile("", "second")
    if sErr != nil {
        err = errors.New(fmt.Sprintf("open second file temp file failed, %s", sErr.Error()))
        return
    }
    _, sCpErr := io.Copy(sTmpFp, sResp.Body)
    if sCpErr != nil {
        err = errors.New(fmt.Sprintf("save second first tmp file failed, %s", sCpErr.Error()))
        return
    }
    
    // close the second one
    sTmpFname := sTmpFp.Name()
    sTmpFp.Close()
    sResp.Body.Close()

    // do conversion
    oTmpFp, oErr := ioutil.TempFile("", "output")
    if oErr != nil {
        err = errors.New(fmt.Sprintf("open output file temp file failed, %s", oErr.Error()))
        return
    }
    oTmpFname := oTmpFp.Name()
    oTmpFp.Close()

    // be sure to delete temp files
    defer os.Remove(fTmpFname)
    defer os.Remove(sTmpFname)

    // prepare for ffmpeg
    mergeCmdParams := []string{
        "-y",
        "-v", "error",
        "-i", fTmpFname,
        "-i", sTmpFname,
        // "-filter_complex", fmt.Sprintf("amix=inputs=2:duration=%s:dropout_transition=2", dstDuration),
        "-filter_complex", fmt.Sprintf("[0:v:0]pad=iw*2:ih[bg]; [bg][1:v:0]overlay=w"),
        "-f", dstFormat,
        oTmpFname,
    }

    // execute command
    mergeCmd := exec.Command("ffmpeg", mergeCmdParams...)
    stdErrPipe, pipeErr := mergeCmd.StderrPipe()
    if pipeErr != nil {
        err = errors.New(fmt.Sprintf("open exec stderr pipe error, %s", pipeErr.Error()))
        return
    }
    if startErr := mergeCmd.Start(); startErr != nil {
        err = errors.New(fmt.Sprintf("start ffmpeg command error, %s", startErr.Error()))
        return
    }
    stdErrData, readErr := ioutil.ReadAll(stdErrPipe)
    if readErr != nil {
        err = errors.New(fmt.Sprintf("read ffmpeg command stderr error, %s", readErr.Error()))
        defer os.Remove(oTmpFname)
        return
    }

    // check stderr output & output file
    if string(stdErrData) != "" {
        log.Error(string(stdErrData))
    }
    if waitErr := mergeCmd.Wait(); waitErr != nil {
        err = errors.New(fmt.Sprintf("wait ffmpeg to exit error, %s", waitErr))
        defer os.Remove(oTmpFname)
        return
    }
    if oFileInfo, statErr := os.Stat(oTmpFname); statErr != nil || oFileInfo.Size() == 0 {
        err = errors.New("audio merge with no valid output result")
        defer os.Remove(oTmpFname)
        return
    }

    // write result
    result = oTmpFname
    resultType = ufop.RESULT_TYPE_OCTECT
    contentType = dstMime

    log.Info(result)

    return
}




