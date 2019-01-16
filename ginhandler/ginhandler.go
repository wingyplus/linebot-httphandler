package ginhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
)

// WebhookHandler is gin wrapper for line-bot-sdk-go/linebot/httphandler.WebhookHandler.
type WebhookHandler struct {
	httphandlerWebhookHandler *httphandler.WebhookHandler
}

// New create WebhookHandler.
func New(channelSecret, channelToken string) (*WebhookHandler, error) {
	handler, err := httphandler.New(channelSecret, channelToken)
	if err != nil {
		return nil, err
	}
	return &WebhookHandler{
		httphandlerWebhookHandler: handler,
	}, nil
}

// HandleEvents register events handler function.
func (wh *WebhookHandler) HandleEvents(fn httphandler.EventsHandlerFunc) {
	wh.httphandlerWebhookHandler.HandleEvents(fn)
}

// HandleError register error handler function.
func (wh *WebhookHandler) HandleError(fn httphandler.ErrorHandlerFunc) {
	wh.httphandlerWebhookHandler.HandleError(fn)
}

// Handle implements gin.HandlerFunc
func (wh *WebhookHandler) Handle(ctx *gin.Context) {
	wh.httphandlerWebhookHandler.ServeHTTP(ctx.Writer, ctx.Request)
}

func (wh *WebhookHandler) NewClient(opts ...linebot.ClientOption) (*linebot.Client, error) {
	return wh.httphandlerWebhookHandler.NewClient(opts...)
}
