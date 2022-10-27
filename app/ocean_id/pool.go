package oceanId

import (
	"OceanID/config"
	"context"
	"github.com/pkg/errors"
	"log"
	"sync/atomic"
	"time"
)

type idPool struct {
	size int64
	// mps Max pool size
	mps uint64
	// nps Min pool size
	nps uint64

	ctx    context.Context
	mister Mist
	// OceanID pool
	oip     chan int64
	counter chan struct{}
}

type IdPool interface {
	BulkGetID(size int64) ([]int64, error)
	GetID() (int64, error)
}

func NewOceanID(ctx context.Context) (IdPool, error) {
	args, err := config.AssertArgs(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "app.ocean_id.pool")
	}
	mister, err := NewMist(args.Get("ID_METADATA_PATH").(string))
	if err != nil {
		return nil, err
	}
	pool := &idPool{
		size:    0,
		mps:     args.Get("MAX_ID_POOL_SIZE").(uint64),
		nps:     args.Get("MIN_ID_POOL_SIZE").(uint64),
		ctx:     ctx,
		mister:  mister,
		oip:     nil,
		counter: nil,
	}
	pool.initPool()
	go pool.mdsTicker()
	go pool.poolListener()
	return pool, nil
}

func (p *idPool) GetID() (int64, error) {
	if atomic.LoadInt64(&p.size) <= 0 {
		return 0, errors.New("get id failed, id pool is empty")
	}
	atomic.AddInt64(&p.size, -1)
	p.callOIP()
	id, ok := <-p.oip
	if !ok {
		return 0, errors.New("get id failed")
	}
	return id, nil
}

func (p *idPool) BulkGetID(size int64) ([]int64, error) {
	pSize := atomic.LoadInt64(&p.size)
	switch {
	case pSize == 0:
		return nil, errors.New("get id failed, id pool is empty")
	case pSize < size:
		return nil, errors.New("request too large")
	}
	bi := make([]int64, size)
	atomic.AddInt64(&p.size, -size)
	for i := 0; i < int(size); i++ {
		bi[i] = <-p.oip
		p.callOIP()
	}
	return bi, nil
}

func (p *idPool) callOIP() {
	p.counter <- struct{}{}
}

func (p *idPool) autoFillPool() {
	fillSize := int64(p.mps)
	pSize := atomic.LoadInt64(&p.size)
	if pSize <= int64(p.nps) {
		fillSize = int64(p.mps) - pSize
	}
	addSize := 0
	for i := 0; i < int(fillSize); i++ {
		select {
		case p.oip <- p.mister.GetID():
			addSize++
		default:
			atomic.AddInt64(&p.size, int64(addSize))
			log.Printf("id pool filled, expected: %d, actual: %d", fillSize, addSize)
			return
		}
	}
	log.Printf("id pool filled! fill size: %d", addSize)
	atomic.AddInt64(&p.size, int64(addSize))
}

func (p *idPool) initPool() {
	p.oip = make(chan int64, p.mps)
	p.counter = make(chan struct{}, p.nps-p.nps/10)
	// first fill this pool
	p.autoFillPool()
}

func (p *idPool) poolListener() {
	counter := 0
	counterLimit := int(p.nps / 10)
	for {
		if counter >= counterLimit {
			// call autoFill
			p.autoFillPool()
		}
		select {
		case <-p.ctx.Done():
			return
		case <-p.counter:
			counter++
		}
	}
}

// metadata save ticker
func (p *idPool) mdsTicker() {
	ticker := time.NewTicker(time.Minute * 1)
	for {
		select {
		case <-p.ctx.Done():
			if err := p.mister.WriteMetadata(); err != nil {
				log.Println(errors.Wrap(err, "write OceanID metadata failed, context done"))
			}
			return
		case <-ticker.C:
			if err := p.mister.WriteMetadata(); err != nil {
				log.Println(errors.Wrap(err, "write OceanID metadata failed"))
			}
		}
	}
}
