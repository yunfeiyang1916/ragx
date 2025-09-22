package repo

import (
	"ragx/app/internal/biz"
	"ragx/app/internal/biz/query"

	"github.com/go-kratos/kratos/v2/log"
)

func NewKnowledgeBaseRepo(data biz.Data, logger log.Logger) biz.KnowledgeBaseRepo {
	return &KnowledgeBaseRepo{
		Data:      data,
		DB:        data.DB(),
		Rdb:       data.Rdb(),
		GormQuery: query.Use(data.DB()),
		Log:       log.NewHelper(logger),
	}
}

func NewKnowledgeChunkRepo(data biz.Data, logger log.Logger) biz.KnowledgeChunkRepo {
	return &KnowledgeChunkRepo{
		Data:      data,
		DB:        data.DB(),
		Rdb:       data.Rdb(),
		GormQuery: query.Use(data.DB()),
		Log:       log.NewHelper(logger),
	}
}

func NewKnowledgeDocumentRepo(data biz.Data, logger log.Logger) biz.KnowledgeDocumentRepo {
	return &KnowledgeDocumentRepo{
		Data:      data,
		DB:        data.DB(),
		Rdb:       data.Rdb(),
		GormQuery: query.Use(data.DB()),
		Log:       log.NewHelper(logger),
	}
}
