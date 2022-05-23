package mid

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/Joggz/services/foundation/web"
)

func Panics() web.Middleware {

	m := func(handler web.Handler) web.Handler {

		h := func (ctx context.Context, w http.ResponseWriter, r *http.Request)  ( err error){
			
			defer func ()  {
				if rec := recover(); rec != nil {

					trace := debug.Stack()
					err = fmt.Errorf("PANIC [%v] TRACE [%s]", rec, string(trace) )
				}
			}()
			return handler(ctx, w, r)
		}
		return h
	}

	return m
}