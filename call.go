package bunch

import "sync/atomic"

func (b *Bunch) NewCalls() *callList {
	list := &callList{
		calls: make([]call, 0, 2),
		b:     b,
	}
	return list
}

type call struct {
	callf  func() (interface{}, error)
	afterf func(interface{}, error)
}

type callList struct {
	calls []call
	b     *Bunch

	callRet interface{}
	err     error

	nextCallIdx int

	canceled int32
}

func (list *callList) callNext() {
	call := list.calls[list.nextCallIdx]
	list.nextCallIdx++
	list.callRet, list.err = call.callf()
	list.b.afterQ <- func() {
		call.afterf(list.callRet, list.err)
		if atomic.LoadInt32(&list.canceled) != 0 {
			return
		}
		if list.nextCallIdx < len(list.calls) {
			list.b.submit(list.callNext)
		}
	}
}

func (list *callList) Call(f func() (interface{}, error)) setAfter {
	ca := call{}
	ca.callf = f
	list.calls = append(list.calls, ca)
	return setAfter{
		callList: list,
		index:    len(list.calls) - 1,
	}
}

func (list *callList) Commit() {
	if list.nextCallIdx < len(list.calls) {
		list.b.submit(list.callNext)
	}
}

func (list *callList) Cancel() {
	atomic.StoreInt32(&list.canceled, 1)
}

type setAfter struct {
	callList *callList
	index    int
}

func (ca setAfter) After(f func(interface{}, error)) *callList {
	ca.callList.calls[ca.index].afterf = f
	return ca.callList
}
