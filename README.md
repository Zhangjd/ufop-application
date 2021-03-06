# ufop-application [![Build Status](https://travis-ci.org/Zhangjd/ufop-application.svg)](https://travis-ci.org/Zhangjd/ufop-application)
基于七牛ufop的自定义处理程序

## Installation

```
go get -u github.com/qiniu/api.v6
go get -u golang.org/x/text/encoding/simplifiedchinese
```

## Build Tutorial

* 设置环境变量 $GOPATH 和 $GOBIN
* 修改build.sh中的环境变量
* 重命名deploy文件夹中 \*.conf.example为 \*.conf, 写入自定义配置
* 本地调试版本请进入src目录, 执行go install qufop.go
* 部署远程版本请进入src目录, 执行./build.sh

## Local Test

Command: qufop \<qufop-config-file>

```
POST /uop HTTP/1.1
Content-Type: application/json
{
    "cmd": "<ufop>/<param>",
    "src": {
        "url": "http://<host>:<port>/<path>",
        "mimetype": "<mimetype>",
        "fsize": <filesize>
    }
}
```

## Remote Ufop Call

Reference: 

http://developer.qiniu.com/docs/v6/api/reference/fop/pfop/pfop.html
http://developer.qiniu.com/docs/v6/api/reference/fop/pfop/prefop.html

1. 触发持久化处理(pfop)，接口返回的\<persistentId>
2. 七牛服务端按顺序完成所有指定的云处理操作后，会将处理结果状态提交到\<persistentNotifyUrl>指向的网址，用作异步回调处理
3. 客户端可以使用\<persistentId>来主动查询持久化处理的执行状态

## Module Usage

### 视频合成 (画面左右放置, 分别占一半大小)

```
Option:
st-videomerge/format/<format>/mime/<mime>/bucket/<bucket>/url/<url>

Relevant command execution: 
ffmpeg -y -i prevideo.mp4 -i sufvideo.mp4  -filter_complex "[0:v:0]pad=iw*2:ih[bg]; [bg][1:v:0]overlay=w" output.mp4

Params - description:
    format - output format (suggest mp4)
    mime   - <base64_encode required> output Mime-type (suggest video/mp4)
    bucket - <base64_encode required> bucket of the second video file in Qiniu
    url    - <base64_encode required> url of the second video
```

### 视频声波合成

```
// TODO
```

## License

This project is published under MIT License. See the LICENSE file for more.




