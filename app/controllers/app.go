package controllers

import (
	"container/list"
	"errors"
	"github.com/MonkeyBuisness/screencast/app/service"
	"github.com/revel/revel"
	"log"
)

var (
	subscribers *list.List
	listener    ServiceListener
)

type App struct {
	*revel.Controller
}

type ServiceListener struct{}

type Subscriber struct {
	socket  revel.ServerWebSocket
	receive chan []byte
	send    chan []byte
	close   chan error
}

type Frame struct {
	Data []byte `json:"image"`
}

func (c App) Ping() revel.Result {
	return c.RenderText("PONG")
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Mirror(ws revel.ServerWebSocket) revel.Result {
	// add new subscriber
	subscriber := NewSubscriber(ws)
	subscribers.PushBack(subscriber)

	defer func(subscriber *Subscriber) {
		subscriber.Close()
		for e := subscribers.Front(); e != nil; e = e.Next() {
			if e.Value.(*Subscriber) == subscriber {
				subscribers.Remove(e)
				break
			}
		}
	}(subscriber)

	// start receive listener
	go func(s *Subscriber) {
		var msg string
		for {
			if err := s.socket.MessageReceiveJSON(&msg); err != nil {
				s.close <- err
				return
			}

			s.receive <- []byte(msg)
		}
	}(subscriber)

	// start subscriber listener
	if err := subscriber.Listen(); err != nil {
		log.Println(err)
	}

	return nil
}

func (s Subscriber) Listen() error {
	for {
		select {
		case data, ok := <-s.send:
			if ok {
				if err := s.socket.MessageSend(data); err != nil {
					return err
				}
			} else {
				return errors.New("send channel closed")
			}
		case data, ok := <-s.receive:
			if ok {
				// TODO
				log.Println(data)
			} else {
				return errors.New("receive channel closed")
			}
		case err := <-s.close:
			return err
		}
	}
}

func (s Subscriber) Close() {
	close(s.receive)
	close(s.send)
	close(s.close)
}

func (s Subscriber) Send(data []byte) {
	s.send <- data
}

func (sl ServiceListener) NewScreenshot(data []byte) {
	sendScreenshot(data)
}

func NewSubscriber(ws revel.ServerWebSocket) *Subscriber {
	return &Subscriber{
		socket:  ws,
		receive: make(chan []byte, 1),
		send:    make(chan []byte, 1),
		close:   make(chan error),
	}
}

func sendScreenshot(data []byte) {
	for e := subscribers.Front(); e != nil; e = e.Next() {
		if subscriber, ok := e.Value.(*Subscriber); ok {
			subscriber.Send(data)
		}
	}
}

func init() {
	// init subscribers list
	subscribers = list.New()

	// init screencast service listener
	listener = ServiceListener{}
	service.SetListener(listener)
}
