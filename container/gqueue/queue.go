package gqueue

import (
	"container/list"
	"goutil/safe/gwg"
	"time"
)

type (
	Queue struct {
		list  *list.List
		wg    *gwg.WaitSwitch
		limit uint64
	}
)

func NewQueue() *Queue {
	q := &Queue{list: list.New(), wg: gwg.NewWaitSwitch()}
	q.wg.SetBlock()
	return q
}

func (q *Queue) SetLimit(n uint64) {
	q.limit = n
}

func (q *Queue) Push(v any) (oldestDeleted bool) {
	defer q.wg.SetUnBlock()

	oldestDeleted = false
	if q.limit > 0 && uint64(q.list.Len()) >= q.limit {
		q.Pop(nil)
		oldestDeleted = true
	}
	q.list.PushBack(v)
	return oldestDeleted
}

func (q *Queue) Pop(waitTimeout *time.Duration) any {
	if waitTimeout != nil {
		q.wg.Wait(waitTimeout)
	}

	var result any = nil
	if q.list.Len() > 0 {
		result = q.list.Front().Value
		q.list.Remove(q.list.Front())
		if q.list.Len() == 0 {
			q.wg.SetBlock()
		}
	}
	return result
}

func (q *Queue) Len() int {
	return q.list.Len()
}
