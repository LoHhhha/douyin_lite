package controller

import (
	"douyin_lite/pojo"
	"douyin_lite/repository"
	"douyin_lite/tools"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CommentResponse struct {
	pojo.Response
	Comment pojo.CommentWithUser `json:"comment"`
}

type CommentListResponse struct {
	pojo.Response
	CommentList []pojo.CommentWithUser `json:"comment_list,omitempty"`
}

func CommentAction(c *gin.Context) {
	token := c.Query("token")
	claim, err := tools.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, pojo.Response{StatusCode: 1, StatusMsg: "token is invalid"})
		tools.ErrorPrint(err)
		return
	}

	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, pojo.Response{StatusCode: 1, StatusMsg: "parse vid failed"})
		tools.ErrorPrint(err)
		return
	}

	actionType, err := strconv.Atoi(c.Query("action_type"))
	if err != nil {
		c.JSON(http.StatusBadRequest, pojo.Response{StatusCode: 1, StatusMsg: "parse action_type failed"})
		tools.ErrorPrint(err)
		return
	}

	if actionType == 1 {
		comment := pojo.Comment{
			VID:        videoId,
			UID:        claim.Id,
			Content:    c.Query("comment_text"),
			CreateDate: tools.Now2mysql(),
		}

		id, err := repository.InsertComment(comment)
		if err != nil {
			c.JSON(http.StatusInternalServerError, pojo.Response{StatusCode: 1, StatusMsg: err.Error()})
			tools.ErrorPrint(err)
			return
		}

		user, found, err := repository.FindUserById(comment.UID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, pojo.Response{StatusCode: 1, StatusMsg: err.Error()})
			tools.ErrorPrint(err)
			return
		}
		if found == false {
			c.JSON(http.StatusBadRequest, pojo.Response{StatusCode: 1, StatusMsg: "missing user id"})
			return
		}

		c.JSON(http.StatusOK, CommentResponse{
			Response: pojo.Response{
				StatusCode: 0,
			},
			Comment: pojo.CommentWithUser{
				ID:         id,
				User:       user,
				Content:    comment.Content,
				CreateDate: comment.CreateDate,
			},
		})
		return
	} else if actionType == 2 {
		commentId, _ := strconv.ParseInt(c.Query("comment_id"), 10, 64)
		err := repository.DeleteComment(commentId, videoId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, pojo.Response{StatusCode: 1, StatusMsg: err.Error()})
			tools.ErrorPrint(err)
			return
		}
		c.JSON(http.StatusOK, pojo.Response{StatusCode: 0})
		return
	}

	c.JSON(http.StatusBadRequest, pojo.Response{StatusCode: 1, StatusMsg: "actionType invalid"})
}

func CommentList(c *gin.Context) {
	// allow not login users check someone FavoriteList
	// attention: if token != nil, we still check it.
	var uid int64 = 0
	token := c.Query("token")
	if token != "" {
		claim, err := tools.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusBadRequest, pojo.Response{StatusCode: 1, StatusMsg: "token is invalid"})
			tools.ErrorPrint(err)
			return
		}
		uid = claim.Id
	}

	vid, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, pojo.Response{StatusCode: 1, StatusMsg: "parse vid failed"})
		tools.ErrorPrint(err)
		return
	}

	comments, err := repository.GetCommentListByVid(vid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, pojo.Response{StatusCode: 1, StatusMsg: err.Error()})
		tools.ErrorPrint(err)
		return
	}

	list := tools.CommentList2CommentWithUserList(comments, uid)

	// update video`s comment_count
	err = repository.UpdateVideoCommentCountByVid(vid, int64(len(list)))
	if err != nil {
		tools.ErrorPrint(err)
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    pojo.Response{StatusCode: 0},
		CommentList: list,
	})
}
