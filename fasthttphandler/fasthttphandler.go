package fasthttphandler

import (
	"bytes"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
	"github.com/valyala/fasthttp"
)

// WebhookHandler is gin adapter for line-bot-sdk-go/linebot/httphandler.WebhookHandler.
type WebhookHandler struct {
	channelSecret string
	channelToken  string

	handleEvents httphandler.EventsHandlerFunc
	handleError  httphandler.ErrorHandlerFunc
}

// New create WebhookHandler.
// TODO(thanabodee): check channel secret and channel token is it empty.
func New(channelSecret, channelToken string) (*WebhookHandler, error) {
	return &WebhookHandler{
		channelSecret: channelSecret,
		channelToken:  channelToken,
	}, nil
}

// HandleEvents register events handler function.
func (wh *WebhookHandler) HandleEvents(fn httphandler.EventsHandlerFunc) {
	wh.handleEvents = fn
}

// HandleError register error handler function.
func (wh *WebhookHandler) HandleError(fn httphandler.ErrorHandlerFunc) {
	wh.handleError = fn
}

// Handle implements gin.HandlerFunc
func (wh *WebhookHandler) Handle(ctx *fasthttp.RequestCtx) {
	req, err := httpRequestFromContext(ctx)
	if err != nil {
		ctx.Response.Header.SetStatusCode(fasthttp.StatusInternalServerError)
	}

	// TODO(thanabodee): handle error if linebot.ParseRequest return an error.
	events, _ := linebot.ParseRequest(wh.channelSecret, req)
	wh.handleEvents(events, req)
}

func httpRequestFromContext(ctx *fasthttp.RequestCtx) (*http.Request, error) {
	req, err := http.NewRequest(string(ctx.Request.Header.Method()), string(ctx.Request.Header.RequestURI()), bytes.NewBuffer(ctx.Request.Body()))
	return req, err
}

func (wh *WebhookHandler) NewClient(opts ...linebot.ClientOption) (*linebot.Client, error) {
	return linebot.New(wh.channelSecret, wh.channelToken, opts...)
}
