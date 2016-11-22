package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/Akagi201/light"
	"github.com/Sirupsen/logrus"
	"github.com/arschles/go-bindata-html-template"
	"github.com/hpcloud/tail"
	"github.com/jessevdk/go-flags"
	"golang.org/x/net/websocket"
)

var opts struct {
	ListenAddr string `long:"listen" default:"0.0.0.0:8327" description:"HTTP address and port to listen at"`
	Template   string `long:"template" default:"data/template/index.html" description:"the template base file"`
	Log        string `long:"log" description:"the log file to tail -f"`
}

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.Level = logrus.InfoLevel
	f := new(logrus.TextFormatter)
	f.TimestampFormat = "2006-01-02 15:04:05"
	f.FullTimestamp = true
	log.Formatter = f
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.New("base", Asset).Parse(opts.Template)
	if err != nil {
		log.Printf("Template parse failed, err: %v", err)
		return
	}
	v := struct {
		Host string
		Log  string
	}{
		r.Host,
		opts.Log,
	}
	if err = t.Execute(w, &v); err != nil {
		log.Printf("Template execute failed, err: %v", err)
		return
	}
}

func handleTail(w http.ResponseWriter, r *http.Request) {
	t, err := tail.TailFile(opts.Log, tail.Config{Follow: true, ReOpen: true, Location: &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END}})
	if err != nil {
		log.Printf("tail file failed, err: %v", err)
		return
	}
	for line := range t.Lines {
		log.Println(line.Text)
		w.Write([]byte(line.Text))
	}
}

func handleFollow(ws *websocket.Conn) {
	t, err := tail.TailFile(opts.Log, tail.Config{Follow: true, ReOpen: true, Location: &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END}})
	if err != nil {
		log.Printf("tail file failed, err: %v", err)
		return
	}
	for line := range t.Lines {
		log.Println(line.Text)
		ws.Write([]byte(line.Text))
	}
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		if !strings.Contains(err.Error(), "Usage") {
			log.Fatalf("error: %v", err)
		} else {
			return
		}
	}

	app := light.New()

	app.Get("/tail", handleTail)
	app.Get("/follow", websocket.Handler(handleFollow))
	app.Get("/", handleHome)

	log.Printf("HTTP listening at: %v", opts.ListenAddr)
	app.Listen(opts.ListenAddr)
}
