/**
 * Author: Zhangjd
 * Date: December 17th, 2015
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
    "math"
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
    duration, parseErr := this.parseVideoDuration("00.mp4")
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
    
    // Create creates the named file with mode 0666 (before umask), truncating it if it already exists
    userFile := "test.wav"
    fout, err := os.Create(userFile)
    defer fout.Close()
    if err != nil {
        fmt.Println(userFile, err)
        return
    }
    fout.Write(output)

    rscode := [] string {"uv8e463l175lsiijdq4t"}
    tempTxtFile, _ := this.createTxtFile(duration, rscode)
    this.mergeWavIntoMp3(tempTxtFile)
    defer os.Remove(tempTxtFile)

    return
}

func (this *WaveMixer) parse(cmd string) (format string, mime string, bucket string, url string, waveArr []string, err error) {
    return
}

func (this *WaveMixer) createTxtFile (duration int, rscode []string) (txtFile string, err error) {
    repeatCount := math.Floor(float64(duration) / float64(len(rscode)) / 1.74)
    tmpTxtFile, sErr := ioutil.TempFile("", "create_sound_")
    if sErr != nil {
        err = errors.New(fmt.Sprintf("open temp file failed, %s", sErr.Error()))
        return
    }
    for i := 0.0; i < repeatCount; i++ {
        _, sCpErr := tmpTxtFile.WriteString("file '/Users/Zhangjd/IdeaProjects/ufop/bin/test.wav' \n")
        if sCpErr != nil {
            err = errors.New(fmt.Sprintf("save second first tmp file failed, %s", sCpErr.Error()))
            return
        }
    }
    txtFile = tmpTxtFile.Name()
    tmpTxtFile.Close()
    return
}

func (this *WaveMixer) mergeWavIntoMp3 (txtFile string) () {
    outputMp3FileName := "/Users/Zhangjd/IdeaProjects/ufop/bin/output.mp3"
    mergeCmdParams := []string{
        "-y",
        "-v", "error",
        "-f", "concat",
        "-i", txtFile,
        "-ar", "44100",
        "-ab", "128k",
        outputMp3FileName,
    }
    mergeCmd := exec.Command("ffmpeg", mergeCmdParams...)
    stdErrPipe, pipeErr := mergeCmd.StderrPipe()
    if pipeErr != nil {
        fmt.Sprintf("open exec stderr pipe error, %s", pipeErr.Error())
    }
    startErr := mergeCmd.Start();
    if startErr != nil {
        fmt.Sprintf("start ffmpeg command error, %s", startErr.Error())
    }
    stdErrData, readErr := ioutil.ReadAll(stdErrPipe)
    if readErr != nil {
        fmt.Sprintf("read ffmpeg command stderr error, %s", readErr.Error())
    }
    if string(stdErrData) != "" {
        log.Info(string(stdErrData))
    }
    mergeCmd.Wait()
}

// 获取视频长度
func (this *WaveMixer) parseVideoDuration(fileName string) (result int, err error) {
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







