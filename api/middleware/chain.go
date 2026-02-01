package middleware

import "net/http"

// Middleware はHTTPミドルウェアの型
type Middleware func(http.Handler) http.Handler

// Chain はミドルウェアチェーンを構築するビルダー
type Chain struct {
	middlewares []Middleware
}

// NewChain は新しいチェーンを作成
func NewChain(middlewares ...Middleware) *Chain {
	return &Chain{middlewares: middlewares}
}

// Use はミドルウェアをチェーンに追加
func (c *Chain) Use(mw Middleware) *Chain {
	c.middlewares = append(c.middlewares, mw)
	return c
}

// Then はハンドラーにミドルウェアチェーンを適用
// ミドルウェアは追加順（左から右）に適用される
func (c *Chain) Then(handler http.Handler) http.Handler {
	// 逆順に適用してネストを構築
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		handler = c.middlewares[i](handler)
	}
	return handler
}

// ThenFunc はHandlerFuncにミドルウェアチェーンを適用
func (c *Chain) ThenFunc(fn http.HandlerFunc) http.Handler {
	return c.Then(fn)
}

// Append は新しいミドルウェアを追加した新しいチェーンを返す（元のチェーンは変更しない）
func (c *Chain) Append(middlewares ...Middleware) *Chain {
	newMiddlewares := make([]Middleware, len(c.middlewares)+len(middlewares))
	copy(newMiddlewares, c.middlewares)
	copy(newMiddlewares[len(c.middlewares):], middlewares)
	return &Chain{middlewares: newMiddlewares}
}

// Wrap はハンドラーに単一のミドルウェアを適用するヘルパー
func Wrap(handler http.Handler, middlewares ...Middleware) http.Handler {
	return NewChain(middlewares...).Then(handler)
}
