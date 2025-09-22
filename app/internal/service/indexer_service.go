package service

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"ragx/app/internal/biz"

	pb "ragx/api/gen"

	"github.com/go-kratos/kratos/v2/transport/http"
)

type IndexerService struct {
	pb.UnimplementedIndexerServiceServer
	knowledgeDocumentUc *biz.KnowledgeDocumentUsecase
}

func RegisterIndexerServiceHTTPServer(s *http.Server, srv *IndexerService) {
	r := s.Route("/")
	r.POST("/api/v1/indexer", srv.UploadIndexer())
}

func NewIndexerServiceService(knowledgeDocumentUc *biz.KnowledgeDocumentUsecase) *IndexerService {
	return &IndexerService{knowledgeDocumentUc: knowledgeDocumentUc}
}

func (s *IndexerService) UploadIndexer() func(ctx http.Context) error {
	return func(ctx http.Context) error {
		http.SetOperation(ctx, pb.IndexerService_UploadIndexer_FullMethodName)
		file, header, err := ctx.Request().FormFile("file")
		if err != nil {
			return err
		}
		defer file.Close()
		//contentType := header.Header.Get("Content-Type")
		//if contentType != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		//	return err
		//}
		// 文件大小限制, 10MB
		//if header.Size > 10*1024*1024 {
		//	return err
		//}

		// 指定保存目录 - 可以根据需要修改这个路径
		saveDir := "./uploads"

		// 确保目录存在
		if err := os.MkdirAll(saveDir, 0755); err != nil {
			return err
		}

		// 构建完整的文件保存路径
		savePath := filepath.Join(saveDir, header.Filename)

		// 创建目标文件
		dst, err := os.Create(savePath)
		if err != nil {
			return err
		}
		defer dst.Close()

		// 将上传的文件内容复制到目标文件
		if _, err := io.Copy(dst, file); err != nil {
			return err
		}

		var req pb.UploadIndexerRequest
		req.KnowledgeName = ctx.Form().Get("knowledge_name")
		req.Uri = savePath
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return s.knowledgeDocumentUc.Create(ctx, req.(*pb.UploadIndexerRequest))
		})
		out, err := h(ctx, &req)
		if err != nil {
			return err
		}
		reply := out.(*pb.UploadIndexerReply)
		return ctx.Result(200, reply)
	}
}
