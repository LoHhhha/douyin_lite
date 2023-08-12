package controller

import (
	"douyin_lite/pojo"
	"douyin_lite/repository"
	"douyin_lite/tools"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type FavoriteListResponse struct {
	pojo.Response
	VideoList []pojo.VideoWithInfo `json:"video_list"`
}

func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	claims, err := tools.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, pojo.Response{StatusCode: 1, StatusMsg: err.Error()})
		tools.ErrorPrint(err)
		return
	}
	uid := claims.Id

	vid, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, pojo.Response{StatusCode: 1, StatusMsg: err.Error()})
		tools.ErrorPrint(err)
		return
	}

	actionType, err := strconv.Atoi(c.Query("action_type"))
	if err != nil {
		c.JSON(http.StatusBadRequest, pojo.Response{StatusCode: 1, StatusMsg: err.Error()})
		tools.ErrorPrint(err)
		return
	}

	if actionType == 1 {
		id, err := repository.FindFavorite(uid, vid)
		if err != nil {
			c.JSON(http.StatusBadRequest, pojo.Response{StatusCode: 1, StatusMsg: err.Error()})
			tools.ErrorPrint(err)
			return
		}
		// favorite again
		if id != -1 {
			c.JSON(http.StatusOK, pojo.Response{StatusCode: 0, StatusMsg: "Favoriting!"})
			return
		}

		_, err = repository.InsertFavorite(uid, vid)
		if err != nil {
			c.JSON(http.StatusBadRequest, pojo.Response{StatusCode: 1, StatusMsg: err.Error()})
			tools.ErrorPrint(err)
			return
		}

		c.JSON(http.StatusOK, pojo.Response{StatusCode: 0})
		return
	} else if actionType == 2 {
		id, err := repository.FindFavorite(uid, vid)
		if err != nil {
			c.JSON(http.StatusBadRequest, pojo.Response{StatusCode: 1, StatusMsg: err.Error()})
			tools.ErrorPrint(err)
			return
		}
		// haven`t favorited
		if id == -1 {
			c.JSON(http.StatusOK, pojo.Response{StatusCode: 0, StatusMsg: "Haven`t favorited!"})
			return
		}

		err = repository.DeleteFavorite(uid, vid)
		if err != nil {
			c.JSON(http.StatusBadRequest, pojo.Response{StatusCode: 1, StatusMsg: err.Error()})
			tools.ErrorPrint(err)
			return
		}

		c.JSON(http.StatusOK, pojo.Response{StatusCode: 0})
		return
	}

	c.JSON(http.StatusBadRequest, pojo.Response{StatusCode: 1, StatusMsg: "actionType invalid"})
}

func FavoriteList(c *gin.Context) {
	token := c.Query("token")
	if !tools.CheckToken(token) {
		c.JSON(http.StatusBadRequest, UserListResponse{
			Response: pojo.Response{StatusCode: 1, StatusMsg: "invalid token"},
		})
		return
	}

	uid, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, pojo.Response{StatusCode: 1, StatusMsg: err.Error()})
		tools.ErrorPrint(err)
		return
	}

	videos, err := repository.GetFavoriteVideoList(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, pojo.Response{StatusCode: 1, StatusMsg: err.Error()})
		tools.ErrorPrint(err)
		return
	}

	// because we get FavoriteList, so IsFavorite is actually true so set tolIsFavorite=1
	list := tools.VideoList2VideoWithInfoList(videos, uid, 1)

	// update user`s favorite_count
	err = repository.UpdateUserFavoriteCountById(uid, int64(len(list)))
	if err != nil {
		tools.ErrorPrint(err)
	}

	c.JSON(http.StatusOK, FavoriteListResponse{
		Response: pojo.Response{
			StatusCode: 0,
		},
		VideoList: list,
	})
}
