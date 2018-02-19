package webo

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"
)

func Test_log_writer(t *testing.T) {
	now := time.Now()
	w := GetDayWriter()
	w.force_time = now
	_, err := fmt.Fprintf(w, "t=%s\n", w.force_time)
	if err != nil {
		log.Fatal(err)
	}
	w.force_time = now.Add(time.Second * 3600)
	_, err = fmt.Fprintf(w, "t=%s\n", w.force_time)
	if err != nil {
		log.Fatal(err)
	}
	w.force_time = now.Add(time.Second * 3600 * 2)
	_, err = fmt.Fprintf(w, "t=%s\n", w.force_time)
	if err != nil {
		log.Fatal(err)
	}
	w.force_time = now.AddDate(0, 0, 1)
	_, err = fmt.Fprintf(w, "t=%s\n", w.force_time)
	if err != nil {
		log.Fatal(err)
	}

	w.Close()
	before := w.force_time.AddDate(0, 0, -7)
	err = os.Chtimes("logs/log_"+strconv.Itoa(int(w.force_time.Weekday()))+".txt", before, before)
	_, err = fmt.Fprintf(w, "seul\n")
	if err != nil {
		log.Fatal(err)
	}
	os.RemoveAll("logs")
}
