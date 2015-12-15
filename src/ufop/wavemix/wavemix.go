/**
 * Author: Zhangjd
 * Date: December 8th, 2015
 * Reference: http://developer.qiniu.com/docs/v6/api/reference/fop/pfop/pfop.html
 * Description: 声波合成模块
 */

package wavemix

import (
    "errors"
    "fmt"
    "io/ioutil"
    "github.com/qiniu/api.v6/auth/digest"
    "github.com/qiniu/log"
    "os"
    "os/exec"
    "regexp"
    "strings"
    "strconv"
    "ufop"
)

const (
    AUDIO_MERGE_MAX_FIRST_FILE_LENGTH  = 100 * 1024 * 1024
    AUDIO_MERGE_MAX_SECOND_FILE_LENGTH = 100 * 1024 * 1024
)

type WaveMixer struct {
    mac                 *digest.Mac
    maxFirstFileLength  uint64
    maxSecondFileLength uint64
}

func (this *WaveMixer) Name() string {
    return "wavemix"
}

func (this *WaveMixer) InitConfig(jobConf string) (err error) {
    return
}

func (this *WaveMixer) Do(req ufop.UfopRequest) (result interface{}, resultType int, contentType string, err error) {
    duration, parseErr := this.parseVideoDuration("/Users/Zhangjd/Downloads/01.mp4")
    if parseErr != nil {
        log.Error(parseErr)
        err = parseErr
        return
    }
    log.Info(duration)

    var wav WavForge
    wav.InitConfig()
    wav.CreateWave()
    output := wav.getWavData()
    // fmt.Println(output)
    
    userFile := "test.wav"
    // Create creates the named file with mode 0666 (before umask), truncating it if it already exists
    fout, err := os.Create(userFile)
    defer fout.Close()
    if err != nil {
        fmt.Println(userFile, err)
        return
    }
    fout.WriteString(output)


    return
}

func (this *WaveMixer) parse(cmd string) (format string, mime string, bucket string, url string, waveArr []string, err error) {
    return
}

// 获取视频长度
func (this *WaveMixer) parseVideoDuration(fileName string) (result int, err error){
    // prepare for ffmpeg
    mergeCmdParams := []string{
        "-i", fileName,
    }

    // execute command
    mergeCmd := exec.Command("ffmpeg", mergeCmdParams...)

    // Wait will close the pipe after seeing the command exit
    stdErrPipe, pipeErr := mergeCmd.StderrPipe()
    if pipeErr != nil {
        err = errors.New(fmt.Sprintf("open exec stderr pipe error, %s", pipeErr.Error()))
        return
    }

    // Starts the specified command but does NOT wait for it to complete
    startErr := mergeCmd.Start();
    if startErr != nil {
        err = errors.New(fmt.Sprintf("start ffmpeg command error, %s", startErr.Error()))
        return
    }

    // Reads from stdErrPipe until an error or EOF and returns the data it read
    stdErrData, readErr := ioutil.ReadAll(stdErrPipe)
    if readErr != nil {
        err = errors.New(fmt.Sprintf("read ffmpeg command stderr error, %s", readErr.Error()))
        return
    }

    // Waits for the command to exit. It must have been started by Start.
    waitErr := mergeCmd.Wait()
    if waitErr != nil {
        // err = errors.New(fmt.Sprintf("wait ffmpeg to exit error, %s", waitErr))
        // return
    }

    // regex
    pattern := "Duration: ([0-9:]+)"
    keyRegx := regexp.MustCompile(pattern)
    matchStr := keyRegx.FindString(string(stdErrData))
    if matchStr == "" {
        err = errors.New("Cannot retrive duration.")
        return
    }
    matchStr = strings.Replace(matchStr, "Duration: ", "", -1)
    arr := strings.Split(matchStr, ":")
    hour, err0 := strconv.Atoi(arr[0])
    minute, err1 := strconv.Atoi(arr[1])
    second, err2 := strconv.Atoi(arr[2])
    if err0 != nil || err1 != nil || err2 != nil {
        err = errors.New("Invalid duration time format.")
        return
    }
    result = hour * 3600 + minute * 60 + second
    return
}







