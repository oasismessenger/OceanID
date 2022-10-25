package oceanID

import (
	"OceanID/config"
	"OceanID/utils"
	"bytes"
	"context"
	"github.com/pkg/errors"
	"log"
	"sync/atomic"
	"time"
)

const (
	SaltBit          uint  = 7
	SaltShift        uint  = 4
	IncrShift        uint  = SaltBit + SaltShift
	DefaultIncrValue int64 = 100
	MaxUint8         int64 = 1<<8 - 1
)

type OceanID struct {
	incr    int64
	maxPool uint64
	minPool uint64

	ctx context.Context
	mdp string
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
		maxPool: args.Get("MAX_ID_POOL_SIZE").(uint64),
		minPool: args.Get("MIN_ID_POOL_SIZE").(uint64),
		ctx:     ctx,
		mdp:     mdPath,
	}
	go oci.idTicker()
	return oci, nil
}

func (i *OceanID) GetId() int64 {
	atomic.AddInt64(&i.incr, 1)
	a, b := utils.Int64(MaxUint8), utils.Int64(MaxUint8)
	return (atomic.LoadInt64(&i.incr) << IncrShift) | (a << SaltShift) | b
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
