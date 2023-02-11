package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/monitor/localstorage"
	"log"
	"net/http"
	"strings"

	"gopkg.in/olahol/melody.v1"
)

func main() {
	m := melody.New()
	l := localstorage.New()

	m.Upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		key := uuid.New().String()
		_ = l.AddToStorage(key)
		cache := l.GetAllFromStorage()
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(struct {
			Data map[string]localstorage.StorageItem
		}{Data: cache})
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html")
	})

	http.HandleFunc("/client/:id/connect", func(w http.ResponseWriter, r *http.Request) {
		err := m.HandleRequest(w, r)
		if err != nil {
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println(err)
		}
	})

	http.HandleFunc("/notify/:id", func(w http.ResponseWriter, r *http.Request) {
		_ = HandleNotification(r, m)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		return
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.BroadcastFilter(msg, func(q *melody.Session) bool {
			return s.Request.URL.Path == q.Request.URL.Path && s != q
		})
	})

	err := http.ListenAndServe(":1111", nil)

	if err != nil {
		log.Fatal(err)
	}
}

func HandleNotification(r *http.Request, m *melody.Melody) error {
	id := strings.TrimPrefix(r.URL.Path, "/notify/")

	response := struct {
		Uri string
		Id  string
	}{Uri: r.URL.Path, Id: id}

	resBody := new(bytes.Buffer)
	json.NewEncoder(resBody).Encode(response)
	err := m.BroadcastFilter(resBody.Bytes(), func(s *melody.Session) bool {
		return s.Request.URL.Path == fmt.Sprintf("/client/%v/connect", id)
	})
	return err
}

//func ValidateRequest(l *localstorage.LocalStorage, r *http.Request) {
//	key := ""
//	if paramId, found := r.URL.Query().Get("id"); !found {
//		fmt.Println("missing required webhook id parameter")
//		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required webhook id parameter"})
//		return
//	} else {
//		key = paramId
//	}
//	if !l.Exists(key) {
//		fmt.Printf("webhook id: %v does not exist", key)
//		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("webhook id: %v does not exist", key)})
//		return
//	}
//}
