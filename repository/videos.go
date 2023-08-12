package repository

import (
	"douyin_lite/pojo"
	"douyin_lite/settings"
)

/*
	user`s work_count update when "publish/list"(controller)
*/

// GetVideoListByCommand
// @param: format(string)           查询格式
// @param: args...([]interface{})   参数列表
// @return: ([]pojo.Video, error)   视频列表和错误信息
func GetVideoListByCommand(format string, args ...interface{}) ([]pojo.Video, error) {
	stmt, err := Database.Prepare(format)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}

	videos := make([]pojo.Video, 0)

	defer rows.Close() // very important: Close rows to release held database links

	for rows.Next() {
		var video pojo.Video
		rows.Scan(&video.Vid, &video.Uid, &video.PlayURL, &video.CoverURL, &video.FavoriteCount,
			&video.CommentCount, &video.Uploadtime, &video.Title)
		videos = append(videos, video)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return videos, err
}

// InsertVideo
// @param: video(pojo.Video)
// @return: err(error)
func InsertVideo(video pojo.Video) error {
	query := "INSERT INTO videos (uid, play_url, cover_url, uploadtime, title) VALUES (?,?,?,?,?)"
	_, err := Database.Exec(query, video.Uid, video.PlayURL, video.CoverURL,
		video.Uploadtime, video.Title)
	return err
}

// UpdateVideoFavoriteCountByVidUseIncrement
// @param: vid(int64), increment(int64)
// @return: err(error)
// use optimistic lock to ensure all is fine
func UpdateVideoFavoriteCountByVidUseIncrement(vid int64, increment int64) error {
	queryGet := "SELECT favorite_count FROM videos WHERE vid=?"
	queryUpdate := "UPDATE videos SET favorite_count=? WHERE vid=? AND favorite_count=?"
	return __UpdateByOptimisticLockUseIncrement__(queryGet, queryUpdate, vid, increment)
}

// UpdateVideoCommentCountByVidUseIncrement
// @param: vid(int64), increment(int64)
// @return: err(error)
// use optimistic lock to ensure all is fine
func UpdateVideoCommentCountByVidUseIncrement(vid int64, increment int64) error {
	queryGet := "SELECT comment_count FROM videos WHERE vid=?"
	queryUpdate := "UPDATE videos SET comment_count=? WHERE vid=? AND comment_count=?"
	return __UpdateByOptimisticLockUseIncrement__(queryGet, queryUpdate, vid, increment)
}

// UpdateVideoCommentCountByVid
// @param: vid(int64), newNumber(int64)
// @return: err(error)
// use optimistic lock to ensure all is fine
func UpdateVideoCommentCountByVid(vid int64, newNumber int64) error {
	queryGet := "SELECT comment_count FROM videos WHERE vid=?"
	queryUpdate := "UPDATE videos SET comment_count=? WHERE vid=? AND comment_count=?"
	return __UpdateByOptimisticLock__(queryGet, queryUpdate, vid, newNumber)
}

// GetPublishVideoList
// @param: uid(int64)
// @return: videoList([]pojo.Video), err(error)
func GetPublishVideoList(uid int64) ([]pojo.Video, error) {
	var ListGetFormat = "SELECT * FROM videos WHERE uid = ?"
	return GetVideoListByCommand(ListGetFormat, uid)
}

// GetFeedVideoList
// @param: lastTime(string), uid(int64)
// @return: videos([]pojo.Video),err(error)
func GetFeedVideoList(lastTime string, uid int64) ([]pojo.Video, error) {
	/* we can use uid to support our recommended algorithm, but now not need to use. */
	var FeedGetFormat = "SELECT * FROM videos WHERE uploadtime < ? ORDER BY uploadtime DESC LIMIT?"
	return GetVideoListByCommand(FeedGetFormat, lastTime, settings.FeedListMaxNum)
}

// GetFavoriteVideoList
// @param: lastTime(string)
// @return: videos([]pojo.Video),err(error)
func GetFavoriteVideoList(uid int64) ([]pojo.Video, error) {
	var FeedGetFormat = "SELECT videos.* FROM favorites join videos on favorites.vid = videos.vid where favorites.uid = ?"
	return GetVideoListByCommand(FeedGetFormat, uid)
}
