package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/kelwang/securefileshare/handler"
)

func main() {
	rootPathFlag := flag.String("path", "", "root path will the folder where your docs located to download")
	passFlag := flag.String("code", "", "passcode for your client to gain access")
	portFlag := flag.Int("port", 0, "if empty, the default port will be 8080")
	flag.Parse()
	rootPath := *rootPathFlag

	passCode := *passFlag
	if passCode == "" {
		log.Fatal("pass code can't be empty")
	}
	port := "8080"
	if pt := *portFlag; pt > 0 {
		port = strconv.Itoa(pt)
	}

	println("secure file share server is started on port " + port)

	http.Handle("/", handler.New(rootPath, passCode))
	http.ListenAndServe(":"+port, nil)
}
