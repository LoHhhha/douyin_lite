package repository

import (
	"database/sql"
	"errors"
)

/*
	videos`s favorite_count update when "favorite/action"(repository)
	user`s total_favorited update when "publish/list"(controller)
	user`s favorite_count update when "favorite/list"(controller) or "favorite/action"(repository)
*/

// InsertFavorite
// @param: uid(int64), vid(int64)
// @return: id(int64), err(error)
// this func will update video`s favorite_count (atomic operation)
func InsertFavorite(uid int64, vid int64) (int64, error) {
	query := "INSERT INTO favorites (uid, vid) VALUES (?, ?)"
	res, err := Database.Exec(query, uid, vid)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		// ensure favorite don`t apply
		_ = DeleteFavorite(uid, vid)
		return 0, err
	}

	// add video`s favorite_count
	err = UpdateVideoFavoriteCountByVidUseIncrement(vid, 1)
	if err != nil {
		// ensure favorite don`t apply
		_ = DeleteFavorite(uid, vid)
		return 0, err
	}

	// add user`s favorite_count
	err = UpdateUserFavoriteCountByIdUseIncrement(uid, 1)
	if err != nil {
		// ensure favorite don`t apply
		_ = UpdateVideoFavoriteCountByVidUseIncrement(vid, -1)
		_ = DeleteFavorite(uid, vid)
		return 0, err
	}
	return id, nil
}

// DeleteFavorite
// @param: uid(int64), vid(int64)
// @return: err(error)
// this func will update video`s favorite_count (atomic operation)
func DeleteFavorite(uid int64, vid int64) error {
	query := "DELETE FROM favorites WHERE uid = ? AND vid = ?"
	res, err := Database.Exec(query, uid, vid)
	if err != nil {
		return err
	}
	if affect, _ := res.RowsAffected(); affect != 0 {
		err = UpdateVideoFavoriteCountByVidUseIncrement(vid, -1)
		err = UpdateUserFavoriteCountByIdUseIncrement(uid, -1)
		if err != nil {
			return err
		}
	}
	return nil
}

// FindFavorite
// @param: uid(int64), vid(int64)
// @return: id(int64), err(error)
func FindFavorite(uid int64, vid int64) (int64, error) {
	query := "SELECT id FROM favorites where uid = ? AND vid = ?"
	var id int64
	err := Database.QueryRow(query, uid, vid).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, nil
		}
		return -1, err
	}
	return id, nil
}
