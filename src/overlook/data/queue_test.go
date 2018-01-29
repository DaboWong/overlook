package data

import (
	"log"
	"testing"
)

type play struct {
	Idx int
}

func Test_Queue(t *testing.T) {

	q := NewQueue(nil)
	go q.Push(&play{0})
	go q.Push(&play{1})
	go q.Push(&play{2})
	go q.Push(&play{3})
	x := NewQueue(q)

	for q.HasMore() {
		p := q.Pop()
		log.Printf("p:%d", p.(*play).Idx)
	}

	log.Println("-------------------------")

	for x.HasMore() {
		p := x.Pop()
		log.Printf("p:%d", p.(*play).Idx)
	}

	x.Clear()

	log.Println("-------------------------")

	log.Println(len(x.element) == 0)
	log.Println(x.element != nil)

}
