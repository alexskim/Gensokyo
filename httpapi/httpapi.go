package httpapi

import (
	"github.com/hoshinonyaruko/gensokyo/config"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hoshinonyaruko/gensokyo/callapi"
	"github.com/hoshinonyaruko/gensokyo/handlers"
	"github.com/tencent-connect/botgo/openapi"
)

// CombinedMiddleware 创建并返回一个带有依赖的中间件闭包
func CombinedMiddleware(api openapi.OpenAPI, apiV2 openapi.OpenAPI) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := config.GetHTTPAccessToken()
		if accessToken != "" {
			tokenHeader := strings.Replace(c.GetHeader("Authorization"), "Bearer ", "", 1)
			tokenQuery, _ := c.GetQuery("access_token")
			if (tokenHeader == "" || tokenHeader != accessToken) && (tokenQuery == "" || tokenQuery != accessToken) {
				c.JSON(http.StatusForbidden, gin.H{"error": "鉴权失败"})
				return
			}
		}
		// 检查路径和处理对应的请求
		if c.Request.URL.Path == "/send_group_msg" {
			handleSendGroupMessage(c, api, apiV2)
			return
		}
		if c.Request.URL.Path == "/send_private_msg" {
			handleSendPrivateMessage(c, api, apiV2)
			return
		}
		if c.Request.URL.Path == "/send_guild_channel_msg" {
			handleSendGuildChannelMessage(c, api, apiV2)
			return
		}
		if c.Request.URL.Path == "/get_group_list" {
			handleGetGroupList(c, api, apiV2)
			return
		}
		if c.Request.URL.Path == "/put_interaction" {
			handlePutInteraction(c, api, apiV2)
			return
		}

		// 调用c.Next()以继续处理请求链
		c.Next()
	}
}

