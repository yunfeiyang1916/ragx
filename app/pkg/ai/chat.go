package ai

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	_ "github.com/joho/godotenv/autoload"
	"os"
)

type ChatModel struct {
	model.ToolCallingChatModel
}

func NewChatModel() *ChatModel {
	key := os.Getenv("OPENAI_API_KEY")
	baseURL := "https://openrouter.ai/api/v1"
	config := &openai.ChatModelConfig{
		APIKey:  key,
		Model:   "gpt-4o",
		BaseURL: baseURL,
	}
	chatModel, err := openai.NewChatModel(context.Background(), config)
	if err != nil {
		panic(err)
	}
	return &ChatModel{
		ToolCallingChatModel: chatModel,
	}
}
