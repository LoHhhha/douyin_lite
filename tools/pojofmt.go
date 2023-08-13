package tools

import (
	"douyin_lite/pojo"
	"douyin_lite/repository"
	"douyin_lite/settings"
)

// VideoList2VideoWithInfoList
// @param: videos([]pojo.Video), uid(int64), [tolIsFavorite(int)]
// @return: res([]pojo.VideoWithInfo)
// tolIsFavorite	0/none: check and set IsFavorite
//					1: 		set IsFavorite=true
//					>=2: 	set IsFavorite=false
/* this func will fill pojo.VideoWithInfo.IsFavorite and pojo.User.IsFollow */
func VideoList2VideoWithInfoList(videos []pojo.Video, uid int64, optional ...int) []pojo.VideoWithInfo {

	var tolIsFavorite = 0
	if len(optional) > 0 {
		tolIsFavorite = optional[0]
	}

	// meet: use to save element had met, so we can reduce connect to database.
	meetFollow := make(map[int64]int64)
	meetUser := make(map[int64]pojo.User)

	res := make([]pojo.VideoWithInfo, len(videos))
	for idx, video := range videos {
		res[idx].ID = video.Vid
		res[idx].PlayURL = settings.UrlPrefix + video.PlayURL
		res[idx].CoverURL = settings.UrlPrefix + video.CoverURL
		res[idx].FavoriteCount = video.FavoriteCount
		res[idx].CommentCount = video.CommentCount
		res[idx].Title = video.Title

		if tolIsFavorite == 0 {
			// check if user had favorite this videos
			var IsFavorite = false
			id, err := repository.FindFavorite(uid, video.Vid)
			if err != nil {
				ErrorPrint(err)
			}
			if id != -1 {
				IsFavorite = true
			}
			res[idx].IsFavorite = IsFavorite
		} else if tolIsFavorite == 1 {
			res[idx].IsFavorite = true
		} else {
			res[idx].IsFavorite = false
		}

		// find video`s author
		var author pojo.User
		var isFoundUser = false
		var err error = nil
		// first find in memory
		preAuthor, ok := meetUser[video.Uid]
		if ok {
			author = preAuthor
			isFoundUser = true
		} else {
			// second find in database
			author, isFoundUser, err = repository.FindUserById(video.Uid)
			if err != nil {
				ErrorPrint(err)
			} else {
				// if found, update memory
				meetUser[video.Uid] = author
			}
		}

		res[idx].Author = author

		// if "uid<=0" let IsFollow=false
		// this feature use in feed
		if uid <= 0 {
			res[idx].Author.IsFollow = false
			continue
		}

		if isFoundUser {
			// check if user is follow this video`s author
			// first find in memory
			preFid, ok := meetFollow[author.Id]
			if !ok {
				// second find in database
				fid, err := repository.FindFollow(author.Id, uid)
				if err != nil {
					ErrorPrint(err)
				} else {
					// if found, update memory
					meetFollow[author.Id] = fid
				}
				if fid != -1 {
					res[idx].Author.IsFollow = true
				}
			} else if preFid != -1 {
				res[idx].Author.IsFollow = true
			}
		}
	}
	return res
}

// CommentList2CommentWithUserList
// @param: videos([]pojo.Video), uid(int64)
// @return: res([]pojo.VideoWithInfo)
/* this func will fill pojo.User.IsFollow, let uid==0 to disable check follow. */
func CommentList2CommentWithUserList(comments []pojo.Comment, uid int64) []pojo.CommentWithUser {
	res := make([]pojo.CommentWithUser, len(comments))
	for idx, comment := range comments {
		res[idx].ID = comment.ID
		res[idx].Content = comment.Content
		res[idx].CreateDate = comment.CreateDate

		user, found, err := repository.FindUserById(comment.UID)
		if err != nil {
			ErrorPrint(err)
		}
		if found && uid != 0 {
			res[idx].User = user
			fid, err := repository.FindFollow(res[idx].User.Id, uid)
			if err != nil {
				ErrorPrint(err)
			}
			if fid != -1 {
				res[idx].User.IsFollow = true
			}
		}
	}
	return res
}
