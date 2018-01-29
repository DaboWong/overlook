package tcp

import (
	"log"
	"net"
	"os"
	"overlook/codec"
	"overlook/ds"
)

type Listener struct {
	ls   net.Listener
	host string

	connector map[*Connector]*Connector

	register chan *Connector
	ungister chan *Connector

	notifier ds.Notifier
}

func (self *Listener) StartAccept(callback ...func(notifier ds.Notifier, conn *Connector)) {

	for {
		conn, err := self.ls.Accept()
		if err != nil {
			log.Println("accept error:", err.Error())
			continue
		}

		connector := &Connector{
			conn:           conn,
			notifier:       ds.NewEventHandler(),
			listener:       self,
			SyncMode:       false,
			write:          make(chan []byte, 256),
			decoder:        codec.NewProtoBufDecoder(),
			name:           "game",
			encoder:        &codec.TKEncoder{},
			IDataContainer: ds.NewDataContainer(),
			exit:           make(chan int, 1),
		}

		ConnMgr.add(connector)

		for i := 0; i < len(callback); i++ {
			if callback[i] != nil {
				callback[i](connector.notifier, connector)
			}
		}

		self.register <- connector

		go connector.Run()
	}
}

func (self *Listener) Run() {
	for {
		select {
		case m, ok := <-self.register:
			if ok {
				self.connector[m] = m
			}
		case m, ok := <-self.ungister:
			if ok {
				delete(self.connector, m)
			}

		}
	}
}

func StartListener(host string, notifier ds.Notifier, callback ...func(notifier ds.Notifier, conn *Connector)) *Listener {

	log.Println("start listen at: ", host)

	l, err := net.Listen("tcp", host)

	if err != nil {
		log.Fatalln("listen error, please check Host is config ok!", host)
		os.Exit(1)
	}

	listener := &Listener{
		ls:        l,
		host:      host,
		connector: make(map[*Connector]*Connector, 0),
		register:  make(chan *Connector, 1),
		ungister:  make(chan *Connector, 1),
		notifier:  notifier,
	}

	go listener.StartAccept(callback...)

	return listener
}
