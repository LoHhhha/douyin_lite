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

type FeedResponse struct {
	pojo.Response
	NextTime  int64                `json:"next_time"`
	VideoList []pojo.VideoWithInfo `json:"video_list"`
}

func Feed(c *gin.Context) {
	found := c.Query("last_time")
	var lastTime string
	if found == "" {
		lastTime = tools.Now2mysql()
	} else {
		unixTimestamp, err := strconv.ParseInt(found, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, pojo.Response{StatusCode: 1, StatusMsg: err.Error()})
			tools.ErrorPrint(err)
			return
		}
		lastTime = tools.Time2mysql(time.Unix(unixTimestamp, 0))
	}

	var uid int64 = -1
	token := c.Query("token")
	claims, err := tools.ParseToken(token)
	if err == nil {
		uid = claims.Id
	}

	videos, err := repository.GetFeedVideoList(lastTime, uid)
	if err != nil {
		tools.ErrorPrint(err)
		c.JSON(http.StatusInternalServerError, pojo.Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}

	var getVideoListSize = len(videos)

	list := tools.VideoList2VideoWithInfoList(videos, uid)

	var nextTime int64
	if getVideoListSize != 0 {
		nextTime = tools.Mysql2Unix(videos[len(videos)-1].Uploadtime)
	} else {
		nextTime = time.Now().Unix()
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response: pojo.Response{StatusCode: 0},
		NextTime: nextTime, VideoList: list,
	})
}
