package ai

import (
	"context"
	"log"
	"ragx/app/pkg/utils"

	openaiEmbedding "github.com/cloudwego/eino-ext/components/embedding/openai"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/indexer"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/elastic/go-elasticsearch/v8"
)

const (
	// 默认API基础URL
	DefaultBaseUrl = "https://openrouter.ai/api/v1"
	// 默认使用的模型
	DefaultModel = "gpt-4o"
	// 默认最大token数
	defaultMaxTokens = 8192
)

type Client struct {
	// 最大token数
	maxTokens int
	// 模型名称
	modelName string
	// 只需要chat模型
	onlyChatModel bool
	// 向量转换嵌入器模型名称
	embeddingModelName string
	// 索引名
	indexName string
	// 索引器的相关配置：es地址
	esAddress  string
	esUsername string
	esPassword string
	// API密钥
	apiKey string
	// API基础URL
	baseUrl string
	/*
		温度参数，控制生成文本的随机性
		| Temperature Value | Randomness   | Applicable Scenarios                      |
		|----------------|------------------------|-------------------------------------------|
		| 0                 | No randomness         | Factual answers, code generation, technical documentation |
		| 0.5 - 0.7      | Moderate randomness | Conversational systems, content creation, recommendation systems |
		| 1                 | High randomness       | Creative writing, brainstorming, advertising copy |
		| 1.5 - 2         | Extreme randomness  | Artistic creation, game design, exploratory tasks |
	*/
	temperature float32

	// 模型，用于生成文本或执行其他模型相关操作
	ChatModel model.ToolCallingChatModel
	// 嵌入器，用于将文本转换为向量表示。主要用于文档检索和相似度计算。
	Embedder embedding.Embedder

	// 加载器，用于加载文档
	Loader document.Loader
	// 转换器，对输入的文档进行各种转换操作，如分割、过滤、合并等，从而得到满足特定需求的文档
	Transformer document.Transformer
	// es客户端
	ESClient *elasticsearch.Client
	// 索引器，用于将文档转换为索引，以便进行快速检索
	Indexer indexer.Indexer
	// 检索器，用于根据查询条件从索引中检索相关文档
	Retriever retriever.Retriever
}

func NewClient(apiKey string, opts ...ClientOption) *Client {
	c := &Client{
		apiKey:             apiKey,
		baseUrl:            DefaultBaseUrl,
		maxTokens:          defaultMaxTokens,
		modelName:          DefaultModel,
		embeddingModelName: "text2vec-large-chinese",
		temperature:        0,
		onlyChatModel:      true,
	}
	// 应用所有选项
	for _, opt := range opts {
		opt(c)
	}
	// 初始化模型
	config := &openai.ChatModelConfig{
		APIKey:  c.apiKey,
		Model:   c.modelName,
		BaseURL: c.baseUrl,
	}
	chatModel, err := openai.NewChatModel(context.Background(), config)
	if err != nil {
		log.Fatalf("new openai chat model failed, err: %+v", err)
	}
	c.ChatModel = chatModel
	if c.onlyChatModel {
		return c
	}
	// 初始化嵌入器
	embeddingConfig := &openaiEmbedding.EmbeddingConfig{
		APIKey:     c.apiKey,
		Model:      c.embeddingModelName,
		Dimensions: utils.Ptr(1024),
		Timeout:    0,
		BaseURL:    c.baseUrl,
	}
	embedder, err := openaiEmbedding.NewEmbedder(context.Background(), embeddingConfig)
	if err != nil {
		log.Fatalf("new openai embedding failed, err: %+v", err)
	}
	c.Embedder = embedder
	// 初始化加载器
	c.Loader = newLoader()
	// 初始化转换器
	c.Transformer = NewMultiTransformer()
	// 初始化es客户端
	c.ESClient, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{c.esAddress},
		Username:  c.esUsername,
		Password:  c.esPassword,
	})
	if err != nil {
		log.Fatalf("new es client failed, err: %+v", err)
	}
	// 初始化索引器
	c.Indexer = newIndexer(c)
	// 初始化检索器
	c.Retriever = newRetriever(c)
	return c
}
