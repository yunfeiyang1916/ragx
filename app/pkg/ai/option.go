package ai

// 是用于配置Client的函数类型
type ClientOption func(*Client)

// 设置API基础URL
func WithBaseUrl(baseUrl string) ClientOption {
	return func(c *Client) {
		c.baseUrl = baseUrl
	}
}

// 设置最大token数
func WithMaxTokens(maxTokens int) ClientOption {
	return func(c *Client) {
		if maxTokens < 1000 {
			c.maxTokens = defaultMaxTokens
		}
		c.maxTokens = maxTokens
	}
}

// 设置模型名称
func WithModel(name string) ClientOption {
	return func(c *Client) {
		c.modelName = name
	}
}

// 设置温度参数，控制生成文本的随机性
func WithTemperature(temperature float32) ClientOption {
	return func(c *Client) {
		c.temperature = temperature
	}
}

// 设置只需要chat模型
func WithOnlyChatModel(onlyChatModel bool) ClientOption {
	return func(c *Client) {
		c.onlyChatModel = onlyChatModel
	}
}

// 设置向量转换嵌入器模型名称
func WithEmbeddingModelName(name string) ClientOption {
	return func(c *Client) {
		c.embeddingModelName = name
	}
}

// 设置索引名称
func WithIndexName(name string) ClientOption {
	return func(c *Client) {
		c.indexName = name
	}
}

// 设置es配置地址
func WithESAddress(esAddress string) ClientOption {
	return func(c *Client) {
		c.esAddress = esAddress
	}
}

// 设置embedding api key
func WithEmbeddingApiKey(apiKey string) ClientOption {
	return func(c *Client) {
		c.embeddingApiKey = apiKey
	}
}
