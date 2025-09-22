package biz

import (
	"context"
	pb "ragx/api/gen"
	"ragx/app/internal/biz/entity"
	"ragx/app/internal/biz/query"
	"ragx/app/pkg/ai"

	"github.com/cloudwego/eino/components/document"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/v2/errors/gerror"
	"gorm.io/gen"
	"gorm.io/gen/field"
)

type KnowledgeDocumentRepo interface {
	Query() *query.Query
	// 批量创建，支持事务
	BatchCreate(context.Context, []*entity.KnowledgeDocument, ...*query.Query) ([]*entity.KnowledgeDocument, error)
	// 创建，支持事务
	Create(context.Context, *entity.KnowledgeDocument, ...*query.Query) (*entity.KnowledgeDocument, error)
	Update(context.Context, *entity.KnowledgeDocument, ...field.Expr) (int64, error)
	UpdateWithTx(context.Context, *query.Query, *entity.KnowledgeDocument, ...field.Expr) (int64, error)
	// 保存全部字段，支持事务
	Save(context.Context, *entity.KnowledgeDocument, ...*query.Query) (int64, error)
	// 删除，支持事务
	Delete(context.Context, int64, ...*query.Query) (int64, error)
	DeleteByConditions(context.Context, ...gen.Condition) (int64, error)
	DeleteByConditionsWithTx(context.Context, *query.Query, ...gen.Condition) (int64, error)
	Get(context.Context, int64, ...field.RelationField) (*entity.KnowledgeDocument, error)
	GetByConditions(context.Context, ...gen.Condition) (*entity.KnowledgeDocument, error)
	// 支持预加载
	GetByConditionsWithPreload(context.Context, []field.RelationField, ...gen.Condition) (*entity.KnowledgeDocument, error)
	List(context.Context, *entity.PageAndOrder, ...gen.Condition) ([]*entity.KnowledgeDocument, int64, error)
	// 只需要列表，不需要总数
	ListWithoutCount(context.Context, *entity.PageAndOrder, ...gen.Condition) ([]*entity.KnowledgeDocument, error)
	ListAll(context.Context, ...gen.Condition) ([]*entity.KnowledgeDocument, error)
	// 支持预加载
	ListAllWithPreload(context.Context, []field.RelationField, ...gen.Condition) ([]*entity.KnowledgeDocument, error)
	Count(context.Context, ...gen.Condition) (int64, error)
}

type KnowledgeDocumentUsecase struct {
	repo     KnowledgeDocumentRepo
	log      *log.Helper
	aiClient *ai.Client
}

func NewKnowledgeDocumentUsecase(repo KnowledgeDocumentRepo, logger log.Logger, aiClient *ai.Client) *KnowledgeDocumentUsecase {
	return &KnowledgeDocumentUsecase{repo: repo, log: log.NewHelper(logger), aiClient: aiClient}
}

func (uc *KnowledgeDocumentUsecase) Create(ctx context.Context, req *pb.UploadIndexerRequest) (*pb.UploadIndexerReply, error) {
	//var obj = entity.KnowledgeDocument{
	//	KnowledgeBaseName: req.KnowledgeName,
	//	FileName:          req.Uri,
	//	Status:            consts.StatusPending,
	//}
	//e, err := uc.repo.Create(ctx, &obj)
	//if err != nil {
	//	uc.log.Errorf("KnowledgeDocumentUsecase.Create err: %+v", gerror.Wrap(err, ""))
	//	return nil, err
	//}
	// 先调用加载器，加载文件内容
	docs, err := uc.aiClient.Loader.Load(ctx, document.Source{URI: req.Uri})
	if err != nil {
		uc.log.Errorf("KnowledgeDocumentUsecase.Create Load err: %+v", gerror.Wrap(err, ""))
		return nil, err
	}
	// 调用转换器，对文档进行分隔、过滤、合并
	docs, err = uc.aiClient.Transformer.Transform(ctx, docs)
	if err != nil {
		uc.log.Errorf("KnowledgeDocumentUsecase.Create Transform err: %+v", gerror.Wrap(err, ""))
		return nil, err
	}
	// 合并文档
	docs, err = ai.DocAddIDAndMerge(ctx, docs)
	if err != nil {
		uc.log.Errorf("KnowledgeDocumentUsecase.Create Transform err: %+v", gerror.Wrap(err, ""))
		return nil, err
	}
	// 调用索引器，将文档索引到向量数据库
	// 设置知识库的名称
	ctx = context.WithValue(ctx, ai.KnowledgeName, req.KnowledgeName)
	ids, err := uc.aiClient.Indexer.Store(ctx, docs)
	if err != nil {
		uc.log.Errorf("KnowledgeDocumentUsecase.Create Index err: %+v", gerror.Wrap(err, ""))
		return nil, err
	}

	return &pb.UploadIndexerReply{
		DocIds: ids,
	}, nil
}

func (uc *KnowledgeDocumentUsecase) Update(ctx context.Context, obj *entity.KnowledgeDocument) (*entity.KnowledgeDocument, error) {
	_, err := uc.repo.Save(ctx, obj)
	if err != nil {
		uc.log.Errorf("%+v", err)
		return nil, err
	}
	return obj, nil
}

func (uc *KnowledgeDocumentUsecase) Delete(ctx context.Context, id int64) error {
	_, err := uc.repo.Delete(ctx, id)
	if err != nil {
		uc.log.Errorf("%+v", err)
		return err
	}
	return nil
}

func (uc *KnowledgeDocumentUsecase) Get(ctx context.Context, id int64) (*entity.KnowledgeDocument, error) {
	e, err := uc.repo.Get(ctx, id)
	if err != nil {
		uc.log.Errorf("%+v", err)
		return nil, err
	}
	return e, nil
}

func (uc *KnowledgeDocumentUsecase) List(ctx context.Context, page *entity.PageAndOrder) ([]*entity.KnowledgeDocument, int64, error) {
	arr, count, err := uc.repo.List(ctx, page)
	if err != nil {
		uc.log.Errorf("%+v", err)
		return nil, 0, err
	}
	return arr, count, nil
}

func (uc *KnowledgeDocumentUsecase) ListAll(ctx context.Context) ([]*entity.KnowledgeDocument, error) {
	arr, err := uc.repo.ListAll(ctx)
	if err != nil {
		uc.log.Errorf("%+v", err)
		return nil, err
	}
	return arr, nil
}
