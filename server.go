package service

import (
	"fmt"
	"github.com/boram-gong/service/svc"
	"github.com/boram-gong/service/svc/endpoint"
	svc_http "github.com/boram-gong/service/svc/http"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	logs "github.com/lestrrat-go/file-rotatelogs"
)

func interruptHandler(ch chan<- error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	terminateError := fmt.Errorf("%s", <-c)
	ch <- terminateError
}

type Service struct {
	Server *gin.Engine
	Port   string
	Log    string
}

func NewService(port, outLogPath string) *Service {
	var service = new(Service)
	service.Log = outLogPath
	service.Port = port
	if outLogPath != "" {
		gin.DisableConsoleColor()
		writer, _ := logs.New(
			outLogPath + "%Y%m%d.log",
		)
		gin.DefaultWriter = io.MultiWriter(writer)
	}
	service.Server = gin.Default()
	service.Server.Use(cors.Default())
	service.Server.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))
	return service
}

func (s *Service) Run() {
	ch := make(chan error)
	go interruptHandler(ch)
	// http
	go func() {
		ch <- s.Server.Run(":" + s.Port)
	}()
	log.Printf("closed:%s \n", <-ch)
}

func (s *Service) AddHTTPHandler(httpMethod, relativePath string, f endpoint.Endpoint, decode svc_http.DecodeRequestFuncFromGin) {
	// json-adapter
	s.Server.Handle(httpMethod, relativePath, func(c *gin.Context) {
		svc_http.NewServer(
			f,
			svc_http.WrapS(c, decode),
			svc.EncodeHTTPGenericResponse,
			serverOptions...,
		).ServeHTTP(c.Writer, c.Request)
	})
}

var (
	serverOptions = []svc_http.ServerOption{
		svc_http.ServerBefore(svc.HeadersToContext),
		svc_http.ServerErrorEncoder(svc.ErrorEncoder),
		svc_http.ServerErrorHandler(svc_http.NewNopErrorHandler()),
		svc_http.ServerAfter(svc_http.SetContentType(svc.ContentType)),
	}
)
