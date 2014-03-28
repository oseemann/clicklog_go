package main

import (
    "os"
    "fmt"
    "time"
)

var timestamp_layout = "2006-01-02T15:04:05.00000"
var log_queue = make(chan LogRecord, 10000)

func log(r LogRecord) {
    log_queue <- r
}

func Logger() {
    for ;; {
        n := len(log_queue)
        if n == 0 {
            time.Sleep(100 * time.Millisecond)
        } else {
            file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_WRONLY, 0660)
            if err == nil {
                for i := 0; i < n; i++ {
                    r := <-log_queue
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
}

// vim: set filetype=go expandtab:
