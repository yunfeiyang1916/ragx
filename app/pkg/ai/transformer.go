package ai

import (
	"context"
	"fmt"
	"log"
	"ragx/app/pkg/utils"
	"strings"

	"github.com/cloudwego/eino-ext/components/document/loader/file"
	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/recursive"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/errors/gerror"
)

// 创建一个新的组合转换器
func NewMultiTransformer() document.Transformer {
	// 初始化基础转换器实例
	trans := &multiTransformer{}
	// 创建上下文对象，用于控制请求生命周期和传递元数据
	ctx := context.Background()

	// 配置递归分割器参数
	config := &recursive.Config{
		ChunkSize:   1000,                                    // 每个文本块大小为1000个字符
		OverlapSize: 100,                                     // 块之间有100个字符的重叠（10%），避免上下文断裂
		Separators:  []string{"\n", "。", "?", "？", "!", "！"}, // 分割符：换行符、中文句号、中英文问号和感叹号
	}

	// 创建递归分割器实例
	recTrans, err := recursive.NewSplitter(ctx, config)
	if err != nil {
		log.Fatalf("create recursive splitter failed, err: %v", gerror.Wrap(err, ""))
	}

	// 配置Markdown文档特殊处理
	mdTrans, err := markdown.NewHeaderSplitter(ctx, &markdown.HeaderConfig{
		// 标题级别映射：Markdown标题符号 -> 内部标题级别标识
		Headers: map[string]string{
			"#":   Title1, // 一级标题
			"##":  Title2, // 二级标题
			"###": Title3, // 三级标题
		},
		TrimHeaders: false, // 保留标题文本内容（不进行修剪）
	})
	if err != nil {
		log.Fatalf("create markdown splitter failed, err: %v", gerror.Wrap(err, ""))
	}

	// 将两种分割器组合到转换器中
	trans.recursive = recTrans // 通用递归分割器
	trans.markdown = mdTrans   // Markdown专用分割器

	// 返回完整的文档转换器
	return trans
}

// 组合转换器，用于对文档进行多个转换操作
type multiTransformer struct {
	// 用于处理Markdown格式的文档
	markdown document.Transformer
	// 用于处理递归文档结构的转换器
	recursive document.Transformer
}

// Transform 文档转换方法，根据文档类型选择合适的分割器进行处理
// ctx: 上下文对象，用于控制请求生命周期和传递元数据
// docs: 待处理的文档切片
// opts: 转换器选项（可变参数）
// 返回值: 处理后的文档切片，error 错误信息
func (m *multiTransformer) Transform(ctx context.Context, docs []*schema.Document, opts ...document.TransformerOption) ([]*schema.Document, error) {
	// 用于判断是否包含Markdown文档
	isMd := false
	// 遍历文档切片，检查是否包含Markdown格式文档
	for _, doc := range docs {
		// 只需要判断第一个文档是否为.md格式即可（优化性能）
		// 检查文档元数据中的扩展名字段
		if doc.MetaData["_extension"] == ".md" {
			isMd = true // 标记为包含Markdown文档
			break       // 找到第一个Markdown文档后立即退出循环
		}
	}
	// 根据文档类型选择相应的分割器进行处理
	if isMd {
		// 如果包含Markdown文档，使用Markdown标题分割器
		// 这样可以保持Markdown的标题层级结构
		return m.markdown.Transform(ctx, docs, opts...)
	}
	// 如果不包含Markdown文档，使用通用递归分割器
	// 适用于普通文本、HTML等其他格式的文档
	return m.recursive.Transform(ctx, docs, opts...)
}

