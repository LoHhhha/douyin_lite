package controller

import (
	"douyin_lite/pojo"
	"douyin_lite/repository"
	"douyin_lite/tools"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserLoginResponse struct {
	pojo.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	pojo.Response
	User pojo.User `json:"user"`
}

// Register
// @input: username(string), password(string)
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	_, found, err := repository.FindUserByNameAndPassword(username, password)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: pojo.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		tools.ErrorPrint(err)
		return
	}

	if found {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: pojo.Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
		return
	}

	// id given by database
	user := pojo.User{
		Id:             -1,
		Name:           username,
		Password:       password,
		FollowerCount:  0,
		FollowCount:    0,
		TotalFavorited: 0,
	}

	id, err := repository.InsertUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UserLoginResponse{
			Response: pojo.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		tools.ErrorPrint(err)
		return
	}

	user.Id = id

	token, err := tools.GetToken(user.Id, user.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UserLoginResponse{
			Response: pojo.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		tools.ErrorPrint(err)
		return
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: pojo.Response{StatusCode: 0},
		UserId:   user.Id,
		Token:    token,
	})
}

// Login
// @input: username(string), password(string)
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	user, found, err := repository.FindUserByNameAndPassword(username, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UserLoginResponse{
			Response: pojo.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		tools.ErrorPrint(err)
		return
	}

	if !found {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: pojo.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
		return
	}

	token, err := tools.GetToken(user.Id, user.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UserLoginResponse{
			Response: pojo.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		tools.ErrorPrint(err)
		return
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: pojo.Response{StatusCode: 0},
		UserId:   user.Id,
		Token:    token,
	})
}

// UserInfo
// @input: user_id(int64), token(string)
// get user_id user`s info
func UserInfo(c *gin.Context) {
	// Only to check if the token is valid
	// later maybe we should use own user_id to check if followed given user_id
	// IsFollow need to be set in here!
	token := c.Query("token")
	if !tools.CheckToken(token) {
		c.JSON(http.StatusBadRequest, UserListResponse{
			Response: pojo.Response{StatusCode: 1, StatusMsg: "invalid token"},
		})
		return
	}

	needId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, UserResponse{
			Response: pojo.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		tools.ErrorPrint(err)
		return
	}

	user, found, err := repository.FindUserById(needId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UserLoginResponse{
			Response: pojo.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		tools.ErrorPrint(err)
		return
	}

	if found {
		c.JSON(http.StatusOK, UserResponse{
			Response: pojo.Response{StatusCode: 0},
			User:     user,
		})
		return
	}
	c.JSON(http.StatusOK, UserResponse{
		Response: pojo.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
	})
}
