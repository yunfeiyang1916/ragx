package server

import (
	pb "ragx/api/gen"
	"ragx/app/internal/conf"
	"ragx/app/internal/service"
	logging "ragx/app/pkg/middleware/log"

	"github.com/go-kratos/kratos/v2/middleware/validate"

	staticHttp "net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, logger log.Logger,
	streamService *service.StreamService,
	kbService *service.KnowledgeBaseService,
	indexerService *service.IndexerService) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			validate.Validator(),
			logging.Server(logger),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	srv := http.NewServer(opts...)
	staticDir := "../../frontend/dist/assets"

	// 创建静态文件服务，只处理/assets/路径的请求
	fs := staticHttp.FileServer(staticHttp.Dir(staticDir))
	srv.HandlePrefix("/assets", staticHttp.StripPrefix("/assets", fs))

	// 如果需要默认首页，可以添加
	srv.HandleFunc("/", func(w staticHttp.ResponseWriter, r *staticHttp.Request) {
		staticHttp.ServeFile(w, r, "../../frontend/dist/index.html")
	})
	service.RegisterStreamServiceHTTPServer(srv, streamService)
	service.RegisterIndexerServiceHTTPServer(srv, indexerService)
	pb.RegisterKnowledgeBaseServiceHTTPServer(srv, kbService)
	return srv
}
