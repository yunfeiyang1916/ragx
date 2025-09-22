package service

import (
	"context"

	pb "ragx/api/gen"
)

type ChatService struct {
	pb.UnimplementedChatServiceServer
}

func NewChatService() *ChatService {
	return &ChatService{}
}

func (s *ChatService) Chat(ctx context.Context, req *pb.ChatRequest) (*pb.ChatReply, error) {
	return &pb.ChatReply{}, nil
}
