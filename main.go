package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func sse(w http.ResponseWriter, r *http.Request) {
	a := []string{"good morning", "hello", "good evening", "good night"}
	flusher, _ := w.(http.Flusher)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Controle", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	t := time.NewTicker(1 * time.Second)

	defer t.Stop()

	go func() {
		cnt := 1

		for {
			select {
			case <-t.C:
				rand.Seed(time.Now().UnixNano())
				m := map[string]interface{}{
					"message": a[rand.Intn(4)],
					"count":   cnt,
				}

				s, _ := json.Marshal(m)
				fmt.Fprintf(w, "data: %v\n\n", string(s))
				cnt++
				flusher.Flush()
			}
		}
	}()

	notify := w.(http.CloseNotifier).CloseNotify()
	<-notify
	log.Println("connection has closed")
}

func main() {
	http.HandleFunc("/event", sse)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.ListenAndServe(":8080", nil)
}
