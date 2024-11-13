package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreatePostHandler(t *testing.T) {

	gin.SetMode(gin.TestMode) // 输出的日志不同呗， 使用TestMode那么输出的日志就更加详细吧
	r := gin.Default()
	url := "/api/v1/post"
	r.POST(url, CreatePostHandler)

	body := `
{
    "community_id":1,
    "title":"[流言板]哈asdasdasd阿搜应被驱逐，我向范德比尔特挺身而出致敬",
    "content":"【西海岸】马丁这是新的测试禁赛，赶紧让范德彪闭嘴"
}   
	`

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(body)))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	// 这里还需要说明的是， 这个函数进行参数校验之后就结束了， 肯定会返回错误， 期望会返回错误

	//assert.Contains(t, w.Body.String(), "当前未登录")
	//方法二

	res := new(ResponseData)
	if err := json.Unmarshal(w.Body.Bytes(), res); err != nil {
		// 将数据反序列化到响应数据中
		t.Fatalf("json.Unmarshal failed. w.body failed, err: %v\n", err)
	}
	assert.Equal(t, res.Code, CodeNeedLogin)
}
