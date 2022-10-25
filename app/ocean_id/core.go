package oceanID

import (
	"bytes"
	"context"
	"log"
	"sync/atomic"
	"time"

	"OceanID/config"
	"OceanID/utils"

	"github.com/pkg/errors"
)

const (
	SaltBit          uint  = 7
	SaltShift        uint  = 4
	IncrShift        uint  = SaltBit + SaltShift
	DefaultIncrValue int64 = 100
	MaxUint8         int64 = 1<<8 - 1
)

type OceanID struct {
	incr     int64
	poolSize int64
	maxPool  int64
	minPool  int64

	ctx context.Context
	// Ocean ID Pool
	oip     chan int64
	counter chan struct{}
	// metadata path
	mdp string
}

type OI interface {
	GetID() (int64, error)
	BulkGetID(size int64) ([]int64, error)
}

func NewOceanID(ctx context.Context) (*OceanID, error) {
	incr := DefaultIncrValue
	args, err := config.AssertArgs(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "app.ocean_id.core")
	}
	mdPath := args.Get("ID_METADATA_PATH").(string)
	if utils.IsExist(mdPath) {
		data, err := utils.ReadFile(mdPath)
		if err != nil {
			return nil, errors.Wrap(err, "app.ocean_id.core")
		}
		incr = utils.Bytes2Int(data)
	} else {
		log.Println("OceanID history file lost! reset to default value: 100")
	}
	oci := &OceanID{
		incr:    incr,
		maxPool: int64(args.Get("MAX_ID_POOL_SIZE").(uint64)),
		minPool: int64(args.Get("MIN_ID_POOL_SIZE").(uint64)),
		ctx:     ctx,
		// oip:     make(chan int64, args.Get("MAX_ID_POOL_SIZE").(uint64)),
		mdp: mdPath,
	}

	oci.newPool()
	go oci.idTicker()
	go oci.poolListener()
	return oci, nil
}

func (i *OceanID) GetID() (int64, error) {
	if atomic.LoadInt64(&i.poolSize) == 0 {
		return 0, errors.New("get id failed, id pool is empty")
	}
	atomic.AddInt64(&i.poolSize, -1)
	i.callOIP()
	return <-i.oip, nil
}

func (i *OceanID) BulkGetID(size int64) ([]int64, error) {
	ps := atomic.LoadInt64(&i.poolSize)
	switch {
	case ps == 0:
		return nil, errors.New("get id failed, id pool is empty")
	case ps < size:
		return nil, errors.New("request too large")
	}
	bi := make([]int64, size)
	atomic.AddInt64(&i.poolSize, -size)
	// TODO deadlock!!!
	for j := 0; j < int(size); j++ {
		bi[j] = <-i.oip
		i.callOIP()
	}
	return bi, nil
}

func (i *OceanID) getId() int64 {
	atomic.AddInt64(&i.incr, 1)
	a, b := utils.Int64(MaxUint8), utils.Int64(MaxUint8)
	return (atomic.LoadInt64(&i.incr) << IncrShift) | (a << SaltShift) | b
}

func (i *OceanID) fillPool() {
	bulkSize := i.maxPool
	remSize := atomic.LoadInt64(&i.poolSize)
	if remSize <= i.minPool {
		bulkSize = i.maxPool - remSize
	}
	for j := 0; j < int(bulkSize); j++ {
		i.oip <- i.getId()
	}
	atomic.AddInt64(&i.poolSize, bulkSize)
}

func (i *OceanID) newPool() {
	i.counter = make(chan struct{}, i.minPool/10)
	i.oip = make(chan int64, i.maxPool)
	// init full this pool
	i.fillPool()
}

func (i *OceanID) callOIP() {
	i.counter <- struct{}{}
}

func (i *OceanID) poolListener() {
	var (
		counter      int64 = 0
		counterLimit int64 = i.minPool / 10
	)
	for {
		if counter >= counterLimit {
			i.fillPool()
		}
		select {
		case <-i.ctx.Done():
			return
		case <-i.counter:
			counter++
		}
	}
}

func (i *OceanID) writeMD() error {
	return utils.Write(i.mdp, bytes.NewBuffer(
		utils.Int2Bytes(
			atomic.LoadInt64(&i.incr),
		),
	))
}

func (i *OceanID) idTicker() {
	ticker := time.NewTicker(time.Minute * 1)
	for {
		select {
		case <-i.ctx.Done():
			if err := i.writeMD(); err != nil {
				log.Println(errors.Wrap(err, "app.ocean_id.core write OceanID metadata failed, context done!"))
			}
			return
		case <-ticker.C:
			if err := i.writeMD(); err != nil {
				log.Println(errors.Wrap(err, "app.ocean_id.core write OceanID metadata failed"))
			}
		}
	}
}
