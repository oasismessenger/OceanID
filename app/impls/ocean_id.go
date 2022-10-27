package impls

import (
	"context"
	"time"

	"OceanID/app/ocean_id"
	"OceanID/schemes/id_service"
)

type OceanID struct {
	oceanId.IdPool
	idService.UnimplementedOceanIDServer
}

func (o *OceanID) SetOI(oi oceanId.IdPool) {
	o.IdPool = oi
}

func (o *OceanID) GenerateID(_ context.Context, request *idService.IDRequest) (*idService.IDReply, error) {
	id, err := o.GetID()
	if err != nil {
		return nil, err
	}
	return &idService.IDReply{
		Id:        id,
		Timestamp: uint64(time.Now().UnixNano()),
		ReplyId:   request.RequestId,
	}, nil
}

func (o *OceanID) BulkGenerateID(_ context.Context, request *idService.IDBulkRequest) (*idService.IDBulkReply, error) {
	ids, err := o.BulkGetID(int64(request.GetBulkSize()))
	if err != nil {
		return nil, err
	}
	return &idService.IDBulkReply{
		Ids:       ids,
		Timestamp: uint64(time.Now().UnixNano()),
		ReplyId:   request.RequestId,
		Size:      uint32(len(ids)),
	}, nil
}

func (o *OceanID) mustEmbedUnimplementedOceanIDServer() {}
