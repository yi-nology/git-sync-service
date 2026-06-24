package git_sync

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/yi-nology/git-sync-service/internal/pkg/response"
)

const maxWebhookBodySize = 10 << 20

func ReceiveWebhook(ctx context.Context, c *app.RequestContext) {
	repoKey := c.Param("repoKey")
	if repoKey == "" {
		response.BadRequest(c, "repoKey is required")
		return
	}

	bodyBytes, _ := c.Body()
	if len(bodyBytes) > maxWebhookBodySize {
		response.Error(c, consts.StatusRequestEntityTooLarge, "request body too large")
		return
	}

	headers := make(map[string]string)
	c.Request.Header.VisitAll(func(k, v []byte) {
		headers[string(k)] = string(v)
	})

	httpReq, err := http.NewRequest(
		string(c.Method()),
		string(c.Path()),
		io.LimitReader(bytes.NewReader(bodyBytes), maxWebhookBodySize),
	)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	for k, v := range headers {
		httpReq.Header.Set(k, v)
	}
	httpReq.RemoteAddr = c.ClientIP()

	err = GetSyncService().ReceiveWebhook(ctx, repoKey, httpReq)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"message": "webhook received",
	})
}
