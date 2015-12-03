package main

import (
    "encoding/json"
    "io"
    "io/ioutil"
    "log"
    "net/http"
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

func main() {
    http.HandleFunc("/uop", demoHandler)
    err := http.ListenAndServe(":9100", nil)
    if err != nil {
        log.Fatal("Demo server failed to start:", err)
    }
}