package impls

import (
	"OceanID/app/ocean_id"
	"OceanID/schemes/id_service"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
	"time"
)

func RequestDecoder[T any](r io.Reader) (*T, error) {
	var request T
	if err := jsoniter.ConfigFastest.NewDecoder(r).Decode(&request); err != nil {
		return nil, err
	}
	return &request, nil
}

type OceanIDHttp struct {
	oceanID.IDPool
}

func (h *OceanIDHttp) SetOI(oi oceanID.IDPool) {
	h.IDPool = oi
}

func (h *OceanIDHttp) GenerateID(w http.ResponseWriter, r *http.Request) error {
	request, err := RequestDecoder[idService.IDRequest](r.Body)
	if err != nil {
		return err
	}
	id, err := h.GetID()
	if err != nil {
		return err
	}
	jsoniter.ConfigFastest.NewEncoder(w).Encode(&idService.IDReply{
		Id:        id,
		Timestamp: uint64(time.Now().UnixNano()),
		ReplyId:   request.RequestId,
	})
	return nil
}

func (h *OceanIDHttp) BulkGenerateID(w http.ResponseWriter, r *http.Request) error {
	request, err := RequestDecoder[idService.IDBulkRequest](r.Body)
	if err != nil {
		return err
	}
	ids, err := h.BulkGetID(int64(request.BulkSize))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	jsoniter.ConfigFastest.NewEncoder(w).Encode(&idService.IDBulkReply{
		Ids:       ids,
		Timestamp: uint64(time.Now().UnixNano()),
		ReplyId:   request.RequestId,
		Size:      uint32(len(ids)),
	})

	return nil
}
