package impls

import (
	"OceanID/app/ocean_id"
	"context"
	"time"

	"OceanID/schemes/id_service"
)

type OceanIDGrpc struct {
	oceanID.IDPool
	idService.UnimplementedOceanIDServer
}

func (o *OceanIDGrpc) SetOI(oi oceanID.IDPool) {
	o.IDPool = oi
}

func (o *OceanIDGrpc) GenerateID(_ context.Context, request *idService.IDRequest) (*idService.IDReply, error) {
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

func (o *OceanIDGrpc) BulkGenerateID(_ context.Context, request *idService.IDBulkRequest) (*idService.IDBulkReply, error) {
	ids, err := o.BulkGetID(int64(request.BulkSize))
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

func (o *OceanIDGrpc) mustEmbedUnimplementedOceanIDServer() {}
