package controller

import (
	"douyin_lite/pojo"
	"douyin_lite/repository"
	"douyin_lite/tools"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserListResponse struct {
	pojo.Response
	UserList []pojo.User `json:"user_list"`
}

// FollowList
// @input: user_id(int64), token(string)
func FollowList(c *gin.Context) {
	token := c.Query("token")
	if !tools.CheckToken(token) {
		c.JSON(http.StatusBadRequest, UserListResponse{
			Response: pojo.Response{StatusCode: 1, StatusMsg: "invalid token"},
		})
		return
	}

	id, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, UserListResponse{
			Response: pojo.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		tools.ErrorPrint(err)
		return
	}

	follows, err := repository.GetUserFollowList(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UserListResponse{
			Response: pojo.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		tools.ErrorPrint(err)
		return
	}

	err = repository.UpdateUserFollowCountById(id, int64(len(follows)))
	if err != nil {
		tools.ErrorPrint(err)
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: pojo.Response{StatusCode: 0},
		UserList: follows,
	})
}

// FollowerList
// @input: user_id(int64), token(string)
func FollowerList(c *gin.Context) {
	token := c.Query("token")
	if !tools.CheckToken(token) {
		c.JSON(http.StatusBadRequest, UserListResponse{
			Response: pojo.Response{StatusCode: 1, StatusMsg: "invalid token"},
		})
		return
	}

	id, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, UserListResponse{
			Response: pojo.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		tools.ErrorPrint(err)
		return
	}

	followers, err := repository.GetUserFollowerList(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UserListResponse{
			Response: pojo.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		tools.ErrorPrint(err)
		return
	}

	err = repository.UpdateUserFollowerCountById(id, int64(len(followers)))
	if err != nil {
		tools.ErrorPrint(err)
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: pojo.Response{StatusCode: 0},
		UserList: followers,
	})
}

// FriendList
// @input: user_id(int64), token(string)
func FriendList(c *gin.Context) {
	token := c.Query("token")
	if !tools.CheckToken(token) {
		c.JSON(http.StatusBadRequest, UserListResponse{
			Response: pojo.Response{StatusCode: 1, StatusMsg: "invalid token"},
		})
		return
	}

	id, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, UserListResponse{
			Response: pojo.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		tools.ErrorPrint(err)
		return
	}

	friends, err := repository.GetUserFriendList(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UserListResponse{
			Response: pojo.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		tools.ErrorPrint(err)
		return
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: pojo.Response{StatusCode: 0},
		UserList: friends,
	})
}

// RelationAction
// @input: to_user_id(int64), action_type(int32), token(string)
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	claims, err := tools.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, UserListResponse{
			Response: pojo.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		tools.ErrorPrint(err)
		return
	}

	ownId := claims.Id

	followId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, UserListResponse{
			Response: pojo.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		tools.ErrorPrint(err)
		return
	}

	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, UserListResponse{
			Response: pojo.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		tools.ErrorPrint(err)
		return
	}

	/* because ownId is following followId, uid=followId and follower_uid=ownId */
	if actionType == 1 {
		id, err := repository.FindFollow(followId, ownId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, pojo.Response{StatusCode: 1, StatusMsg: err.Error()})
			tools.ErrorPrint(err)
			return
		}
		// follow again
		if id != -1 {
			c.JSON(http.StatusOK, pojo.Response{StatusCode: 0, StatusMsg: "Following!"})
			return
		}

		err = repository.InsertFollow(followId, ownId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, pojo.Response{StatusCode: 1, StatusMsg: err.Error()})
			tools.ErrorPrint(err)
			return
		}

		c.JSON(http.StatusOK, pojo.Response{StatusCode: 0})
		return
	} else if actionType == 2 {
		id, err := repository.FindFollow(followId, ownId)
		if err != nil {
			c.JSON(http.StatusOK, pojo.Response{StatusCode: 1, StatusMsg: err.Error()})
			tools.ErrorPrint(err)
			return
		}
		// haven`t followed
		if id == -1 {
			c.JSON(http.StatusOK, pojo.Response{StatusCode: 0, StatusMsg: "Haven`t Followed!"})
			return
		}

		err = repository.DeleteFollowById(id)
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
