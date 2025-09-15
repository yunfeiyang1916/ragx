package service

import (
	"context"

	pb "ragx/api/gen"
)

type ChatServiceService struct {
	pb.UnimplementedChatServiceServer
}

func NewChatServiceService() *ChatServiceService {
	return &ChatServiceService{}
}

func (s *ChatServiceService) Chat(ctx context.Context, req *pb.ChatRequest) (*pb.ChatReply, error) {
	return &pb.ChatReply{}, nil
}
