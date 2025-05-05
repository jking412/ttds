package app

import (
	"awesomeProject/internal/usecase"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func CreateContainerHandler(c *gin.Context) {

	// Get user ID from context (assuming it's set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}

	// Get template ID from query parameter
	templateID, err := strconv.ParseUint(c.Query("template_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template_id"})
		return
	}

	err = usecase.NewContainerService().CreateContainer(userID.(uint), uint(templateID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "container creation started"})
}

func GetContainerStatusHandler(c *gin.Context) {

	// Get user ID from context (assuming it's set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}

	templateID, err := strconv.ParseUint(c.Param("template_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template_id"})
		return
	}

	channel, err := usecase.NewContainerService().GetChannel(userID.(uint), uint(templateID))
	if err != nil {
		// 如果是不存在，而不是错误，则返回一个写入Running的channel，目前不判断，直接视为不存在而不是错误
		channel = make(chan string, 100)
		go func() {
			channel <- "Running"
			time.Sleep(time.Second * 3)
			close(channel)
		}()
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	for {
		select {
		case msg, ok := <-channel:
			if !ok {
				return
			}
			fmt.Fprintf(c.Writer, "%s", msg)
			c.Writer.Flush()
		case <-c.Writer.CloseNotify():
			return
		}
	}
}

func GetContainerHandler(c *gin.Context) {

	// Get user ID from context (assuming it's set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}

	templateID, err := strconv.ParseUint(c.Param("template_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template_id"})
		return
	}

	container, err := usecase.NewContainerService().GetContainer(userID.(uint), uint(templateID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, container)
}
