package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/exists"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"log"
	"ragx/app/pkg/utils"

	"github.com/cloudwego/eino-ext/components/indexer/es8"
	"github.com/cloudwego/eino/components/indexer"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/errors/gerror"
)

const (
	// 知识库名称
	KnowledgeName = "_knowledge_name"
	// 扩展数据字段
	FieldExtra = "ext"
	// 内容字段
	FieldContent = "content"
	// 内容向量字段
	FieldContentVector = "content_vector"
	// 问答内容字段
	FieldQAContent = "qa_content"
	// 问答内容向量字段
	FieldQAContentVector = "qa_content_vector"

	Title1 = "h1"
	Title2 = "h2"
	Title3 = "h3"
)

var (
	// ext 里面需要存储的数据
	ExtKeys = []string{"_extension", "_file_name", "_source", "h1", "h2", "h3"}
)

// 创建一个新的索引器
func newIndexer(c *Client) indexer.Indexer {
	if err := createIndexIfNotExists(c); err != nil {
		log.Fatalf("create index failed, err: %+v", err)
	}
	config := &es8.IndexerConfig{
		// es客户端
		Client: c.ESClient,
		// 索引名称
		Index: c.indexName,
		// 批量处理大小
		BatchSize: 10,
		Embedding: c.Embedder,
		// 将文档转换为索引字段的函数
		DocumentToFields: func(ctx context.Context, doc *schema.Document) (field2Value map[string]es8.FieldValue, err error) {
			var knowledgeName string
			// 从上下文中获取知识库名称
			if value, ok := ctx.Value(KnowledgeName).(string); ok {
				knowledgeName = value
			} else {
				err = fmt.Errorf("必须提供知识库名称")
				return
			}
			// 处理文档元数据
			if doc.MetaData != nil {
				extra, err := json.Marshal(getExtData(doc))
				if err != nil {
					return nil, err
				}
				doc.MetaData[FieldExtra] = string(extra)
			}
			// 返回字段映射，包含内容、扩展数据、知识库名称和问答内容
			return map[string]es8.FieldValue{
				// 内容字段
				FieldContent: {
					// 文档内容
					Value: doc.Content,
					// 文档内容向量字段
					EmbedKey: FieldContentVector,
				},
				// 扩展数据字段
				FieldExtra: {
					Value: doc.MetaData[FieldExtra],
				},
				// 知识库名称字段
				KnowledgeName: {
					Value: knowledgeName,
				},
				// 问答内容字段
				//FieldQAContent: {
				//	Value: doc.MetaData[FieldQAContent],
				//	// 问答内容向量字段
				//	EmbedKey: FieldQAContentVector,
				//},
			}, nil
		},
	}
	idx, err := es8.NewIndexer(context.Background(), config)
	if err != nil {
		log.Fatalf("new indexer failed, err: %+v", gerror.Wrap(err, ""))
	}
	return idx
}

// es索引是否存在，不存在则创建
func createIndexIfNotExists(c *Client) error {
	ctx := context.Background()
	indexExists, err := exists.NewExistsFunc(c.ESClient)(c.indexName).Do(ctx)
	if err != nil {
		return err
	}
	if indexExists {
		return nil
	}
	_, err = create.NewCreateFunc(c.ESClient)(c.indexName).Request(&create.Request{
		Mappings: &types.TypeMapping{
			Properties: map[string]types.Property{
				FieldContent:  types.NewTextProperty(),
				FieldExtra:    types.NewTextProperty(),
				KnowledgeName: types.NewKeywordProperty(),
				FieldContentVector: &types.DenseVectorProperty{
					Dims:       utils.Ptr(1024), // same as embedding dimensions
					Index:      utils.Ptr(true),
					Similarity: utils.Ptr("cosine"),
				},
				FieldQAContentVector: &types.DenseVectorProperty{
					Dims:       utils.Ptr(1024), // same as embedding dimensions
					Index:      utils.Ptr(true),
					Similarity: utils.Ptr("cosine"),
				},
			},
		},
	}).Do(ctx)

	return err
}

// 获取文档的扩展数据
func getExtData(doc *schema.Document) map[string]any {
	if doc.MetaData == nil {
		return nil
	}
	res := make(map[string]any)
	for _, key := range ExtKeys {
		if v, ok := doc.MetaData[key]; ok {
			res[key] = v
		}
	}
	return res
}
