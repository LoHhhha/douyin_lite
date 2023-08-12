package repository

import "douyin_lite/pojo"

/*
	video`s comment_count update when "comment/action"(repository) or "comment/list"(controller).
*/

// InsertComment
// this func will update video`s comment_count (atomic operation)
func InsertComment(comment pojo.Comment) (int64, error) {
	query := "insert into comments(vid, uid, content, commentdate) values(?, ?, ?, ?)"
	res, err := Database.Exec(query, comment.VID, comment.UID, comment.Content, comment.CreateDate)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId() // get cid
	if err != nil {
		// ensure favorite don`t apply
		_ = DeleteComment(comment.ID, comment.VID)
		return 0, err
	}
	err = UpdateVideoCommentCountByVidUseIncrement(comment.VID, 1)
	if err != nil {
		// ensure favorite don`t apply
		_ = DeleteComment(comment.ID, comment.VID)
		return 0, err
	}
	return id, nil
}

// DeleteComment
// this func will update video`s comment_count (atomic operation)
func DeleteComment(id int64, vid int64) error {
	query := "delete from comments where id = ?"
	res, err := Database.Exec(query, id)
	if err != nil {
		return err
	}
	if affect, _ := res.RowsAffected(); affect != 0 {
		err = UpdateVideoCommentCountByVidUseIncrement(vid, -1)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetCommentListByVid
// @param: vid(int64)
// @return: ret([]pojo.Comment), err(error)
/* Attention: don`t think about too much Comment */
func GetCommentListByVid(vid int64) ([]pojo.Comment, error) {
	query := "SELECT id,uid,content,commentdate FROM comments WHERE vid=?"
	comments, err := Database.Query(query, vid)
	if err != nil {
		return nil, err
	}
	defer comments.Close()

	ret := make([]pojo.Comment, 0)
	for comments.Next() {
		var comment pojo.Comment
		err := comments.Scan(&comment.ID, &comment.UID, &comment.Content, &comment.CreateDate)
		if err != nil {
			return nil, err
		}

		comment.VID = vid
		ret = append(ret, comment)
	}

	return ret, nil
}
