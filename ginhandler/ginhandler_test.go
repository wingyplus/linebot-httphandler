package ginhandler

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testRequestBody = `{
    "events": [
        {
            "replyToken": "nHuyWiB7yP5Zw52FIkcQobQuGDXCTA",
            "type": "message",
            "timestamp": 1462629479859,
            "source": {
                "type": "user",
                "userId": "u206d25c2ea6bd87c17655609a1c37cb8"
            },
            "message": {
                "id": "325708",
                "type": "text",
                "text": "Hello, world"
            }
        }
    ]
}
`

const (
	testChannelSecret = "testsecret"
	testChannelToken  = "testtoken"
)

func TestWebhookHandler(t *testing.T) {
	handler, err := New(testChannelSecret, testChannelToken)
	require.NoError(t, err)

	handler.HandleEvents(func(events []*linebot.Event, r *http.Request) {
		require.NotNil(t, events)
		require.NotNil(t, r)
		bot, err := handler.NewClient()
		require.NoError(t, err)
		require.NotNil(t, bot)
	})
	handler.HandleError(func(err error, r *http.Request) {
		assert.Equal(t, linebot.ErrInvalidSignature, err)
		assert.NotNil(t, r)
	})

	g := gin.New()
	g.POST("/webhook", handler.Handle)

	ts := httptest.NewTLSServer(g)
	defer ts.Close()

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	// valid signature
	t.Run("valid signature", func(t *testing.T) {
		body := []byte(testRequestBody)
		req, err := http.NewRequest("POST", ts.URL+"/webhook", bytes.NewReader(body))
		require.NoError(t, err)

		// generate signature
		mac := hmac.New(sha256.New, []byte(testChannelSecret))
		mac.Write(body)

		req.Header.Set("X-Line-Signature", base64.StdEncoding.EncodeToString(mac.Sum(nil)))
		res, err := httpClient.Do(req)
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("invalid signature", func(t *testing.T) {
		body := []byte(testRequestBody)
		req, err := http.NewRequest("POST", ts.URL+"/webhook", bytes.NewReader(body))
		require.NoError(t, err)

		req.Header.Set("X-LINE-Signature", "invalidsignatue")
		res, err := httpClient.Do(req)
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}
