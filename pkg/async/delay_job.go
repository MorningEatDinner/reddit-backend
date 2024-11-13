package async

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/xiaorui/reddit-async/reddit-backend/pkg/base"
	"go.uber.org/zap"
)

func SendAsyncTask(body map[string]interface{}) error {
	urlStr := "/api/async/v0/job/async"
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		zap.L().Error("json.Marshal error...", zap.Error(err))
		return err
	}
	req := base.Request{
		Url:    fmt.Sprintf("%s%s", baseUrl, urlStr),
		Body:   io.NopCloser(strings.NewReader(string(jsonBytes))),
		Method: "POST",
		Params: map[string]string{},
	}
	_, _, _, err = base.Ask(req)
	if err != nil {
		zap.L().Error(" base.Ask error...", zap.Error(err))
		return err
	}
	return nil
}
