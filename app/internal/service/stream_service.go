package service

import (
	"context"
	pb "ragx/api/gen"

	"github.com/go-kratos/kratos/v2/transport/http"
)

type StreamService struct {
}

func RegisterStreamServiceHTTPServer(s *http.Server, srv *StreamService) {

}

func NewStreamService() *StreamService {
	return &StreamService{}
}

func (s *StreamService) ChatStream() func(ctx http.Context) error {
	return func(ctx http.Context) error {
		http.SetOperation(ctx, pb.ChatService_ChatStream_FullMethodName)

		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			//return s.patientUc.Import(ctx, file)
		})
		out, err := h(ctx, file)
		if err != nil {
			return err
		}
		reply := out.(*pb.ChatReply)
		return ctx.Result(200, reply)
	}
}