// 添加文档ID并合并
// 主要功能：为文档添加唯一ID，并对Markdown文档进行智能合并
// ctx: 上下文对象，用于控制请求生命周期
// docs: 待处理的文档切片
// 返回值: 处理后的文档切片和可能的错误
func DocAddIDAndMerge(ctx context.Context, docs []*schema.Document) (output []*schema.Document, err error) {
	// 为所有没有ID的文档生成唯一ID
	for _, doc := range docs {
		if doc.ID == "" {
			doc.ID = utils.NewUUID() // 使用UUID生成器创建唯一标识符
		}
	}

	// 如果不是Markdown文档，直接返回（不进行合并操作）
	if len(docs) == 0 || docs[0].MetaData[file.MetaKeyExtension] != ".md" {
		return docs, nil
	}

	// 创建新的文档切片，用于存储合并后的文档
	ndocs := make([]*schema.Document, 0, len(docs))
	var nd *schema.Document // 当前正在合并的文档指针
	maxLen := 512           // 合并后文档的最大长度限制

	// 遍历所有文档进行智能合并
	for _, doc := range docs {
		// 检查1: 如果不是同一个源文件，结束当前合并
		if nd != nil && doc.MetaData[file.MetaKeySource] != nd.MetaData[file.MetaKeySource] {
			ndocs = append(ndocs, nd)
			nd = nil
		}

		// 检查2: 如果合并后长度超过限制，结束当前合并
		if nd != nil && len(nd.Content)+len(doc.Content) > maxLen {
			ndocs = append(ndocs, nd)
			nd = nil
		}

		// 检查3: 如果不是同一个一级标题，结束当前合并
		if nd != nil && doc.MetaData[Title1] != nd.MetaData[Title1] {
			ndocs = append(ndocs, nd)
			nd = nil
		}

		// 检查4: 如果不是同一个二级标题（且当前文档有二级标题），结束当前合并
		if nd != nil && nd.MetaData[Title2] != nil && doc.MetaData[Title2] != nd.MetaData[Title2] {
			ndocs = append(ndocs, nd)
			nd = nil
		}

		// 开始新的合并或继续当前合并
		if nd == nil {
			nd = doc // 开始新的合并组
		} else {
			// 合并二级标题
			mergeTitle(nd, doc, Title2)
			// 合并三级标题
			mergeTitle(nd, doc, Title3)
			// 合并内容
			nd.Content += doc.Content
		}
	}

	// 处理最后一个合并组
	if nd != nil {
		ndocs = append(ndocs, nd)
	}

	// 为所有合并后的文档添加Markdown标题格式
	for _, ndoc := range ndocs {
		ndoc.Content = getMdContentWithTitle(ndoc)
	}

	return ndocs, nil
}

// getMdContentWithTitle 为Markdown文档生成带标题的内容
// doc: 文档对象
// 返回值: 格式化后的内容字符串（标题 + 内容）
func getMdContentWithTitle(doc *schema.Document) string {
	if doc.MetaData == nil {
		return doc.Content // 没有元数据，直接返回内容
	}

	title := ""
	// 检查所有可能的标题级别（h1到h6）
	list := []string{"h1", "h2", "h3", "h4", "h5", "h6"}
	for _, v := range list {
		if d, e := doc.MetaData[v].(string); e && len(d) > 0 {
			title += fmt.Sprintf("%s:%s ", v, d) // 格式化标题信息
		}
	}

	if len(title) == 0 {
		return doc.Content // 没有标题信息，直接返回内容
	}

	// 返回带标题格式的内容
	return title + "\n" + doc.Content
}

// mergeTitle 合并两个文档的标题信息
// orgDoc: 原始文档（合并目标）
// addDoc: 要合并的文档
// key: 标题级别键名（如Title1, Title2等）
func mergeTitle(orgDoc, addDoc *schema.Document, key string) {
	// 如果两个文档的标题相同，不需要合并
	if orgDoc.MetaData[key] == addDoc.MetaData[key] {
		return
	}

	var title []string
	// 收集原始文档的标题
	if orgDoc.MetaData[key] != nil {
		title = append(title, orgDoc.MetaData[key].(string))
	}
	// 收集要合并文档的标题
	if addDoc.MetaData[key] != nil {
		title = append(title, addDoc.MetaData[key].(string))
	}

	// 如果有标题信息，用逗号分隔合并
	if len(title) > 0 {
		orgDoc.MetaData[key] = strings.Join(title, ",")
	}
}
