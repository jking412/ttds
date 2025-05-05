package app

import (
	"awesomeProject/internal/usecase"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
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

	channel, err := usecase.NewContainerService().GetChannel(userID.(uint), uint(templateID), usecase.ContainerCreate)
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
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-channel; ok {
			// print time and msg
			//fmt.Printf("%s: %s\n", time.Now().Format("2006-01-02 15:04:05"), msg)
			c.SSEvent("message", msg)
			return true
		}
		return false
	})

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

	accessUrl := fmt.Sprintf("http://%s:3000/?tkn=%s", container.IPAddress, container.Token)

	c.JSON(http.StatusOK, gin.H{
		"container": container,
		"accessUrl": accessUrl,
	})
}

func CheckContainerHandler(c *gin.Context) {
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

	err = usecase.NewContainerService().CheckContainer(userID.(uint), uint(templateID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// TODO: 应该先直接获取channel一次，如果失败再重复获取，这样用户请求会更快得到结果

	// 重复获取channel，持续3s,每0.5s获取一次，如果3s内没有获取到，则返回错误
	tick := time.NewTicker(time.Millisecond * 500)
	timeout := time.After(time.Second * 3)

	var channel chan string

	for {
		var ok = false
		select {
		case <-tick.C:
			channel, err = usecase.NewContainerService().GetChannel(userID.(uint), uint(templateID), usecase.ContainerExec)
			if err != nil {
				logrus.Warnf("GetChannel failed: %v", err)
			} else {
				ok = true
			}
		case <-timeout:
			tick.Stop()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "timeout to get channel"})
			return
		}
		if ok {
			break
		}
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-channel; ok {
			// print time and msg
			//fmt.Printf("%s: %s\n", time.Now().Format("2006-01-02 15:04:05"), msg)
			c.SSEvent("message", msg)
			return true
		}
		return false
	})
}