// handleSendGroupMessage 处理发送群聊消息的请求
func handleSendGroupMessage(c *gin.Context, api openapi.OpenAPI, apiV2 openapi.OpenAPI) {
	var retmsg string
	var req struct {
		GroupID    int64  `json:"group_id" form:"group_id"`
		UserID     *int64 `json:"user_id,omitempty" form:"user_id"`
		Message    string `json:"message" form:"message"`
		AutoEscape bool   `json:"auto_escape" form:"auto_escape"`
	}

	// 根据请求方法解析参数
	if c.Request.Method == http.MethodGet {
		// 从URL查询参数解析
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		// 从JSON或表单数据解析
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	// 使用解析后的参数处理请求
	client := &HttpAPIClient{}
	// 创建 ActionMessage 实例
	message := callapi.ActionMessage{
		Action: "send_group_msg",
		Params: callapi.ParamsContent{
			GroupID: strconv.FormatInt(req.GroupID, 10), // 注意这里需要转换类型，因为 GroupID 是 int64
			Message: req.Message,
		},
	}
	// 如果 UserID 存在，则加入到参数中
	if req.UserID != nil {
		message.Params.UserID = strconv.FormatInt(*req.UserID, 10)
	}
	// 调用处理函数
	retmsg, err := handlers.HandleSendGroupMsg(client, api, apiV2, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回处理结果
	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, retmsg)
}

// handleSendPrivateMessage 处理发送私聊消息的请求
func handleSendPrivateMessage(c *gin.Context, api openapi.OpenAPI, apiV2 openapi.OpenAPI) {
	var retmsg string
	var req struct {
		GroupID    int64  `json:"group_id" form:"group_id"`
		UserID     int64  `json:"user_id" form:"user_id"`
		Message    string `json:"message" form:"message"`
		AutoEscape bool   `json:"auto_escape" form:"auto_escape"`
	}

	// 根据请求方法解析参数
	if c.Request.Method == http.MethodGet {
		// 从URL查询参数解析
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		// 从JSON或表单数据解析
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	// 使用解析后的参数处理请求
	// TODO: 添加请求处理逻辑
	// 例如：api.SendGroupMessage(req.GroupID, req.Message, req.AutoEscape)
	client := &HttpAPIClient{}
	// 创建 ActionMessage 实例
	message := callapi.ActionMessage{
		Action: "send_private_msg",
		Params: callapi.ParamsContent{
			GroupID: strconv.FormatInt(req.GroupID, 10), // 注意这里需要转换类型，因为 GroupID 是 int64
			UserID:  strconv.FormatInt(req.UserID, 10),
			Message: req.Message,
		},
	}
	// 调用处理函数
	retmsg, err := handlers.HandleSendPrivateMsg(client, api, apiV2, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回处理结果
	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, retmsg)
}

// handleSendGuildChannelMessage 处理发送消频道息的请求
func handleSendGuildChannelMessage(c *gin.Context, api openapi.OpenAPI, apiV2 openapi.OpenAPI) {
	var retmsg string
	var req struct {
		GuildID    string `json:"guild_id" form:"guild_id"`
		ChannelID  string `json:"channel_id" form:"channel_id"`
		Message    string `json:"message" form:"message"`
		AutoEscape bool   `json:"auto_escape" form:"auto_escape"`
	}

	// 根据请求方法解析参数
	if c.Request.Method == http.MethodGet {
		// 从URL查询参数解析
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		// 从JSON或表单数据解析
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	// 使用解析后的参数处理请求
	// TODO: 添加请求处理逻辑
	// 例如：api.SendGroupMessage(req.GroupID, req.Message, req.AutoEscape)
	client := &HttpAPIClient{}
	// 创建 ActionMessage 实例
	message := callapi.ActionMessage{
		Action: "send_guild_channel_msg",
		Params: callapi.ParamsContent{
			GuildID:   req.GuildID,
			ChannelID: req.ChannelID, // 注意这里需要转换类型，因为 GroupID 是 int64
			Message:   req.Message,
		},
	}
	// 调用处理函数
	retmsg, err := handlers.HandleSendGuildChannelMsg(client, api, apiV2, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回处理结果
	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, retmsg)
}

// 定义了一个符合 Client 接口的 HttpAPIClient 结构体
type HttpAPIClient struct {
	// 可添加所需字段
}

// 实现 Client 接口的 SendMessage 方法
// 假client中不执行任何操作，只是返回 nil 来符合接口要求
func (c *HttpAPIClient) SendMessage(message map[string]interface{}) error {
	// 不实际发送消息
	// log.Printf("SendMessage called with: %v", message)

	// 返回nil占位符
	return nil
}

// handleGetGroupList 处理获取群列表
func handleGetGroupList(c *gin.Context, api openapi.OpenAPI, apiV2 openapi.OpenAPI) {
	var retmsg string

	// 使用解析后的参数处理请求
	client := &HttpAPIClient{}
	// 创建 ActionMessage 实例
	message := callapi.ActionMessage{
		Action: "get_group_list",
	}

	// 调用处理函数
	retmsg, err := handlers.GetGroupList(client, api, apiV2, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回处理结果
	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, retmsg)
}

// handlePutInteraction 处理put_interaction的请求
func handlePutInteraction(c *gin.Context, api openapi.OpenAPI, apiV2 openapi.OpenAPI) {
	var req struct {
		Echo     string `json:"echo" form:"echo"`           // Echo值用于标识interaction
		PostType string `json:"post_type" form:"post_type"` // PostType用于设置code参数
	}

	// 根据请求方法解析参数
	if c.Request.Method == http.MethodGet {
		// 从URL查询参数解析
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		// 从JSON或表单数据解析
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	// 创建 ActionMessage 实例
	message := callapi.ActionMessage{
		Action:   "put_interaction",
		Echo:     req.Echo,
		PostType: req.PostType,
	}

	// 调用处理函数
	client := &HttpAPIClient{} // 假设HttpAPIClient实现了callapi.Client接口
	retmsg, err := handlers.HandlePutInteraction(client, api, apiV2, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回处理结果
	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, retmsg)
}
