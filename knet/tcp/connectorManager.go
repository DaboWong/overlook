package tcp

import (
	"log"
	"sync"
)

type ConnectorManager struct {
	reg   chan *Connector
	unreg chan *Connector
	conns map[string]*Connector
	mutex sync.Mutex
}

var ConnMgr = ConnectorManager{
	reg:   make(chan *Connector, 1),
	unreg: make(chan *Connector, 1),
	conns: make(map[string]*Connector, 0),
	mutex: sync.Mutex{},
}

func (self *ConnectorManager) add(con *Connector) {
	self.reg <- con
}

func (self *ConnectorManager) remove(con *Connector) {
	self.unreg <- con
}

func (self *ConnectorManager) Run() {
	for {
		select {
		case con, ok := <-self.reg:
			if ok {
				self.mutex.Lock()
				log.Printf("add conn, name :%s, %s", con.Name(), con.conn.LocalAddr().String())
				if _, ok := self.conns[con.Name()]; ok {
					log.Fatalln("register error, connection exist!")
				}
				self.conns[con.Name()] = con

				//
				for name, c := range self.conns {
					log.Printf("now conns: %s %s", name, c.conn.LocalAddr().String())
				}
				//
				self.mutex.Unlock()
			}
		case con, ok := <-self.unreg:
			if ok && con != nil {
				self.mutex.Lock()
				log.Printf("remove conn, name :%s, %s", con.Name(), con.conn.LocalAddr().String())
				delete(self.conns, con.Name())

				//
				for name, c := range self.conns {
					log.Printf("now conns: %s %s", name, c.conn.LocalAddr().String())
				}
				//
				self.mutex.Unlock()
			}
		}
	}
}

func (self *ConnectorManager) Get(name string) (*Connector, bool) {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	v, ok := self.conns[name]
	log.Printf("Get Connector %s %+v", name, v.conn)
	return v, ok
}

func init() {
	log.Printf("conn manager ....start")
	go ConnMgr.Run()
}
