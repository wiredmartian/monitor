package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
	"net/http"
)

func main() {
	r := gin.Default()
	m := melody.New()

	m.Upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "public/index.html")
	})

	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	r.GET("/channel/:id/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	r.POST("/notify/:id", func(c *gin.Context) {
		_ = HandleNotification(c, m)
		c.Status(200)
		return
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.BroadcastFilter(msg, func(q *melody.Session) bool {
			fmt.Println(s.Request.URL.Path)
			return s.Request.URL.Path == q.Request.URL.Path && s != q
		})
	})

	r.Run(":5000")
}

func HandleNotification(c *gin.Context, m *melody.Melody) error {
	c.Status(200)

	id, found := c.Params.Get("id")
	if !found {
		return fmt.Errorf("id not found in params")
	}
	response := struct {
		Uri string
		Id  string
	}{Uri: c.Request.URL.Path, Id: id}

	resBody := new(bytes.Buffer)
	json.NewEncoder(resBody).Encode(response)

	err := m.BroadcastFilter(resBody.Bytes(), func(s *melody.Session) bool {
		return s.Request.URL.Path == fmt.Sprintf("/channel/%v/ws", id)
	})
	if err != nil {
		return err
	}
	return nil
}
