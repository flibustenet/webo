package webo

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type LogSender struct {
	sync.RWMutex
	name      string
	poste     string
	version   string
	url_log   string
	ch        chan []byte
	last_send time.Time
}

func NewLogSender(name, poste, version, url_log string) *LogSender {
	ch := make(chan []byte, 555)
	ls := &LogSender{name: name, poste: poste, version: version, url_log: url_log, ch: ch}
	go ls.go_send()
	return ls
}
func (w *LogSender) go_send() {
	buf := &bytes.Buffer{}
	tickChan := time.NewTicker(time.Second * 60).C
	for {
		select {
		case <-tickChan:
			if buf.Len() > 0 {
				w.send(buf.String())
				buf.Reset()
			}
		case b := <-w.ch:
			buf.Write(b)
		}
	}
}
func (w *LogSender) Write(output []byte) (int, error) {
	w.ch <- output
	return len(output), nil
}
func (w *LogSender) send(buf string) {
	_, err := http.PostForm(w.url_log+"/"+w.name+"_"+w.version+"/archive", url.Values{"name": {w.name},
		"version": {w.version},
		"poste":   {w.poste},
		"msg":     {buf}})
	if err != nil {
		fmt.Println(err)
	}
}
