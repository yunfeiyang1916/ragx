package service

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/schema"
	"github.com/go-kratos/kratos/v2/log"
	"io"
	pb "ragx/api/gen"
	"ragx/app/internal/biz"
	"ragx/app/pkg/utils"
	"time"

	"github.com/go-kratos/kratos/v2/transport/http"
)

type StreamService struct {
	chatUc *biz.ChatUsecase
	log    *log.Helper
}

func RegisterStreamServiceHTTPServer(s *http.Server, srv *StreamService) {
	r := s.Route("/")
	r.POST("/api/v1/chat/stream", srv.ChatStream())
}

func NewStreamService(chatUc *biz.ChatUsecase, logger log.Logger) *StreamService {
	return &StreamService{
		chatUc: chatUc,
		log:    log.NewHelper(logger),
	}
}

func (s *StreamService) ChatStream() func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in pb.ChatRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, pb.ChatService_ChatStream_FullMethodName)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return s.chatUc.ChatStream(ctx, req.(*pb.ChatRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		sr := out.(*schema.StreamReader[*schema.Message])
		defer sr.Close()
		httpResp := ctx.Response()
		// 设置响应头
		httpResp.Header().Set("Content-Type", "text/event-stream")
		httpResp.Header().Set("Cache-Control", "no-cache")
		httpResp.Header().Set("Connection", "keep-alive")
		httpResp.Header().Set("X-Accel-Buffering", "no") // 禁用Nginx缓冲
		httpResp.Header().Set("Access-Control-Allow-Origin", "*")
		sd := &pb.StreamData{
			Id:      utils.NewUUID(),
			Created: time.Now().Unix(),
		}
		i := 0
		for {
			message, err := sr.Recv()
			if err == io.EOF { // 流式输出结束
				break
			}
			if err != nil {
				s.log.Errorf("recv failed: %v", err)
				httpResp.Write([]byte(fmt.Sprintf("event: error\ndata: %s\n\n", err.Error())))
				break
			}
			sd.Content = message.Content
			bytes, _ := sonic.Marshal(sd)
			_, err = httpResp.Write([]byte(fmt.Sprintf("data:%s\n", string(bytes))))
			if err != nil {
				s.log.Errorf("write failed: %v", err)
			}
			s.log.Infof("message[%d]: %+v\n", i, message)
			i++
		}
		// 发送结束信号
		_, err = httpResp.Write([]byte(fmt.Sprintf("data:%s\n", "[DONE]")))
		if err != nil {
			s.log.Errorf("write failed: %v", err)
		}
		return err
	}
}
