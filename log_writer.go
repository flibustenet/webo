package webo

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type DayWriter struct {
	sync.RWMutex
	cur_day    time.Weekday
	fp         *os.File
	force_time time.Time // uniquement pour unitest
}

var daywriter = &DayWriter{}

func GetDayWriter() *DayWriter {
	return daywriter
}

func (w *DayWriter) now() time.Time {
	if !w.force_time.IsZero() {
		return w.force_time
	}
	return time.Now()
}

// création du fichier si inexistant ou si celui existant
// est plus vieux de 24h
func (w *DayWriter) Write(output []byte) (int, error) {
	w.Lock()
	defer w.Unlock()
	now := w.now()
	if w.fp == nil || w.cur_day != now.Weekday() {
		_, err := os.Stat("logs")
		if os.IsNotExist(err) {
			err = os.Mkdir("logs", 0700)
		}
		if err != nil {
			return 0, err
		}

		path := fmt.Sprintf("logs/log_%d.txt", now.Weekday())
		info, err := os.Stat(path)
		if err != nil && !os.IsNotExist(err) {
			return 0, err
		}
		if os.IsNotExist(err) || now.Sub(info.ModTime()).Hours() > 24. {
			w.fp, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0600)
			if err != nil {
				return 0, fmt.Errorf("Impossible de créer le log %s : %s", path, err)
			}
		} else {
			w.fp, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
			if err != nil {
				return 0, fmt.Errorf("Impossible d'ouvrir le log %s : %s", path, err)
			}
		}
		w.cur_day = now.Weekday()
	}
	i, err := w.fp.Write(output)
	if err != nil {
		return i, err
	}
	err = w.fp.Sync()
	return i, err
}
func (w *DayWriter) Close() error {
	w.Lock()
	defer w.Unlock()
	if w.fp != nil {
		err := w.fp.Close()
		w.fp = nil
		return err
	}
	return nil
}
