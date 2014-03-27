package main

import (
    "fmt"
    "time"
    "net/http"
)

type LogRecord struct {
    ts time.Time
    client string
    target string
    referer string
}

var thisHost = "fry:8080"

func root(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi!")
}

func make_logr(request *http.Request, dest string) (LogRecord) {

    var referer string = "-"
    r_referer := request.Header["Referer"]

    if len(r_referer) > 0 {
        referer = r_referer[0]
    }

    logr := LogRecord{
        time.Now(),
        request.RemoteAddr,
        dest,
        referer,
    }

    return(logr)
}

func clickHandler(w http.ResponseWriter, r *http.Request) {
    redirect_dest := r.FormValue("r")

    // check sanity of target url
    // TODO

    // log request
    logr := make_logr(r, redirect_dest)
    log(logr)
    
    http.Redirect(w, r, redirect_dest, 302)
}

func main() {
    go Logger()
    http.HandleFunc("/", root)
    http.HandleFunc("/l", clickHandler)
    http.ListenAndServe(":8080", nil)
}

// vim: set filetype=go expandtab:
