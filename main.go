package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
)

var tmplLanding *template.Template
var tmplRestart *template.Template

type data struct {
	Title string
	Text  string
}

type link struct {
	Desc string
	Url  string
}

type links struct {
	Title string
	Links []link
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func funcName() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return fmt.Sprintf("%s:%d %s", filepath.Base(frame.File), frame.Line, frame.Function)
}

func homepage(w http.ResponseWriter, r *http.Request) {
	// Check if exact match ("/" is catch-all)
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		log.Printf("Path %s does not have a route, 404 sent", r.URL.Path)
		return
	}

	data := links{
		Title: "Select service",
		Links: []link{
			{"Restart IPTV", "/restart-iptv"},
			// {"Reboot router", "/reboot"},
		},
	}

	err := tmplLanding.Execute(w, data)
	check(err)

	log.Printf("Endpoint %s hit: %s", r.URL.Path, funcName())
}

func restartService(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("systemctl", "restart", "udm-iptv")
	bufStdout := new(bytes.Buffer)
	bufStderr := new(bytes.Buffer)

	cmd.Stdout = bufStdout
	cmd.Stderr = bufStderr

	err := cmd.Run()

	var text *bytes.Buffer
	if err != nil {
		text = bufStderr
	} else {
		text = bufStdout
	}
	data := data{
		Title: "Restart service",
		Text:  text.String(),
	}

	err = tmplRestart.Execute(w, data)
	check(err)

	log.Printf("Endpoint %s hit: %s", r.URL.Path, funcName())
}

func handleRequests(certPath string, keyPath string, listenPort int) {
	http.HandleFunc("/", homepage)
	http.HandleFunc("/restart-iptv", restartService)
	// TODO: Reboot
	// http.HandleFunc("/reboot", restartService)

	if certPath != "" && keyPath != "" {
		log.Printf("Starting HTTPS server on port %d", listenPort)
		log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", listenPort),
			certPath,
			keyPath,
			nil))
	} else {
		log.Printf("Starting HTTP server on port %d", listenPort)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil))
	}
}

func main() {
	certPath := flag.String("cert", "", "path to certificate")
	keyPath := flag.String("key", "", "path to certificate key")
	listenPort := flag.Int("port", 12345, "port to listen to")
	// Only support no flags or all flags
	flag.Parse()
	if (*certPath != "" && *keyPath == "") || (*certPath == "" && *keyPath != "") {
		panic("You have to specify paths to both cert and key if you want to serve via TLS")
	}
	var err error
	tmplLanding, err = template.New("landing.gohtml").Funcs(sprig.FuncMap()).ParseFiles("landing.gohtml", "head.gohtml", "stylesheet.gohtml")
	check(err)
	tmplRestart, err = template.New("restart.gohtml").Funcs(sprig.FuncMap()).ParseFiles("restart.gohtml", "head.gohtml", "stylesheet.gohtml")
	check(err)

	handleRequests(*certPath, *keyPath, *listenPort)
}
