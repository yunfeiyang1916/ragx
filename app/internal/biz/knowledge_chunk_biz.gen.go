package biz

import (
	"context"
	"ragx/app/internal/biz/entity"
	"ragx/app/internal/biz/query"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type KnowledgeChunkRepo interface {
	Query() *query.Query
	// 批量创建，支持事务
	BatchCreate(context.Context, []*entity.KnowledgeChunk, ...*query.Query) ([]*entity.KnowledgeChunk, error)
	// 创建，支持事务
	Create(context.Context, *entity.KnowledgeChunk, ...*query.Query) (*entity.KnowledgeChunk, error)
	Update(context.Context, *entity.KnowledgeChunk, ...field.Expr) (int64, error)
	UpdateWithTx(context.Context, *query.Query, *entity.KnowledgeChunk, ...field.Expr) (int64, error)
	// 保存全部字段，支持事务
	Save(context.Context, *entity.KnowledgeChunk, ...*query.Query) (int64, error)
	// 删除，支持事务
	Delete(context.Context, int64, ...*query.Query) (int64, error)
	DeleteByConditions(context.Context, ...gen.Condition) (int64, error)
	DeleteByConditionsWithTx(context.Context, *query.Query, ...gen.Condition) (int64, error)
	Get(context.Context, int64, ...field.RelationField) (*entity.KnowledgeChunk, error)
	GetByConditions(context.Context, ...gen.Condition) (*entity.KnowledgeChunk, error)
	// 支持预加载
	GetByConditionsWithPreload(context.Context, []field.RelationField, ...gen.Condition) (*entity.KnowledgeChunk, error)
	List(context.Context, *entity.PageAndOrder, ...gen.Condition) ([]*entity.KnowledgeChunk, int64, error)
	// 只需要列表，不需要总数
	ListWithoutCount(context.Context, *entity.PageAndOrder, ...gen.Condition) ([]*entity.KnowledgeChunk, error)
	ListAll(context.Context, ...gen.Condition) ([]*entity.KnowledgeChunk, error)
	// 支持预加载
	ListAllWithPreload(context.Context, []field.RelationField, ...gen.Condition) ([]*entity.KnowledgeChunk, error)
	Count(context.Context, ...gen.Condition) (int64, error)
}

type KnowledgeChunkUsecase struct {
	repo KnowledgeChunkRepo
	log  *log.Helper
}

func NewKnowledgeChunkUsecase(repo KnowledgeChunkRepo, logger log.Logger) *KnowledgeChunkUsecase {
	return &KnowledgeChunkUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *KnowledgeChunkUsecase) Create(ctx context.Context, obj *entity.KnowledgeChunk) (*entity.KnowledgeChunk, error) {
	e, err := uc.repo.Create(ctx, obj)
	if err != nil {
		uc.log.Errorf("%+v", err)
		return nil, err
	}
	return e, nil
}

func (uc *KnowledgeChunkUsecase) Update(ctx context.Context, obj *entity.KnowledgeChunk) (*entity.KnowledgeChunk, error) {
	_, err := uc.repo.Save(ctx, obj)
	if err != nil {
		uc.log.Errorf("%+v", err)
		return nil, err
	}
	return obj, nil
}

func (uc *KnowledgeChunkUsecase) Delete(ctx context.Context, id int64) error {
	_, err := uc.repo.Delete(ctx, id)
	if err != nil {
		uc.log.Errorf("%+v", err)
		return err
	}
	return nil
}

func (uc *KnowledgeChunkUsecase) Get(ctx context.Context, id int64) (*entity.KnowledgeChunk, error) {
	e, err := uc.repo.Get(ctx, id)
	if err != nil {
		uc.log.Errorf("%+v", err)
		return nil, err
	}
	return e, nil
}

func (uc *KnowledgeChunkUsecase) List(ctx context.Context, page *entity.PageAndOrder) ([]*entity.KnowledgeChunk, int64, error) {
	arr, count, err := uc.repo.List(ctx, page)
	if err != nil {
		uc.log.Errorf("%+v", err)
		return nil, 0, err
	}
	return arr, count, nil
}

func (uc *KnowledgeChunkUsecase) ListAll(ctx context.Context) ([]*entity.KnowledgeChunk, error) {
	arr, err := uc.repo.ListAll(ctx)
	if err != nil {
		uc.log.Errorf("%+v", err)
		return nil, err
	}
	return arr, nil
}
