package main

import (
    "fmt"
    "time"
    "net/url"
    "net/http"
)

type LogRecord struct {
    ts time.Time
    client string
    target string
    referer string
    useragent string
}

var thisHost = "fry:8080"

func root(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi!")
}

func make_logrecord(request *http.Request, dest string) (LogRecord) {

    // Referer
    var referer string = "-"
    raw_referer := request.Header["Referer"]
    if len(raw_referer) > 0 {
        referer = raw_referer[0]
    }

    // User Agent
    var useragent string = "-"
    raw_useragent := request.Header["User-Agent"]
    if len(raw_useragent) > 0 {
        useragent = raw_useragent[0]
    }

    logr := LogRecord{
        time.Now(),
        request.RemoteAddr,
        dest,
        referer,
        useragent,
    }

    return(logr)
}

func clickHandler(w http.ResponseWriter, r *http.Request) {
    raw_dest := r.FormValue("r")
    dest, err := url.QueryUnescape(raw_dest)
    if err != nil {
        // default dest in case of error
        dest = "http://www.google.com/"
    }

    // check sanity of target url
    // TODO

    // log request
    logrecord := make_logrecord(r, dest)
    Log(logrecord)
    
    http.Redirect(w, r, dest, 302)
}

func main() {
    go Logger()
    http.HandleFunc("/", root)
    http.HandleFunc("/l", clickHandler)
    http.ListenAndServe(":8080", nil)
}

// vim: set filetype=go expandtab:
