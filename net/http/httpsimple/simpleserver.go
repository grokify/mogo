package httpsimple

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/apex/gateway"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

const (
	EngineAwsLambda = "awslambda"
	EngineNetHttp   = "nethttp"
	EngineFastHttp  = "fasthttp"
)

type SimpleServer interface {
	PortInt() int
	HttpEngine() string
	Router() http.Handler
	RouterFast() *fasthttprouter.Router
}

func Serve(svc SimpleServer) {
	engine := strings.ToLower(strings.TrimSpace(svc.HttpEngine()))
	if len(engine) == 0 {
		engine = EngineNetHttp
	}
	switch engine {
	case EngineNetHttp:
		log.Fatal(
			http.ListenAndServe(
				portAddress(svc.PortInt()),
				svc.Router()))
	case EngineAwsLambda:
		log.Fatal(
			gateway.ListenAndServe(
				portAddress(svc.PortInt()),
				svc.Router()))
	case EngineFastHttp:
		router := svc.RouterFast()
		if router == nil {
			log.Fatal("E_NO_FASTROUTER_FOR_ENGINE_FASTHTTP")
		}
		log.Fatal(
			fasthttp.ListenAndServe(
				portAddress(svc.PortInt()),
				router.Handler))
	default:
		log.Fatal(fmt.Sprintf("E_ENGINE_NOT_FOUND [%s]", engine))
	}
}

func portAddress(port int) string { return ":" + strconv.Itoa(port) }
