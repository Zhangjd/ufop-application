# ufop-demo
七牛ufop服务demo

### Configure
```
go get -u github.com/qiniu/api.v6
go get -u golang.org/x/text/encoding/simplifiedchinese
```

### Local test
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

### TODO
视频合成逻辑

### License
This project is published under MIT License. See the LICENSE file for more.




