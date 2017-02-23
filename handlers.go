package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/hpcloud/tail"
	"golang.org/x/net/websocket"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.New("base").Parse(opts.Template)
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
