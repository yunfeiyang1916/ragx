package biz

import (
	"context"
	"github.com/cloudwego/eino/schema"
	"github.com/go-kratos/kratos/v2/log"
	pb "ragx/api/gen"
	"ragx/app/internal/consts"
	"ragx/app/pkg/ai"
)

type ChatUsecase struct {
	chatModel *ai.ChatModel
	log       *log.Helper
}

func NewChatUsecase(chatModel *ai.ChatModel, logger log.Logger) *ChatUsecase {
	return &ChatUsecase{
		chatModel: chatModel,
		log:       log.NewHelper(logger),
	}
}

// 将检索到的上下文和问题转换为消息列表
func (c *ChatUsecase) docsMessages(ctx context.Context, req *pb.ChatRequest) ([]*schema.Message, error) {
	// todo 从历史获取
	//chatHistory, err := x.eh.GetHistory(convID, 100)
	//if err != nil {
	//	return
	//}
	//// 插入一条用户数据
	//err = x.eh.SaveMessage(&schema.Message{
	//	Role:    schema.User,
	//	Content: question,
	//}, convID)
	//if err != nil {
	//	return
	//}
	// 使用模板生成消息
	tmpl := consts.PromptTemplate()
	messages, err := tmpl.Format(ctx, map[string]any{
		"role": "你是一个专业的AI助手，能够根据提供的参考信息准确回答用户问题。",
		//"docs":    req.Docs,
		"question": req.Question,
		// 对话历史（这个例子里模拟两轮对话历史）
		//"chat_history": []*schema.Message{
		//	schema.UserMessage("你好"),
		//	schema.AssistantMessage("嘿！我是你的程序员鼓励师！记住，每个优秀的程序员都是从 Debug 中成长起来的。有什么我可以帮你的吗？", nil),
		//	schema.UserMessage("我觉得自己写的代码太烂了"),
		//	schema.AssistantMessage("每个程序员都经历过这个阶段！重要的是你在不断学习和进步。让我们一起看看代码，我相信通过重构和优化，它会变得更好。记住，Rome wasn't built in a day，代码质量是通过持续改进来提升的。", nil),
		//},
	})
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (c *ChatUsecase) Chat(ctx context.Context, req *pb.ChatRequest) (*pb.ChatReply, error) {
	return &pb.ChatReply{}, nil
}

func (c *ChatUsecase) ChatStream(ctx context.Context, req *pb.ChatRequest) (*schema.StreamReader[*schema.Message], error) {
	// 转换为消息列表
	messages, err := c.docsMessages(ctx, req)
	if err != nil {
		return nil, err
	}
	// 流式调用模型
	sr, err := c.chatModel.Stream(ctx, messages)
	if err != nil {
		return nil, err
	}
	return sr, err
}
