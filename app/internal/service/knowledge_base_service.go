package service

import (
	"context"
	"ragx/app/internal/biz"

	pb "ragx/api/gen"
)

type KnowledgeBaseService struct {
	pb.UnimplementedKnowledgeBaseServiceServer
	uc *biz.KnowledgeBaseUsecase
}

func NewKnowledgeBaseService(uc *biz.KnowledgeBaseUsecase) *KnowledgeBaseService {
	return &KnowledgeBaseService{uc: uc}
}

func (s *KnowledgeBaseService) CreateKnowledgeBase(ctx context.Context, req *pb.CreateKnowledgeBaseRequest) (*pb.IDReply, error) {
	return s.uc.Create(ctx, req)
}
func (s *KnowledgeBaseService) UpdateKnowledgeBase(ctx context.Context, req *pb.CreateKnowledgeBaseRequest) (*pb.IDReply, error) {
	return s.uc.Update(ctx, req)
}
func (s *KnowledgeBaseService) DeleteKnowledgeBase(ctx context.Context, req *pb.IDReply) (*pb.IDReply, error) {
	if err := s.uc.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &pb.IDReply{Id: req.Id}, nil
}
func (s *KnowledgeBaseService) GetKnowledgeBase(ctx context.Context, req *pb.IDReply) (*pb.KnowledgeBase, error) {
	return &pb.KnowledgeBase{}, nil
}
func (s *KnowledgeBaseService) ListKnowledgeBase(ctx context.Context, req *pb.ListKnowledgeBaseRequest) (*pb.ListKnowledgeBaseReply, error) {
	return s.uc.List(ctx, req)
}
