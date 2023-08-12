package repository

import (
	"database/sql"
	"douyin_lite/pojo"
	"errors"
)

/*
	user`s follow_count, follower_count update when "follow/list", "follower/list"(controller).
	user`s isFollow set when it use.
*/

// GetUserFollowList
// @param: id(int64)
// @return: user_list([]pojo.User), err(error)
// this func will update user`s follow_count
// /* Attention: don`t think about too much follow */
func GetUserFollowList(id int64) ([]pojo.User, error) {
	query := "select uid from follows where follower_uid = ?"
	follows, err := Database.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer follows.Close()

	ret := make([]pojo.User, 0)
	for follows.Next() {
		var followersId int64
		err := follows.Scan(&followersId)
		if err != nil {
			return nil, err
		}

		user, found, err := FindUserById(followersId)
		if err != nil {
			return nil, err
		}

		// Follow that say user.IsFollow=true
		user.IsFollow = true

		if found {
			ret = append(ret, user)
		}
	}
	return ret, nil
}

// GetUserFollowerList
// @param: id(int64)
// @return: user_list([]pojo.User), err(error)
//this func will update user`s follower_count
/* Attention: don`t think about too much follower */
func GetUserFollowerList(id int64) ([]pojo.User, error) {
	query := "select follower_uid from follows where uid = ?"
	followers, err := Database.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer followers.Close()

	ret := make([]pojo.User, 0)
	for followers.Next() {
		var followersId int64
		err := followers.Scan(&followersId)
		if err != nil {
			return nil, err
		}

		user, found, err := FindUserById(followersId)
		if err != nil {
			return nil, err
		}

		// check if Follow this user
		if fid, _ := FindFollow(user.Id, id); fid != -1 {
			user.IsFollow = true
		} else {
			user.IsFollow = false
		}

		if found {
			ret = append(ret, user)
		}
	}
	return ret, nil
}

// GetUserFriendList
// @param: id(int64)
// @return: user_list([]pojo.User), err(error)
/* Attention: don`t think about too much friend */
func GetUserFriendList(id int64) ([]pojo.User, error) {
	query := "select follower_uid from follows where uid = ? and follower_uid in (select uid from follows where follower_uid = ?)"
	friends, err := Database.Query(query, id, id)
	if err != nil {
		return nil, err
	}
	defer friends.Close()

	ret := make([]pojo.User, 0)
	for friends.Next() {
		var followersId int64
		err := friends.Scan(&followersId)
		if err != nil {
			return nil, err
		}

		user, found, err := FindUserById(followersId)
		if err != nil {
			return nil, err
		}

		// Friend that say user.IsFollow=true
		user.IsFollow = true

		if found {
			ret = append(ret, user)
		}
	}
	return ret, nil
}

// InsertFollow
// @param: uid(int64), followerUid(int64)
// @return: err(error)
func InsertFollow(uid int64, followerUid int64) error {
	query := "INSERT INTO follows (uid, follower_uid) VALUES (?,?)"
	_, err := Database.Exec(query, uid, followerUid)
	return err
}

// DeleteFollowById
// @param: id(int64)
// @return: err(error)
/* Attention: Except that the follow is really exist */
func DeleteFollowById(id int64) error {
	query := "DELETE FROM follows WHERE id=?"
	_, err := Database.Exec(query, id)
	return err
}

// FindFollow
// @param: uid(int64), follower_uid(int64)
// @return:id(int64), err(error)
// Get (uid,follower_uid)`s id
func FindFollow(Id int64, followerId int64) (int64, error) {
	query := "SELECT id FROM follows where uid = ? AND follower_uid = ?"
	var id int64
	err := Database.QueryRow(query, Id, followerId).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, nil
		}
		return -1, err
	}
	return id, nil
}
