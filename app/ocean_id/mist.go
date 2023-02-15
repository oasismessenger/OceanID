package oceanID

import (
	"bytes"
	"log"
	"sync/atomic"

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

type mist struct {
	incr   int64
	mdPath string
}

type Mist interface {
	GetID() int64
	WriteMetadata() error
}

func NewMist(mdp string) (Mist, error) {
	incr := DefaultIncrValue
	if utils.IsExist(mdp) {
		data, err := utils.ReadFile(mdp)
		if err != nil {
			return nil, errors.Wrap(err, "app.ocean_id.mist can't read OceanId metadata")
		}
		incr = utils.Bytes2Int(data)
	} else {
		log.Println("OceanID history file lost! reset to default value: 100")
	}
	return &mist{incr: incr, mdPath: mdp}, nil
}

func (m *mist) GetID() int64 {
	atomic.AddInt64(&m.incr, 1)
	a, b := utils.Int64(MaxUint8), utils.Int64(MaxUint8)
	return (atomic.LoadInt64(&m.incr) << IncrShift) | (a << SaltShift) | b
}

func (m *mist) WriteMetadata() error {
	return utils.Write(m.mdPath, bytes.NewBuffer(
		utils.Int2Bytes(
			atomic.LoadInt64(&m.incr),
		),
	))
}
