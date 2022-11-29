package sse

import (
	"sync"
	"time"
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
	Key   string
	users []*User

	userMux sync.Mutex
}

type User struct {
	MessageChan chan string
	UserId      int
}

func (channel *Channel) PublishMsg(msg string) {
	wg := sync.WaitGroup{}
	wg.Add(len(channel.users) - 1)
	for i := 0; i < len(channel.users); i++ {
		go func(u *User) {
			select {
			case u.MessageChan <- msg:
			case <-time.After(1):
				channel.removeUser(u)
			}
			wg.Done()
		}(channel.users[i])

	}
}

func (channel *Channel) removeUser(user *User) {
	channel.userMux.Lock()
	defer channel.userMux.Unlock()

	for i := 0; i < len(channel.users); i++ {
		u := channel.users[i]
		if u == user {
			channel.users = append(channel.users[:i], channel.users[i+1:]...)
			i--
		}
	}
	if len(channel.users) == 0 {

	}
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
			channel.PublishMsg(pubMsg.Message)
		}
	}
}

func GetOrCreateChannel(channelKey string) *Channel {
	portalInstance := GetPortalInstance()
	portalInstance.channelMux.Lock()
	defer portalInstance.channelMux.Unlock()

	channel, ok := portalInstance.Channels[channelKey]
	if !ok {
		channel = &Channel{
			Key:   channelKey,
			users: make([]*User, 0),
		}
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
	println("init")
}

var PortalInstance *Portal
