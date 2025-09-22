package biz

import (
	"context"
	pb "ragx/api/gen"
	"ragx/app/internal/biz/entity"
	"ragx/app/internal/biz/query"
	"ragx/app/pkg/utils"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type KnowledgeBaseRepo interface {
	Query() *query.Query
	// 批量创建，支持事务
	BatchCreate(context.Context, []*entity.KnowledgeBase, ...*query.Query) ([]*entity.KnowledgeBase, error)
	// 创建，支持事务
	Create(context.Context, *entity.KnowledgeBase, ...*query.Query) (*entity.KnowledgeBase, error)
	Update(context.Context, *entity.KnowledgeBase, ...field.Expr) (int64, error)
	UpdateWithTx(context.Context, *query.Query, *entity.KnowledgeBase, ...field.Expr) (int64, error)
	// 保存全部字段，支持事务
	Save(context.Context, *entity.KnowledgeBase, ...*query.Query) (int64, error)
	// 删除，支持事务
	Delete(context.Context, int64, ...*query.Query) (int64, error)
	DeleteByConditions(context.Context, ...gen.Condition) (int64, error)
	DeleteByConditionsWithTx(context.Context, *query.Query, ...gen.Condition) (int64, error)
	Get(context.Context, int64, ...field.RelationField) (*entity.KnowledgeBase, error)
	GetByConditions(context.Context, ...gen.Condition) (*entity.KnowledgeBase, error)
	// 支持预加载
	GetByConditionsWithPreload(context.Context, []field.RelationField, ...gen.Condition) (*entity.KnowledgeBase, error)
	List(context.Context, *entity.PageAndOrder, ...gen.Condition) ([]*entity.KnowledgeBase, int64, error)
	// 只需要列表，不需要总数
	ListWithoutCount(context.Context, *entity.PageAndOrder, ...gen.Condition) ([]*entity.KnowledgeBase, error)
	ListAll(context.Context, ...gen.Condition) ([]*entity.KnowledgeBase, error)
	// 支持预加载
	ListAllWithPreload(context.Context, []field.RelationField, ...gen.Condition) ([]*entity.KnowledgeBase, error)
	Count(context.Context, ...gen.Condition) (int64, error)
}

type KnowledgeBaseUsecase struct {
	repo KnowledgeBaseRepo
	log  *log.Helper
}

func NewKnowledgeBaseUsecase(repo KnowledgeBaseRepo, logger log.Logger) *KnowledgeBaseUsecase {
	return &KnowledgeBaseUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *KnowledgeBaseUsecase) Create(ctx context.Context, req *pb.CreateKnowledgeBaseRequest) (*pb.IDReply, error) {
	obj := &entity.KnowledgeBase{}
	utils.Copy(obj, req)
	e, err := uc.repo.Create(ctx, obj)
	if err != nil {
		uc.log.Errorf("%+v", err)
		return nil, err
	}
	return &pb.IDReply{Id: e.ID}, nil
}

func (uc *KnowledgeBaseUsecase) Update(ctx context.Context, req *pb.CreateKnowledgeBaseRequest) (*pb.IDReply, error) {
	obj := &entity.KnowledgeBase{}
	utils.Copy(obj, req)
	_, err := uc.repo.Save(ctx, obj)
	if err != nil {
		uc.log.Errorf("%+v", err)
		return nil, err
	}
	return &pb.IDReply{Id: obj.ID}, nil
}

func (uc *KnowledgeBaseUsecase) Delete(ctx context.Context, id int64) error {
	_, err := uc.repo.Delete(ctx, id)
	if err != nil {
		uc.log.Errorf("%+v", err)
		return err
	}
	return nil
}

func (uc *KnowledgeBaseUsecase) Get(ctx context.Context, id int64) (*entity.KnowledgeBase, error) {
	e, err := uc.repo.Get(ctx, id)
	if err != nil {
		uc.log.Errorf("%+v", err)
		return nil, err
	}
	return e, nil
}

func (uc *KnowledgeBaseUsecase) List(ctx context.Context, req *pb.ListKnowledgeBaseRequest) (*pb.ListKnowledgeBaseReply, error) {
	cond := make([]gen.Condition, 0)
	if req.Name != "" {
		cond = append(cond, query.KnowledgeBase.Name.Like("%"+req.Name+"%"))
	}
	if req.Status != 0 {
		cond = append(cond, query.KnowledgeBase.Status.Eq(req.Status))
	}
	if req.Category != "" {
		cond = append(cond, query.KnowledgeBase.Category.Like("%"+req.Category+"%"))
	}
	arr, err := uc.repo.ListAll(ctx, cond...)
	if err != nil {
		uc.log.Errorf("%+v", err)
		return nil, err
	}
	var res pb.ListKnowledgeBaseReply
	utils.Copy(&res.List, arr)
	return &res, nil
}

func (uc *KnowledgeBaseUsecase) ListAll(ctx context.Context) ([]*entity.KnowledgeBase, error) {
	arr, err := uc.repo.ListAll(ctx)
	if err != nil {
		uc.log.Errorf("%+v", err)
		return nil, err
	}
	return arr, nil
}
