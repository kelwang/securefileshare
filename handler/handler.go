package handler

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/kelwang/securefileshare/ui"
)

const maxRetry = 3

// New http.Handler
func New(rootPath, secret, passCode string) http.Handler {
	return &handler{
		rootPath: rootPath,
		secret:   []byte(secret),
		passCode: passCode,
		tried:    0,
	}
}

type handler struct {
	rootPath string
	secret   []byte
	passCode string
	tried    int
}

// ServeHTTP will implement the net http.Handler interface
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i := strings.Index(r.URL.Path[1:], "/")
	var action func(h *handler, w http.ResponseWriter, r *http.Request) (err error)
	var ok bool
	if i == -1 {
		action = list
		goto run
	}
	action, ok = route[r.URL.Path[1:i+1]]
	if !ok {
		action = list
	}
run:
	err := action(h, w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Some error happened")
	}
}

var route = map[string]func(h *handler, w http.ResponseWriter, r *http.Request) (err error){
	"download": download,
	"destroy":  destroy,
}

func list(h *handler, w http.ResponseWriter, r *http.Request) (err error) {
	if !h.verifyRequest(r) {
		tmpl, err := template.New("password").Parse(ui.PasswordPage)
		if err != nil {
			log.Fatal("bad template")
		}
		err = tmpl.Execute(w, nil)
		return err
	}

	return nil
}

func download(h *handler, w http.ResponseWriter, r *http.Request) (err error) {
	if !h.verifyRequest(r) {
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
	if h.verifyRequest(r) {
		log.Fatal("server is distroyed")
	}
	err = errors.New("unauthorized request")
	return

}

func (h *handler) authRequest(r *http.Request) bool {
	if err := r.ParseForm(); err != nil {
		return false
	}
	code := strings.TrimSpace(r.Form.Get("code"))
	if code == "" {
		return false
	}
	if code == h.passCode {
		return true
	}
	if h.tried+1 == maxRetry {
		log.Fatal("max retry has been reached")
	}
	h.tried++
	return false
}

func (h *handler) verifyRequest(r *http.Request) bool {
	return h.authRequest(r)
}
