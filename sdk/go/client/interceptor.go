package client

import (
	"context"
	"net/http"
)

// Interceptor представляет middleware для запросов и ответов
type Interceptor interface {
	// BeforeRequest вызывается перед отправкой запроса
	BeforeRequest(ctx context.Context, req *http.Request) error
	
	// AfterResponse вызывается после получения ответа
	AfterResponse(ctx context.Context, req *http.Request, resp *http.Response) error
}

// AddInterceptor добавляет interceptor для обработки запросов/ответов
func (c *Client) AddInterceptor(interceptor Interceptor) {
	if interceptor != nil {
		c.interceptors = append(c.interceptors, interceptor)
	}
}

// applyInterceptors применяет все interceptors перед запросом
func (c *Client) applyInterceptorsBefore(ctx context.Context, req *http.Request) error {
	for _, interceptor := range c.interceptors {
		if err := interceptor.BeforeRequest(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// applyInterceptorsAfter применяет все interceptors после ответа
func (c *Client) applyInterceptorsAfter(ctx context.Context, req *http.Request, resp *http.Response) error {
	for _, interceptor := range c.interceptors {
		if err := interceptor.AfterResponse(ctx, req, resp); err != nil {
			return err
		}
	}
	return nil
}

