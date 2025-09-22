package ai

import (
	"context"
	"fmt"
	"log"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino-ext/components/retriever/es8"
	"github.com/cloudwego/eino-ext/components/retriever/es8/search_mode"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/schema"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

// 创建并配置一个基于Elasticsearch 8的检索器
func newRetriever(c *Client) retriever.Retriever {
	// 配置ES8检索器参数
	retrieverConfig := &es8.RetrieverConfig{
		Client: c.ESClient,  // Elasticsearch客户端
		Index:  c.indexName, // 索引名称
		// 设置搜索模式为稠密向量相似度搜索（余弦相似度）
		SearchMode: search_mode.SearchModeDenseVectorSimilarity(
			search_mode.DenseVectorSimilarityTypeCosineSimilarity, // 使用余弦相似度算法
			FieldContentVector, // 指定用于向量搜索的字段
		),
		Embedding: c.Embedder, // 嵌入器，用于将文本转换为向量
		// 自定义结果解析器，用于处理ES返回的命中结果
		ResultParser: func(ctx context.Context, hit types.Hit) (doc *schema.Document, err error) {
			// 创建基础文档对象
			doc = &schema.Document{
				ID:       *hit.Id_,         // 文档ID
				MetaData: map[string]any{}, // 元数据字典
			}

			// 解析ES返回的源数据
			var src map[string]any
			if err = sonic.Unmarshal(hit.Source_, &src); err != nil {
				return nil, err
			}

			// 遍历源数据中的每个字段，根据字段类型进行相应处理
			for field, val := range src {
				switch field {
				case FieldContent: // 内容字段
					doc.Content = val.(string) // 设置文档内容
				case FieldContentVector: // 内容向量字段
					var v []float64
					// 将接口切片转换为float64切片
					for _, item := range val.([]interface{}) {
						v = append(v, item.(float64))
					}
					doc.WithDenseVector(v) // 设置文档的稠密向量
				case FieldQAContentVector, FieldQAContent: // QA相关字段
					// 这两个字段不返回给客户端，跳过处理

				case FieldExtra: // 额外信息字段
					if val == nil {
						continue // 空值跳过
					}
					doc.MetaData[FieldExtra] = val.(string) // 设置额外信息元数据
				case KnowledgeName: // 知识库名称字段
					doc.MetaData[KnowledgeName] = val.(string) // 设置知识库名称元数据
				default: // 未知字段
					return nil, fmt.Errorf("unexpected field=%s, val=%v", field, val)
				}
			}

			// 如果命中结果有评分，设置文档评分
			if hit.Score_ != nil {
				doc.WithScore(float64(*hit.Score_))
			}

			return doc, nil
		},
	}
	// 创建ES8检索器实例
	rtr, err := es8.NewRetriever(context.Background(), retrieverConfig)
	if err != nil {
		log.Fatalf("new es retriever failed, err: %+v", err) // 创建失败时记录错误并退出
	}
	return rtr
}
