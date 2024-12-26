package bunch

import (
	"container/list"
	"time"

	ants "github.com/panjf2000/ants/v2"
	"github.com/tidwall/spinlock"
)

type Bunch struct {
	config *config
	pool   *ants.Pool

	afterQ chan func()

	waitSubmits      *list.List
	waitSubmitMtx    *spinlock.Locker
	waitSubmitSignal chan byte
}

type config struct {
	poolSize int
}

type Option func(*config)

var (
	WithPoolSize = func(size int) Option {
		return func(c *config) {
			c.poolSize = size
		}
	}
)

func New(opts ...Option) *Bunch {
	config := &config{
		poolSize: 1000,
	}
	for _, opt := range opts {
		opt(config)
	}
	pool, err := ants.NewPool(config.poolSize, ants.WithNonblocking(true))
	if err != nil {
		panic(err)
	}

	b := &Bunch{
		config:           config,
		pool:             pool,
		afterQ:           make(chan func(), config.poolSize),
		waitSubmits:      list.New(),
		waitSubmitMtx:    &spinlock.Locker{},
		waitSubmitSignal: make(chan byte, 12),
	}

	go b.run()

	return b
}

func (b *Bunch) run() {
	ticker := time.NewTicker(time.Millisecond * 30)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			b.resubmit()
		case <-b.waitSubmitSignal:
			b.resubmit()
		}
	}
}

func (b *Bunch) AfterQ() <-chan func() {
	return b.afterQ
}

func (b *Bunch) notifyResubmit() {
	select {
	case b.waitSubmitSignal <- 0:
	default:
	}
}

func (b *Bunch) submit(f func()) error {
	err := b.pool.Submit(func() {
		f()
		b.notifyResubmit()
	})

	if err == nil {
		return nil
	}
	if err == ants.ErrPoolOverload {
		b.waitSubmitMtx.Lock()
		b.waitSubmits.PushBack(f)
		b.waitSubmitMtx.Unlock()
		return nil
	}
	return err
}

func (b *Bunch) resubmit() {
	if b.waitSubmits.Len() == 0 {
		return
	}

	b.waitSubmitMtx.Lock()
	defer b.waitSubmitMtx.Unlock()

	for elem := b.waitSubmits.Front(); elem != nil; {
		f := elem.Value.(func())
		if b.pool.Submit(func() {
			f()
			b.notifyResubmit()
		}) == nil {
			next := elem.Next()
			b.waitSubmits.Remove(elem)
			elem = next
		} else {
			break
		}
	}
}
