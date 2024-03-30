package fastweb

import (
	"github.com/oklog/run"
	"github.com/valyala/fasthttp"
)

func responseData(ctx *fasthttp.RequestCtx, data []byte) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(data)
}

func RegisterRuntime(g *run.Group) {
	srv := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			switch string(ctx.Path()) {
			case "/":
				responseData(ctx, []byte("{\"msg\":\"hello\"}"))
			}
		},
	}
	g.Add(func() error {
		return srv.ListenAndServe(":8080")
	}, func(err error) {
		srv.Shutdown()
	})
}
