package service

import (
	"context"
	body "github.com/boram-gong/service/body"
	"github.com/gin-gonic/gin"
	"testing"
)

func Demo(ctx context.Context, request interface{}) (interface{}, error) {
	resp := body.NewCommonResp()

	return resp, nil
}

func DecodeDefault(c *gin.Context) (interface{}, error) {
	return nil, nil
}

func TestServer(t *testing.T) {
	server := NewService("9999", "")
	server.AddHTTPHandler("GET", "/test", Demo, DecodeDefault)
	server.Run()

}
