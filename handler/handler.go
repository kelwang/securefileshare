package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
)

// New http.Handler
func New(rootPath, secret, passCode string) http.Handler {
	return &handler{
		rootPath: rootPath,
		secret:   []byte(secret),
		passCode: passCode,
	}
}

type handler struct {
	rootPath string
	secret   []byte
	passCode string
}

// ServeHTTP will implement the net http.Handler interface
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i := strings.Index(r.URL.Path[1:], "/")
	if i == -1 {
		return
	}
	action, ok := route[r.URL.Path[1:i+1]]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "bad request")
		return
	}
	err := action(h, w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Some error happened")
	}
}

var route = map[string]func(h *handler, w http.ResponseWriter, r *http.Request) (err error){
	"download": download,
}

func download(h *handler, w http.ResponseWriter, r *http.Request) (err error) {
	if !verifyRequest(r) {
		err = errors.New("unauthorized request")
		return
	}
	p := r.URL.Path[len("/download/"):]
	defer func(er *error) {
		if rr := recover(); rr != nil {
			*er = errors.New(string(debug.Stack()))
		}
	}(&err)
	w.Header().Set("Content-Type", "application/force-download")
	http.ServeFile(w, r, p)
	return

}

func destroy(h *handler, w http.ResponseWriter, r *http.Request) (err error) {
	if verifyRequest(r) {
		log.Fatal("server is distroyed")
	}
	err = errors.New("unauthorized request")
	return

}

func verifyRequest(r *http.Request) bool {
	return false
}
