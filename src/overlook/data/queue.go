package data

import "log"

type Queue struct {
	element []interface{}
}

func NewQueue(other *Queue) *Queue {
	queue := &Queue{
		element: make([]interface{}, 0),
	}

	if other != nil {
		queue.cloneFrom(other)
	}
	return queue
}

func (self *Queue) cloneFrom(other *Queue) {
	for i := 0; i < len(other.element); i++ {
		self.element = append(self.element, other.element[i])
	}
}

func (self *Queue) Push(v interface{}) *Queue {
	self.element = append(self.element, v)
	return self
}

func (self *Queue) Pop() interface{} {
	if self.element != nil {
		v := self.element[0]
		newElem := make([]interface{}, 0)
		newElem = append(newElem, self.element[1:]...)
		self.element = newElem
		return v
	}
	log.Fatal("element is nil , use : func NewQueue(other*Queue) *Queue to create a new queue.")
	return nil
}

func (self *Queue) HasMore() bool {
	return len(self.element) > 0
}

func (self *Queue) Clear() {
	self.element = make([]interface{}, 0)
}
