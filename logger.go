package main

import (
    "os"
    "fmt"
    "time"
    "sync"
)

var logbuf1 = make([]LogRecord, 1000)
var logbuf2 = make([]LogRecord, 1000)
var curr_ = &logbuf1
var logbuf_ []LogRecord = (*curr_)[:0]
var L sync.Mutex

var timestamp_layout = "2006-01-02T15:04:05.00000"

func log(r LogRecord) {
    L.Lock()
    logbuf_ = append(logbuf_, r)
    L.Unlock()
}

func Logger() {
    for ;; {
        if len(logbuf_) == 0 {
            time.Sleep(10 * time.Millisecond)
            continue
        }

        mybuf := logbuf_
        // switch buffers
        if curr_ == &logbuf1 {
            curr_ = &logbuf2
            logbuf_ = (*curr_)[:0]
        } else {
            curr_ = &logbuf1
            logbuf_ = (*curr_)[:0]
        }

        // dump buffer
        file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_WRONLY, 0660)
        if err == nil {
            for _, r := range mybuf {
                fmt.Fprintf(file,
                    "%s\t%s\t%s\t%s\t%s\n",
                    r.ts.Format(timestamp_layout),
                    r.client,
                    r.target,
                    r.referer,
                    r.useragent)
            }
            file.Close()
        } else {
            fmt.Println(err) 
        }
    }
}

// vim: set filetype=go expandtab:
