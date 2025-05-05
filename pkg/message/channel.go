package message

import (
	"fmt"
	"sync"
)

var once sync.Once
var _ Manager = (*ChannelManager)(nil)

type ChannelManager struct {
	channels map[string]chan string
	mu       sync.RWMutex
}

var channelManager *ChannelManager

func NewChannelManager() *ChannelManager {
	once.Do(func() {
		channelManager = &ChannelManager{
			channels: make(map[string]chan string),
			mu:       sync.RWMutex{},
		}
	})
	return channelManager
}

func (c *ChannelManager) CreateChannel(name string) (chan string, error) {
	c.mu.Lock()
	if _, exists := c.channels[name]; exists {
		return nil, fmt.Errorf("channel %s already exists", name)
	}
	ch := make(chan string, 100) // 设置合理的缓冲区大小
	c.channels[name] = ch
	c.mu.Unlock()
	return ch, nil
}

func (c *ChannelManager) RemoveChannel(name string) error {
	c.mu.Lock()
	ch, exists := c.channels[name]
	if !exists {
		return fmt.Errorf("channel %s not found", name)
	}
	close(ch)
	delete(c.channels, name)
	c.mu.Unlock()
	return nil
}

func (c *ChannelManager) SendMessage(channel string, message string) error {
	ch, exists := c.channels[channel]
	if !exists {
		return fmt.Errorf("channel %s not found", channel)
	}
	ch <- message
	return nil
}

func (c *ChannelManager) GetChannel(channel string) (chan string, error) {
	ch, exists := c.channels[channel]
	if !exists {
		return nil, fmt.Errorf("channel %s not found", channel)
	}
	return ch, nil
}
