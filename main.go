package main

import (
	"strconv"

	"github.com/mathetake/proxy-wasm-go/runtime"
)

const validApiKey = "mashery"

func init() {

}

func main() {
	runtime.SetNewHttpContext(newContext)
}

type httpHeaders struct {
	// you must embed the default context so that you need not to reimplement all the methods by yourself
	runtime.DefaultContext
	contextID uint32
}

func newContext(contextID uint32) runtime.HttpContext {
	return &httpHeaders{contextID: contextID}
}

/**

On receiving the header named "apikey" we check if the value exists in an external redis server.
If yes, the request will be forwarded to the backend. For testing, the following values are supported: "mashery"
*/

 // override default
func (ctx *httpHeaders) OnHttpRequestHeaders(_ int) runtime.Action {
	hs, st := ctx.GetHttpRequestHeaders()
	if st != runtime.StatusOk {
		runtime.LogCritical("failed to get request headers")
		return runtime.ActionContinue
	}

	var apiKeyFound bool

	for _, h := range hs {
		runtime.LogInfo("request header: " + h[0] + ": " + h[1])
		if h[0] == "apikey" {
			apiKeyFound = true

			if h[1] != validApiKey {
				return ctx.InvalidAPIKey()
			}

			runtime.LogDebug("valid apikey found, proceed with the request")
		}
	}

	if apiKeyFound {
		ctx.DispatchHttpCall("httpbin", hs, "", [][2]string{}, 50000)
		return runtime.ActionPause
	} else {
		return ctx.InvalidAPIKey()
	}
}

// override default
func (ctx *httpHeaders) OnHttpCallResponse(_ uint32, _ int, bodySize int, _ int) {
	_, st := ctx.GetHttpCallResponseBody(0, bodySize)
	if st != runtime.StatusOk {
		runtime.LogCritical("failed to get response body")
	} else {
		runtime.LogInfo("access granted")
	}
	ctx.ResumeHttpRequest()
	return
}

// override default
func (ctx *httpHeaders) OnLog() {
	runtime.LogInfo(strconv.FormatUint(uint64(ctx.contextID), 10) + " finished")
}

func (ctx *httpHeaders) InvalidAPIKey() runtime.Action {
	runtime.LogCritical("invalid apikey value!")
	msg := "access forbidden"
	runtime.LogInfo(msg)
	ctx.SendHttpResponse(403, [][2]string{
		{"powered-by", "proxy-wasm-go!!"},
	}, msg)

	return runtime.ActionPause
}