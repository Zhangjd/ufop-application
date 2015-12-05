package foo

import (
    "encoding/json"
    "io"
    "io/ioutil"
    "net/http"
    "github.com/qiniu/api.v6/auth/digest"
    "github.com/qiniu/log"
    "ufop"
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

func demoHandler(w http.ResponseWriter, req *http.Request) {
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        w.WriteHeader(400)
        log.Println("read request body failed:", err)
        return
    }
    var args ReqArgs
    err = json.Unmarshal(body, &args)
    if err != nil {
        w.WriteHeader(400)
        log.Println("invalid request body:", err)
        return
    }
    resp, err := http.Get(args.Src.Url)
    if err != nil {
        w.WriteHeader(400)
        log.Println("fetch resource failed:", err)
        return
    }
    defer resp.Body.Close()
    var buf = make([]byte, 512)
    io.ReadFull(resp.Body, buf)
    contentType := http.DetectContentType(buf)
    w.Write([]byte(contentType))
}

func (this *VideoMerger) Do(req ufop.UfopRequest) (result interface{}, resultType int, contentType string, err error) {
    http.HandleFunc("/uop", demoHandler)
    httpErr := http.ListenAndServe(":9100", nil)
    if httpErr != nil {
        log.Fatal("Demo server failed to start:", err)
    }
    return
}




