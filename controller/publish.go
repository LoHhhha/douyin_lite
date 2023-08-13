package controller

import (
	"douyin_lite/pojo"
	"douyin_lite/repository"
	"douyin_lite/settings"
	"douyin_lite/tools"
	"net/http"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// empty struct did not take up space
var videoExt = map[string]struct{}{
	".mp4": {},
	".avi": {},
	".mov": {},
	".flv": {},
}

type PublishListResponse struct {
	pojo.Response
	VideoList []pojo.VideoWithInfo `json:"video_list"`
}

func Publish(c *gin.Context) {
	token := c.PostForm("token")
	claim, err := tools.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, pojo.Response{StatusCode: 1, StatusMsg: "token is invalid"})
		tools.ErrorPrint(err)
		return
	}

	title := c.PostForm("title")

	file, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusBadRequest, pojo.Response{StatusCode: 1, StatusMsg: err.Error()})
		tools.ErrorPrint(err)
		return
	}

	// check if video extension is valid
	ext := filepath.Ext(file.Filename)
	if _, ok := videoExt[ext]; !ok {
		c.JSON(http.StatusBadRequest, pojo.Response{StatusCode: 1, StatusMsg: "Video format error"})
		tools.ErrorPrint(err)
		return
	}

	// get video/cover`s path
	// use "timestamp_Filename" to get the path name.
	now := time.Now()
	var builder strings.Builder // using strings.Builder is more efficient
	builder.WriteString("./public/videos/")
	builder.WriteString(strconv.FormatInt(now.Unix(), 10)) // generate unique filename
	builder.WriteString("_")
	builder.WriteString(file.Filename)

	videoSavePath := builder.String()

	coverSavePath := strings.Replace(videoSavePath, "./public/videos", "./public/covers", 1)
	coverSavePath = strings.TrimSuffix(coverSavePath, ext)
	coverSavePath += ".jpg"

	err = c.SaveUploadedFile(file, videoSavePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, pojo.Response{StatusCode: 1, StatusMsg: err.Error()})
		tools.ErrorPrint(err)
		return
	}

	// Take the first frame from the video as the cover image
	var cmd *exec.Cmd
	if settings.ServerIsLinux {
		cmd = exec.Command(
			"./tools/ffmpeg",
			"-i", videoSavePath,
			"-ss", "00:00:01",
			"-vframes",
			"1",
			coverSavePath,
		)
	} else {
		cmd = exec.Command(
			"./tools/ffmpeg.exe",
			"-i",
			videoSavePath,
			"-ss", "00:00:01",
			"-vframes",
			"1",
			coverSavePath,
		)
	}
	err = cmd.Run()
	if err != nil {
		c.JSON(http.StatusInternalServerError, pojo.Response{StatusCode: 1, StatusMsg: err.Error()})
		tools.ErrorPrint(err)
		return
	}

	playURL := strings.Replace(videoSavePath, "./public/videos", "/static/videos", 1)
	coverURL := strings.Replace(coverSavePath, "./public/covers", "/static/covers", 1)

	video := pojo.Video{
		Vid:        -1,
		Uid:        claim.Id,
		PlayURL:    playURL,
		CoverURL:   coverURL,
		Title:      title,
		Uploadtime: tools.Now2mysql(),
	}

	err = repository.InsertVideo(video)
	if err != nil {
		c.JSON(http.StatusInternalServerError, pojo.Response{StatusCode: 1, StatusMsg: err.Error()})
		tools.ErrorPrint(err)
		return
	}

	c.JSON(http.StatusOK, pojo.Response{StatusCode: 0})
}

func PublishList(c *gin.Context) {
	// allow not login users check someone PublishList
	// attention: if token != nil, we still check it.
	token := c.Query("token")
	if token != "" && tools.CheckToken(token) == false {
		c.JSON(http.StatusBadRequest, pojo.Response{StatusCode: 1, StatusMsg: "token is invalid"})
		return
	}

	userId := c.Query("user_id")
	uid, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, pojo.Response{StatusCode: 1, StatusMsg: "uid is invalid"})
		tools.ErrorPrint(err)
		return
	}

	videos, err := repository.GetPublishVideoList(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, pojo.Response{StatusCode: 1, StatusMsg: err.Error()})
		tools.ErrorPrint(err)
		return
	}

	// update user`s total_favorited
	var totalFavorited int64 = 0
	for _, video := range videos {
		totalFavorited += video.FavoriteCount
	}
	err = repository.UpdateUserTotalFavoritedById(uid, totalFavorited)
	if err != nil {
		tools.ErrorPrint(err)
	}

	// update user`s work_count
	err = repository.UpdateUserWorkCountById(uid, int64(len(videos)))
	if err != nil {
		tools.ErrorPrint(err)
	}

	list := tools.VideoList2VideoWithInfoList(videos, uid)

	c.JSON(http.StatusOK, PublishListResponse{
		Response: pojo.Response{
			StatusCode: 0,
		},
		VideoList: list,
	})
}
