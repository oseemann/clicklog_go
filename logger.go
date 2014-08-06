package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

var timestamp_layout = "2006-01-02T15:04:05.00000"
var log_queue = make(chan LogRecord, 10000)
var zip_it = make(chan string)

func Log(r LogRecord) {
	log_queue <- r
}

func Zipper(fname string) {
	cmd := exec.Command("gzip", fname)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func CheckZip(log_rotate_limit int64, zip_log_files bool) {
	stat, err := os.Stat("log.txt")
	if err == nil {
		if stat.Size() > log_rotate_limit {
			log.Print("Log size: ", stat.Size())
			tmpfile, err := ioutil.TempFile(".", "log_rotate_")
			if err == nil {
				log.Print("Rotate temp file: %s", tmpfile.Name())
				tmpfile.Close()
				err := os.Rename("log.txt", tmpfile.Name())
				if err == nil {
					file, _ := os.Create("log.txt")
					file.Close()
					log.Print("Recreated empty log.txt")
				} else {
					log.Print("Error moving log.txt: ", err)
					os.Exit(1)
				}
				if zip_log_files {
					go Zipper(tmpfile.Name())
				}
			} else {
				log.Print("Cannot create temp file for rotating ", err)
				os.Exit(1)
			}
		}
	} else {
		log.Print("Cannot stat log.txt: ", err)
		os.Exit(1)
	}
}

func Logger(log_rotate_limit int64, zip_log_files bool) {
	for {
		n := len(log_queue)
		if n == 0 {
			time.Sleep(100 * time.Millisecond)
		} else {
			CheckZip(log_rotate_limit, zip_log_files)

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
				log.Print(err)
			}
		}
	}
}

// vim: set filetype=go expandtab:
