# ufop-application
基于七牛ufop的自定义处理实现

## Configure

```
go get -u github.com/qiniu/api.v6
go get -u golang.org/x/text/encoding/simplifiedchinese
```

## Local test

```
POST /uop HTTP/1.1
Content-Type: application/json
{
    "cmd": "<ufop-command>/<param>",
    "src": {
        "url": "http://<host>:<port>/<path>",
        "mimetype": "<mimetype>",
        "fsize": <filesize>
    }
}
```

## Remote 

Reference: 

http://developer.qiniu.com/docs/v6/api/reference/fop/pfop/pfop.html
http://developer.qiniu.com/docs/v6/api/reference/fop/pfop/prefop.html

1. 触发持久化处理(pfop)，接口返回的\<persistentId>
2. 七牛服务端按顺序完成所有指定的云处理操作后，会将处理结果状态提交到\<persistentNotifyUrl>指向的网址，供用户
3. 用户可以使用\<persistentId>来主动查询持久化处理的执行状态

## Module usage

### 视频画面左右合成

```
Option:
st-videomerge/format/<format>/mime/<mime>/bucket/<bucket>/url/<url>

Command execution: 
ffmpeg -y -i prevideo.mp4 -i sufvideo.mp4  -filter_complex "[0:v:0]pad=iw*2:ih[bg]; [bg][1:v:0]overlay=w" output.mp4

Params - description:
    format - output format (suggest mp4)
    mime   - <base64_encode required> output Mime-type (suggest video/mp4)
    bucket - <base64_encode required> bucket of the second video file in Qiniu
    url    - <base64_encode required> url of the second video
```

### 微波合成
```
TODO
```

## License

This project is published under MIT License. See the LICENSE file for more.




