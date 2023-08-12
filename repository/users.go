package repository

import (
	"database/sql"
	"douyin_lite/pojo"
	"errors"
)

// InsertUser
// @param: user(*pojo.User)
// @return: id(int64), err(error)
func InsertUser(user *pojo.User) (int64, error) {
	query := "INSERT INTO users (name,password, follow_count, follower_count, total_favorited, work_count, favorite_count) VALUES (?, ?, ?, ?, ?, ?, ?)"
	res, err := Database.Exec(
		query,
		user.Name,
		user.Password,
		user.FollowCount,
		user.FollowerCount,
		user.TotalFavorited,
		user.WorkCount,
		user.FavoriteCount)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// UpdateUserFollowCountById
// @param: id(int64), newNumber(int64)
// @return: err(error)
/* Attention: Except that the user is really exist */
func UpdateUserFollowCountById(id int64, newNumber int64) error {
	queryGet := "SELECT follow_count FROM users WHERE id=?"
	queryUpdate := "UPDATE users SET follow_count=? WHERE id=? AND follow_count=?"
	return __UpdateByOptimisticLock__(queryGet, queryUpdate, id, newNumber)
}

// UpdateUserFollowerCountById
// @param: id(int64), newNumber(int64)
// @return: err(error)
/* Attention: Except that the user is really exist */
func UpdateUserFollowerCountById(id int64, newNumber int64) error {
	queryGet := "SELECT follower_count FROM users WHERE id=?"
	queryUpdate := "UPDATE users SET follower_count=? WHERE id=? AND follower_count=?"
	return __UpdateByOptimisticLock__(queryGet, queryUpdate, id, newNumber)
}

// UpdateUserTotalFavoritedById
// @param: id(int64), newNumber(int64)
// @return: err(error)
/* Attention: Except that the user is really exist */
func UpdateUserTotalFavoritedById(id int64, newNumber int64) error {
	queryGet := "SELECT total_favorited FROM users WHERE id=?"
	queryUpdate := "UPDATE users SET total_favorited=? WHERE id=? AND total_favorited=?"
	return __UpdateByOptimisticLock__(queryGet, queryUpdate, id, newNumber)
}

// UpdateUserFavoriteCountByIdUseIncrement
// @param: id(int64), increment(int64)
// @return: err(error)
/* Attention: Except that the user is really exist */
func UpdateUserFavoriteCountByIdUseIncrement(id int64, increment int64) error {
	queryGet := "SELECT favorite_count FROM users WHERE id=?"
	queryUpdate := "UPDATE users SET favorite_count=? WHERE id=? AND favorite_count=?"
	return __UpdateByOptimisticLockUseIncrement__(queryGet, queryUpdate, id, increment)
}

// UpdateUserFavoriteCountById
// @param: id(int64), increment(int64)
// @return: err(error)
/* Attention: Except that the user is really exist */
func UpdateUserFavoriteCountById(id int64, newNumber int64) error {
	queryGet := "SELECT favorite_count FROM users WHERE id=?"
	queryUpdate := "UPDATE users SET favorite_count=? WHERE id=? AND favorite_count=?"
	return __UpdateByOptimisticLock__(queryGet, queryUpdate, id, newNumber)
}

// UpdateUserWorkCountByIdUseIncrement
// @param: id(int64), increment(int64)
// @return: err(error)
/* Attention: Except that the user is really exist */
func UpdateUserWorkCountByIdUseIncrement(id int64, increment int64) error {
	queryGet := "SELECT work_count FROM users WHERE id=?"
	queryUpdate := "UPDATE users SET work_count=? WHERE id=? AND work_count=?"
	return __UpdateByOptimisticLockUseIncrement__(queryGet, queryUpdate, id, increment)
}

// UpdateUserWorkCountById
// @param: id(int64), increment(int64)
// @return: err(error)
/* Attention: Except that the user is really exist */
func UpdateUserWorkCountById(id int64, newNumber int64) error {
	queryGet := "SELECT work_count FROM users WHERE id=?"
	queryUpdate := "UPDATE users SET work_count=? WHERE id=? AND work_count=?"
	return __UpdateByOptimisticLock__(queryGet, queryUpdate, id, newNumber)
}

// FindUserByNameAndPassword
// @param: name(string), password(string)
// @return: (user)controller.User, found(bool), err(error)
// found = false when err != nil
/* Attention: isFollow is always false, it needs to check in "follows" table */
func FindUserByNameAndPassword(name string, password string) (pojo.User, bool, error) {
	var user pojo.User
	query := "SELECT id,name, follow_count, follower_count, total_favorited, work_count, favorite_count FROM users WHERE name = ? AND password = ?"
	err := Database.QueryRow(query, name, password).Scan(
		&user.Id,
		&user.Name,
		&user.FollowCount,
		&user.FollowerCount,
		&user.TotalFavorited,
		&user.WorkCount,
		&user.FavoriteCount,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, false, nil
		}
		return user, false, err
	}
	return user, true, nil
}

// FindUserById
// @param: id(int64)
// @return: (user)controller.User, found(bool), err(error)
// found = false when err != nil
/* Attention: isFollow is always false, it needs to check in "follow" table */
func FindUserById(id int64) (pojo.User, bool, error) {
	var user pojo.User
	query := "SELECT id, name, follow_count, follower_count, total_favorited, work_count, favorite_count FROM users WHERE id = ?"
	err := Database.QueryRow(query, id).Scan(
		&user.Id,
		&user.Name,
		&user.FollowCount,
		&user.FollowerCount,
		&user.TotalFavorited,
		&user.WorkCount,
		&user.FavoriteCount,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, false, nil
		}
		return user, false, err
	}
	return user, true, nil
}
