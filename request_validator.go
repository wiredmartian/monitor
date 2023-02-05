package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/monitor/localstorage"
	"net/http"
)

type ReqValidator interface {
	validateRequest(l *localstorage.LocalStorage, c *gin.Context)
}

func validateRequest(l *localstorage.LocalStorage, c *gin.Context) {
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
