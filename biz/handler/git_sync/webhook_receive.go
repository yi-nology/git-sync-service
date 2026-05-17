package git_sync

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

const maxWebhookBodySize = 10 << 20

func ReceiveWebhook(ctx context.Context, c *app.RequestContext) {
	repoKey := c.Param("repoKey")
	if repoKey == "" {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"error": "repoKey is required",
		})
		return
	}

	bodyBytes, _ := c.Body()
	if len(bodyBytes) > maxWebhookBodySize {
		c.JSON(consts.StatusRequestEntityTooLarge, map[string]interface{}{
			"error": "request body too large",
		})
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
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	for k, v := range headers {
		httpReq.Header.Set(k, v)
	}
	httpReq.RemoteAddr = c.ClientIP()

	err = GetSyncService().ReceiveWebhook(ctx, repoKey, httpReq)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"success": true,
		"message": "webhook received",
	})
}
