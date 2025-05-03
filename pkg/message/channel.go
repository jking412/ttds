package message

import (
	"fmt"
	"sync"
)

var once sync.Once
var _ Manager = (*ChannelManager)(nil)

type ChannelManager struct {
	channels map[string]chan string
}

var channelManager *ChannelManager

func NewChannelManager() *ChannelManager {
	// TODO: channel 还要考虑更多的情况，比如并发，比如多个接收者，比如前端信息接收到一半以后，又刷新怎么办
	once.Do(func() {
		channelManager = &ChannelManager{
			channels: make(map[string]chan string),
		}
	})
	return channelManager
}

func (c *ChannelManager) CreateChannel(name string) (chan string, error) {
	if _, exists := c.channels[name]; exists {
		return nil, fmt.Errorf("channel %s already exists", name)
	}
	ch := make(chan string, 100) // 设置合理的缓冲区大小
	c.channels[name] = ch
	return ch, nil
}

func (c *ChannelManager) RemoveChannel(name string) error {
	ch, exists := c.channels[name]
	if !exists {
		return fmt.Errorf("channel %s not found", name)
	}
	close(ch)
	delete(c.channels, name)
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
