package ai

import (
	"context"
	"log"
	"ragx/app/pkg/utils"

	"github.com/cloudwego/eino-ext/components/document/loader/file"
	"github.com/cloudwego/eino-ext/components/document/loader/url"
	"github.com/cloudwego/eino-ext/components/document/parser/html"
	"github.com/cloudwego/eino-ext/components/document/parser/pdf"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/components/document/parser"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/errors/gerror"
)

// 构建组合加载器
func newLoader() document.Loader {
	ctx := context.Background()
	l := &multiLoader{}
	// 创建解析器
	htmlParser, err := html.NewParser(ctx, &html.Config{
		Selector: utils.Ptr("body"),
	})
	if err != nil {
		log.Fatalf("new html parser failed, err: %+v", gerror.Wrap(err, ""))
	}

	pdfParser, err := pdf.NewPDFParser(ctx, &pdf.Config{})
	if err != nil {
		log.Fatalf("new pdf parser failed, err: %+v", gerror.Wrap(err, ""))
	}

	// 创建扩展解析器
	p, err := parser.NewExtParser(ctx, &parser.ExtParserConfig{
		// 注册特定扩展名的解析器
		Parsers: map[string]parser.Parser{
			".html": htmlParser,
			".pdf":  pdfParser,
		},
		// 设置默认解析器，用于处理未知格式
		FallbackParser: parser.TextParser{},
	})
	if err != nil {
		log.Fatalf("new ext parser failed, err: %+v", gerror.Wrap(err, ""))
	}

	fileLoader, err := file.NewFileLoader(ctx, &file.FileLoaderConfig{
		// 是否使用名称作为文档ID
		UseNameAsID: false,
		Parser:      p,
	})
	if err != nil {
		log.Fatalf("new file loader failed, err: %+v", gerror.Wrap(err, ""))
	}
	l.fileLoader = fileLoader
	urlLoader, err := url.NewLoader(context.Background(), &url.LoaderConfig{
		Parser: p,
	})
	if err != nil {
		log.Fatalf("new url loader failed, err: %+v", gerror.Wrap(err, ""))
	}
	l.urlLoader = urlLoader
	return l
}

// 组合加载器，根据源类型选择文件加载器或URL加载器
type multiLoader struct {
	fileLoader document.Loader // 文件加载器，用于加载本地文件
	urlLoader  document.Loader // URL加载器，用于加载远程URL内容
}

// 方法根据源URI的类型选择相应的加载器进行文档加载
// ctx: 上下文，用于控制请求生命周期
// src: 文档源，包含URI等信息
// opts: 加载选项，可配置加载行为
// 返回值: 加载的文档列表和可能的错误
func (m *multiLoader) Load(ctx context.Context, src document.Source, opts ...document.LoaderOption) ([]*schema.Document, error) {
	// 判断源URI是否为URL
	if utils.IsURL(src.URI) {
		// 如果是URL，使用URL加载器
		return m.urlLoader.Load(ctx, src, opts...)
	}
	// 否则使用文件加载器
	return m.fileLoader.Load(ctx, src, opts...)
}
