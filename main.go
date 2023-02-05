package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/monitor/localstorage"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func main() {
	r := gin.Default()
	m := melody.New()
	l := localstorage.New()

	m.Upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	r.POST("/webhook", func(c *gin.Context) {
		key := uuid.New().String()
		_ = l.AddToStorage(key)
		cache := l.GetAllFromStorage()
		c.JSON(http.StatusOK, gin.H{"data": cache})
	})

	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "public/index.html")
	})

	r.GET("/client/:id/connect", func(c *gin.Context) {
		ValidateRequest(l, c)
		m.HandleRequest(c.Writer, c.Request)
	})

	r.POST("/notify/:id", func(c *gin.Context) {
		ValidateRequest(l, c)
		_ = HandleNotification(c, m)
		c.Status(200)
		return
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.BroadcastFilter(msg, func(q *melody.Session) bool {
			var data map[string]interface{}
			json.Unmarshal(msg, data)
			fmt.Println(data)
			return s.Request.URL.Path == q.Request.URL.Path && s != q
		})
	})

	r.Run(":1111")
}

func HandleNotification(c *gin.Context, m *melody.Melody) error {
	id, _ := c.Params.Get("id")

	response := struct {
		Uri string
		Id  string
	}{Uri: c.Request.URL.Path, Id: id}

	resBody := new(bytes.Buffer)
	json.NewEncoder(resBody).Encode(response)
	err := m.BroadcastFilter(resBody.Bytes(), func(s *melody.Session) bool {
		return s.Request.URL.Path == fmt.Sprintf("/client/%v/connect", id)
	})
	return err
}

func ValidateRequest(l *localstorage.LocalStorage, c *gin.Context) {
	key := ""
	if paramId, found := c.Params.Get("id"); !found {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required webhook id parameter"})
		return
	} else {
		key = paramId
	}
	if !l.Exists(key) {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("webhook id: %v does not exist", key)})
		return
	}
}
