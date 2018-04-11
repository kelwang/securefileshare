package handler

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/kelwang/securefileshare/ui"
)

const maxRetry = 3

// New http.Handler
func New(rootPath, passCode string) http.Handler {
	return &handler{
		rootPath: rootPath,
		passCode: passCode,
		tried:    0,
		session:  make(map[string]int64),
	}
}

type handler struct {
	rootPath string
	passCode string
	tried    int
	session  map[string]int64
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
		log.Println(err)
		fmt.Fprint(w, "Some error happened")
	}
}

var route = map[string]func(h *handler, w http.ResponseWriter, r *http.Request) (err error){
	"download": download,
	"destroy":  destroy,
}

func list(h *handler, w http.ResponseWriter, r *http.Request) (err error) {
	if !h.verifyRequest(w, r) {
		tmpl, err := template.New("password").Parse(ui.PasswordPage)
		if err != nil {
			log.Fatal("bad template")
		}
		err = tmpl.Execute(w, maxRetry-h.tried)
		return err
	}
	tmpl, err := template.New("download").Parse(ui.DownloadPage)
	if err != nil {
		log.Fatal("bad template")
	}
	fs, err := getFiles(h.rootPath)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, struct {
		FS      []string
		Session string
	}{
		FS:      fs,
		Session: r.URL.Query().Get("session"),
	})
}

func download(h *handler, w http.ResponseWriter, r *http.Request) (err error) {
	if !h.verifyRequest(w, r) {
		err = errors.New("unauthorized request")
		return
	}
	p := r.URL.Path[len("/download/"):]
	defer func(er *error) {
		if rr := recover(); rr != nil {
			*er = errors.New(string(debug.Stack()))
		}
	}(&err)
	println(p)
	w.Header().Set("Content-Type", "application/force-download")
	http.ServeFile(w, r, h.rootPath+"/"+p)
	return

}

func destroy(h *handler, w http.ResponseWriter, r *http.Request) (err error) {
	if h.verifyRequest(w, r) {
		log.Fatal("server is distroyed")
	}
	err = errors.New("unauthorized request")
	return

}

func (h *handler) authRequest(w http.ResponseWriter, r *http.Request) bool {
	if err := r.ParseForm(); err != nil {
		return false
	}
	code := strings.TrimSpace(r.Form.Get("code"))
	if code == "" {
		return false
	}
	if code == h.passCode {
		//set auth cookie
		h.setAuthCookie(w, r)
		return true
	}
	if h.tried+1 == maxRetry {
		log.Fatal("max retry has been reached")
	}
	h.tried++
	return false
}

func (h *handler) verifyRequest(w http.ResponseWriter, r *http.Request) bool {
	if session := r.URL.Query().Get("session"); session != "" {
		if exp, ok := h.session[session]; ok {
			if exp >= time.Now().Unix() {
				return true
			}
		}
	}
	return h.authRequest(w, r)
}

func getFiles(dir string) (files []string, err error) {
	var fs []os.FileInfo
	fs, err = ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	for _, v := range fs {
		if !v.IsDir() && v.Name()[0] != '.' {
			files = append(files, v.Name())
		}
	}
	return
}

func (h *handler) setAuthCookie(w http.ResponseWriter, r *http.Request) {
	str := randStringBytes(20)
	h.session[str] = time.Now().Unix() + 900
	http.Redirect(w, r, "/?session="+str, http.StatusFound)
}

const letterBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
