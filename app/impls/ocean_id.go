package impls

import (
	idService "OceanID/schemes/id_service"
	"context"
)

type OceanID struct {
	idService.UnimplementedOceanIDServer
}

func (o OceanID) GenerateID(_ context.Context, request *idService.IDRequest) (*idService.IDReply, error) {

	return nil, nil
}

func (o OceanID) BulkGenerateID(_ context.Context, request *idService.IDBulkRequest) (*idService.IDBulkReply, error) {

	return nil, nil
}

func (o OceanID) mustEmbedUnimplementedOceanIDServer() {}
