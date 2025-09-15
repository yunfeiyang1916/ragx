package encoder

import (
	"bufio"
	"io"
	"time"

	pb "ragx/api/gen"

	kratosHttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/uuid"
)

// 实现HTTP响应编码
func StreamEncodeResponse(w kratosHttp.ResponseWriter, r *kratosHttp.Request, v interface{}) error {
	// 检查是否是文件响应
	resp, ok := v.(*pb.StreamData)
	if !ok {
		// 非文件响应，使用默认编码器
		return kratosHttp.DefaultResponseEncoder(w, r, v)
	}

	// 设置下载头
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // 禁用Nginx缓冲
	w.Header().Set("Access-Control-Allow-Origin", "*")
	sd := &pb.StreamData{
		Id:      uuid.NewString(),
		Created: time.Now().Unix(),
	}
	//if len(docs) > 0 {
	//	sd.Document = docs
	//	marshal, _ := sonic.Marshal(sd)
	//	writeSSEDocuments(httpResp, string(marshal))
	//}
	sd.Document = nil // 置空，发一次就够了

	// 创建带缓冲的写入器
	bufWriter := bufio.NewWriterSize(w, 64*1024) // 64KB 缓冲区
	defer bufWriter.Flush()

	var writer io.Writer = bufWriter

	// 检查是否支持gzip压缩
	// if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
	// 	w.Header().Set("Content-Encoding", "gzip")
	// 	gz := gzip.NewWriter(bufWriter)
	// 	defer gz.Close()
	// 	writer = gz
	// }

	writer.Write([]byte(resp.Content))

	return nil
}
