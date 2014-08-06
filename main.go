package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

type LogRecord struct {
	ts        time.Time
	client    string
	target    string
	referer   string
	useragent string
}

type Options struct {
	port             int
	log_rotate_limit int64
	zip_log_files    bool
}

var thisHost = "fry:8080"

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi!")
}

func make_logrecord(request *http.Request, dest string) LogRecord {

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

	return (logr)
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

func ReadOptions() Options {
	opt := Options{}

	flag.IntVar(&opt.port, "port", 8080, "TCP Port to listen on")
	flag.Int64Var(&opt.log_rotate_limit, "log_rotate_limit", 100000000,
		"Size limit that triggers log file rotation")
	flag.BoolVar(&opt.zip_log_files, "zip_log", true, "Zip rotated log files")

	flag.Parse()

	log.Print("log_rotate_limit = ", opt.log_rotate_limit)
	log.Print("zip_log_files = ", opt.zip_log_files)

	return opt
}

func main() {
	opt := ReadOptions()
	go Logger(opt.log_rotate_limit, opt.zip_log_files)
	http.HandleFunc("/", root)
	http.HandleFunc("/l", clickHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", opt.port), nil))
}

// vim: set filetype=go expandtab:
