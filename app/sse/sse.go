package sse_server

import (
	"sse_demo/util"
	"sync"
	"time"
)

const (
	BroadcastChannelKey = "broadcast_sse_channel"
)

type PubMessage struct {
	ChannelKey string
	Message    string
}

type Portal struct {
	PubMessageChan chan PubMessage
	Channels       map[string]*Channel

	channelMux sync.Mutex
}

type Channel struct {
	Key     string
	clients map[string]*Client

	ClientMux sync.Mutex
	portal    *Portal
}

type Client struct {
	MessageChan chan string
	ClientId    string

	channel *Channel
}

func (channel *Channel) PublishMsg(msg string) {

	if len(channel.clients) == 0 {
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(len(channel.clients))

	for _, v := range channel.clients {
		go func(c *Client) {
			select {
			case c.MessageChan <- msg:
				util.MyLogger.Info("sent msg: " + msg)
			case <-time.After(time.Second * 5):
				channel.removeClient(c)
			}
			wg.Done()
		}(v)
	}
}

func (channel *Channel) removeClient(client *Client) {
	channel.ClientMux.Lock()
	defer channel.ClientMux.Unlock()

	delete(channel.clients, client.ClientId)
	close(client.MessageChan)
	// do not delete
	//if len(channel.clients) == 0 {
	//	channel.portal.removeChannel(channel)
	//}
}

func (channel *Channel) Subscribe(client *Client) {
	channel.ClientMux.Lock()
	defer channel.ClientMux.Unlock()

	channel.clients[client.ClientId] = client
}

func (p *Portal) removeChannel(channel *Channel) {
	p.channelMux.Lock()
	defer p.channelMux.Unlock()

	delete(p.Channels, channel.Key)
}

func (p *Portal) loop() {
	for {
		select {
		case pubMsg := <-p.PubMessageChan:
			channelKey := pubMsg.ChannelKey
			channel, ok := p.Channels[channelKey]
			if !ok {
				break
			}
			go channel.PublishMsg(pubMsg.Message)
		}
	}
}

func GetOrCreateChannel(channelKey string) *Channel {
	portalInstance := GetPortalInstance()
	portalInstance.channelMux.Lock()
	defer portalInstance.channelMux.Unlock()

	channel, ok := portalInstance.Channels[channelKey]
	if !ok {
		println("init channel " + channelKey)
		channel = &Channel{
			Key:     channelKey,
			clients: make(map[string]*Client),
			portal:  portalInstance,
		}
		portalInstance.Channels[channelKey] = channel
	}
	return channel
}

func GetPortalInstance() *Portal {

	return PortalInstance
}

func init() {
	PortalInstance = &Portal{
		PubMessageChan: make(chan PubMessage),
		Channels:       make(map[string]*Channel),
	}
	go PortalInstance.loop()
	println("init portalInstance")
}

var PortalInstance *Portal
