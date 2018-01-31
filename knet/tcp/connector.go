package tcp

import (
	"bytes"
	"encoding/binary"
	"knet/codec"
	"knet/ds"
	log "log"
	"net"
	codec2 "overlook/codec"
	"tkbase/msg"
	"sync"
	"sync/atomic"
	"time"
)

const (
	max_pack_size int = 65533

	//心跳间隔（秒）
	heart_beat_interval time.Duration = 30
)

var headLen = binary.Size(&msg.TKHeader{})

type Connector struct {
	ds.IDataContainer

	conn net.Conn

	write chan []byte

	notifier ds.Notifier

	listener *Listener

	decoder codec.Decoder

	SyncMode bool

	name string

	encoder codec.Encoder

	recChan chan *ds.Event

	recChanGuard sync.Mutex

	Component []func(n ds.Notifier, con *Connector)

	Host string

	autoId int32

	exit chan int
}

func (self *Connector) Read() {

	result := true

	if self.SyncMode {
		result = self.read()
	} else {
		for result {
			result = self.read()
		}
	}

	if !result {
		self.Close()
	}
}

func (self *Connector) read() bool {
	head := make([]byte, headLen)
	//read head
	n, err := self.conn.Read(head)
	if err != nil {
		log.Printf("[%s]read socket error read count:%d,  error: %s ", self.Name(), n, err.Error())
		return false
	}
	reader := bytes.NewBuffer(head)
	header := &msg.TKHeader{}
	binary.Read(reader, binary.LittleEndian, header)

	var body []byte = nil

	//read body
	if header.Length > 0 {
		body = make([]byte, header.Length)
		_, err = self.conn.Read(body)
		if err != nil {
			log.Printf("[%s]read socket error, service exit! %s", self.Name(), err.Error())
			return false
		}
	}

	//	//decode package
	//	if self.decoder != nil {
	//		e, r := self.decoder.Decode(body, header)
	//		if r != nil {
	//			log.Printf("[%s]error decode msg %s %+v", self.Name(), r.Error(), header)
	//			return false
	//		} else if self.notifier != nil {
	//			log.Printf("[RECV][%s]<-[SVR] id: %d, type: %T, header: %+v, msg: %+v", self.Name(), header.Type, e, header, e)
	//			self.notifier.Notify(e.ID, e, header.Serial)
	//
	//			if self.recChan != nil {
	//				self.recChan <- e
	//			}
	//			return true
	//		}
	//	}

	self.decodeAndNotify(body, header)

	//return false
	return true
}

func (self *Connector) decodeAndNotify(body []byte, header *msg.TKHeader) bool {
	//decode package
	if self.decoder != nil {
		e, r := self.decoder.Decode(body, header)
		if r != nil {
			log.Printf("[%s]error decode msg %s %+v", self.Name(), r.Error(), header)
			return false
		} else if self.notifier != nil {
			log.Printf("[RECV][%s]<-[SVR] id: %d, type: %T, header: %+v, msg: %+v", self.Name(), header.Type, e, header, e)

			if ok := self.notifier.Notify(e.ID, e, header.Serial); !ok {
				self.notifier.Notify(msg.MatchToWebNtfID, e)
			}

			if self.recChan != nil {
				self.recChan <- e
			}
			return true
		}
	}

	return false
}

func (self *Connector) Write() {

	for {
		select {
		case b, ok := <-self.write:
			if ok {
				_, err := self.conn.Write(b)
				if err != nil {
					log.Printf("[%s]write error: :%s", self.Name(), err.Error())
					self.Close()
					return
				}
			}
		}
	}
}

func (self *Connector) Send(id int32, serial int32, v interface{}, ch chan *ds.Event) {
	b, err := self.encoder.Encode(id, v, serial)
	if err != nil {
		log.Fatalln("send message error, id:", id, "error:", err.Error())
	} else {
		log.Printf("[SEND][%s]->[SVR] id: %d, type: %T %+v", self.Name(), id, v, v)
		self.send(b)
		if self.SyncMode {
			self.setRecChan(ch)
			self.Read()
		}
	}
}

func (self *Connector) send(b []byte) {
	if self.SyncMode {
		_, err := self.conn.Write(b)
		if err != nil {
			log.Printf("send error: :%s", err.Error())
			self.Close()
		}
	} else {
		self.write <- b
	}
}

func (self *Connector) Run() {

	if self.notifier != nil {
		self.notifier.Notify("Connected", &ds.Event{
			ID:    "Connected",
			Value: self,
		})
	}

	if !self.SyncMode {
		//write
		go self.Write()

		//read
		go self.Read()
	}

	//心跳
	go self.watcher()
}

func (self *Connector) Close() {
	self.conn.Close()
	self.exit <- 1

	ConnMgr.remove(self)

	self.notifier.Notify("OnClosed", &ds.Event{
		ID:    "OnClosed",
		Value: self,
	})

	log.Println("conn closed, name:", self.Name())
}

func (self *Connector) watcher() {

	for {
		select {
		case <-time.After(heart_beat_interval * time.Second):
			//log.Printf("name :%s heart beat %s", self.Name(), time.Now().Format("2006-01-02 15:04:05 -0700 MST"))
			self.notifier.Notify("HeartBeat", nil)
		case <-self.exit:
			return
		}
	}
}

func (self *Connector) setRecChan(ch chan *ds.Event) {
	self.recChan = ch
}

var autoId int32

func StartConnector(host string, name string, notifier ds.Notifier, syncMode bool, callback ...func(n ds.Notifier, con *Connector)) {

	tcpAddr, err := net.ResolveTCPAddr("tcp4", host)
	if err != nil {
		log.Fatalln("Fatal error: %s", err.Error())
	}

	con, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatalln("Fatal error: %s", err.Error())
	}

	log.Printf("connect to %s done! ", host)

	connector := &Connector{
		conn:           con,
		write:          make(chan []byte, 256),
		notifier:       notifier,
		SyncMode:       syncMode,
		decoder:        codec2.NewProtoBufDecoder(),
		name:           name,
		encoder:        &codec2.TKEncoder{},
		IDataContainer: ds.NewDataContainer(),
		Component:      make([]func(n ds.Notifier, con *Connector), 0),
		Host:           host,
		autoId:         autoId,
		exit:           make(chan int, 1),
	}

	autoId = atomic.AddInt32(&autoId, 1)

	connector.Component = append(connector.Component, callback...)

	if callback != nil {
		for i := 0; i < len(callback); i++ {
			if v := callback[i]; v != nil {
				v(notifier, connector)
			}
		}
	}

	ConnMgr.add(connector)

	connector.Run()
}

func (self *Connector) GetNotifier() ds.Notifier {
	return self.notifier
}

func (self *Connector) Name() string {
	return self.name
}
