// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package clt

import (
	"github.com/gorilla/websocket"
	"knet/codec"
	"knet/ds"
	"log"
	"net/http"
	"overlook/data"
	"sync/atomic"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	notifier ds.Notifier

	autoId int32

	decoder codec.Decoder

	ds.IDataContainer

	Callbacks []func(notifier ds.Notifier, client *Client)
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (self *Client) readPump() {
	defer func() {
		self.hub.unregister <- self
		self.conn.Close()
	}()
	self.conn.SetReadLimit(maxMessageSize)
	self.conn.SetReadDeadline(time.Now().Add(pongWait))
	self.conn.SetPongHandler(func(string) error { self.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := self.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		//message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		//c.hub.broadcast <- message

		if self.decoder != nil {
			evt, err := self.decoder.Decode(message)
			if err != nil {
				log.Println("error decode:", err.Error())
			} else if self.notifier != nil {
				log.Printf("[RECV][%s]->[SVR] %T, %+v", "ws", evt, evt)
				self.notifier.Notify(evt.ID, evt)
			}
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("err send:%s", err.Error())
			}
			log.Printf("[SEND][%s]<-[SVR] %+v", "ws", string(message))

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (self *Client) Send(m []byte) {
	self.send <- m
}
func (self *Client) GetAutoId() int32 {
	return self.autoId
}

func (self *Client) GetNotifier() ds.Notifier {
	return self.notifier
}

var autoId int32 = 0

// serveWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, decoder codec.Decoder, callback ...func(notifier ds.Notifier, client *Client)) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		hub:            hub,
		conn:           conn,
		send:           make(chan []byte, 256),
		autoId:         autoId,
		decoder:        decoder,
		notifier:       ds.NewEventHandler(),
		IDataContainer: ds.NewDataContainer(),
		Callbacks:      make([]func(notifier ds.Notifier, client2 *Client), 0),
	}

	autoId = atomic.AddInt32(&autoId, 1)

	client.hub.register <- client
	log.Println("register web client, id:", client.autoId)

	client.AddData(data.NewWatch(client))

	client.Callbacks = append(client.Callbacks, callback...)

	for _, value := range callback {
		if value != nil {
			value(client.notifier, client)
		}
	}

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
