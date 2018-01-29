package ds

import (
	"log"
)

type Event struct {
	ID    string
	Value interface{}
}

type Notifier interface {
	Add(name string, f func(evt *Event, v ...interface{}))
	Notify(name string, evt *Event, v ...interface{}) bool
}

type delegate struct {
	Callback func(evt *Event, v ...interface{})
}

func (self *delegate) Call(evt *Event, v ...interface{}) {
	if self.Callback != nil {
		self.Callback(evt, v...)
	}
}

type Delegates struct {
	callbacks map[*delegate]*delegate
}

func (self *Delegates) Add(f func(evt *Event, v ...interface{})) {
	del := &delegate{Callback: f}

	if _, ok := self.callbacks[del]; ok {
		return
	}
	self.callbacks[del] = del
}

func (self *Delegates) Invoke(evt *Event, v ...interface{}) {
	if self.callbacks != nil {
		for _, callback := range self.callbacks {
			if callback != nil {
				callback.Call(evt, v...)
			}
		}
	}
}

type EventHandler struct {
	delegates map[string]*Delegates
}

func NewEventHandler() Notifier {
	return &EventHandler{
		delegates: make(map[string]*Delegates, 0),
	}
}

func (self *EventHandler) Add(name string, f func(evt *Event, v ...interface{})) {
	//log.Printf("add event: %s", name)
	v, ok := self.delegates[name]
	if ok {
		v.Add(f)
	} else {
		d := &Delegates{callbacks: make(map[*delegate]*delegate)}
		d.Add(f)
		self.delegates[name] = d
	}
}

func (self *EventHandler) Notify(name string, evt *Event, v ...interface{}) bool {
	log.Printf("Event Notify：EventID:[%s]", name)
	if d, ok := self.delegates[name]; ok {
		d.Invoke(evt, v...)
		return true
	}
	return false
}
