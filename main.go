package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/kelwang/securefileshare/handler"
)

func main() {
	rootPathFlag := flag.String("p", "", "root path will the folder where your docs located to download")
	secretFlag := flag.String("s", "", "secret string used for encryption")
	passFlag := flag.String("c", "", "passcode for your client to gain access")

	flag.Parse()
	rootPath := *rootPathFlag
	secret := *secretFlag

	if secret == "" {
		log.Fatal("secret can't be empty")
	}

	passCode := *passFlag
	if passCode == "" {
		log.Fatal("pass code can't be empty")
	}

	http.Handle("/", handler.New(rootPath, secret, passCode))
	http.ListenAndServe(":8080", nil)
}
