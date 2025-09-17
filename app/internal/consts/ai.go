package consts

import (
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func PromptTemplate() prompt.ChatTemplate {
	// 创建模板，使用 FString 格式
	return prompt.FromMessages(schema.FString,
		// 系统消息模板
		schema.SystemMessage("{role}"+
			"请严格遵守以下规则：\n"+
			"1. 回答必须基于提供的参考内容，不要依赖外部知识\n"+
			"2. 如果参考内容中有明确答案，直接使用参考内容回答\n"+
			"3. 如果参考内容不完整或模糊，可以合理推断但需说明\n"+
			"4. 如果参考内容完全不相关或不存在，如实告知用户'根据现有资料无法回答'\n"+
			"5. 保持回答专业、简洁、准确\n"+
			"6. 必要时可引用参考内容中的具体数据或原文\n\n"+
			"当前提供的参考内容：\n"+
			//"{docs}\n\n"+
			""),
		// 插入需要的对话历史（新对话的话这里不填）
		// optional=false 表示必需的消息列表，在模版输入中找不到对应变量会报错，这里不需要报错
		schema.MessagesPlaceholder("chat_history", true),
		// 用户消息模板
		schema.UserMessage("问题: {question}"),
	)
}
