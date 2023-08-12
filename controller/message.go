package controller

import (
	"douyin_lite/pojo"
	"douyin_lite/repository"
	"douyin_lite/tools"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type ChatResponse struct {
	pojo.Response
	MessageList []pojo.Message `json:"message_list"`
}

// MessageAction
// @input: token(string), to_user_id(int64), action_type(int32), content(string)
func MessageAction(c *gin.Context) {
	token := c.Query("token")
	claims, err := tools.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, pojo.Response{StatusCode: 1, StatusMsg: err.Error()})
		tools.ErrorPrint(err)
		return
	}
	fromUserId := claims.Id

	toUserId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, pojo.Response{StatusCode: 1, StatusMsg: err.Error()})
		tools.ErrorPrint(err)
		return
	}

	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, pojo.Response{StatusCode: 1, StatusMsg: err.Error()})
		tools.ErrorPrint(err)
		return
	}

	content := c.Query("content")

	if actionType == 1 {
		// "actionType == 1" mean that send message

		message := pojo.Message{
			FromUserId: fromUserId,
			ToUserId:   toUserId,
			Content:    content,
			CreateTime: strconv.FormatInt(time.Now().Unix(), 10),
		}

		err := repository.InsertMessage(&message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, pojo.Response{StatusCode: 1, StatusMsg: err.Error()})
			tools.ErrorPrint(err)
			return
		}

		c.JSON(http.StatusOK, pojo.Response{StatusCode: 0})
		return
	}

	c.JSON(http.StatusBadRequest, pojo.Response{StatusCode: 1, StatusMsg: "Error action_type"})
}

// MessageChat
// @input: token(string), to_user_id(int64), pre_msg_time:timestamp(string)
func MessageChat(c *gin.Context) {
	token := c.Query("token")
	claims, err := tools.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, pojo.Response{StatusCode: 1, StatusMsg: err.Error()})
		tools.ErrorPrint(err)
		return
	}
	fromUserId := claims.Id

	toUserId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, pojo.Response{StatusCode: 1, StatusMsg: err.Error()})
		tools.ErrorPrint(err)
		return
	}

	preMsgTime, err := strconv.ParseInt(c.Query("pre_msg_time"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, pojo.Response{StatusCode: 1, StatusMsg: err.Error()})
		tools.ErrorPrint(err)
		return
	}

	messages, err := repository.GetMessageList(toUserId, fromUserId, preMsgTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, pojo.Response{StatusCode: 1, StatusMsg: err.Error()})
		tools.ErrorPrint(err)
		return
	}

	// give time to app
	// add this code will let so many blank.
	// don`t do this will let app always give "请求信息失败".
	//	if len(messages) == 0 {
	//		messages = append(messages, pojo.Message{CreateTime: strconv.FormatInt(time.Now().Unix(), 10)})
	//	}

	c.JSON(http.StatusOK, ChatResponse{
		Response:    pojo.Response{StatusCode: 0},
		MessageList: messages,
	})
}
