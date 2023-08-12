package repository

import (
	"douyin_lite/pojo"
	"strconv"
)

// GetMessageList
// @param: toUserId(int64), fromUserId(int64), time(time.Time)
// @return: messages([]pojo.Message), err(error)
/* Attention: don`t think about too much message */
func GetMessageList(toUserId int64, fromUserId int64, preTime int64) ([]pojo.Message, error) {
	query := "select id,to_user_id,from_user_id,content,create_time from messages WHERE ((to_user_id = ? AND from_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)) AND create_time > ? ORDER BY create_time"
	messages, err := Database.Query(query, toUserId, fromUserId, toUserId, fromUserId, preTime)
	if err != nil {
		return nil, err
	}
	defer messages.Close()

	ret := make([]pojo.Message, 0)

	for messages.Next() {
		var message pojo.Message
		var timestamp int64
		err := messages.Scan(&message.Id, &message.ToUserId, &message.FromUserId, &message.Content, &timestamp)
		if err != nil {
			return nil, err
		}
		message.CreateTime = strconv.FormatInt(timestamp, 10)

		ret = append(ret, message)
	}

	return ret, nil
}

// InsertMessage
// @param: message(pojo.Message)
// @return: err(error)
func InsertMessage(message *pojo.Message) error {
	query := "INSERT INTO messages (to_user_id,from_user_id,content,create_time) VALUES (?,?,?,?)"
	_, err := Database.Exec(query, message.ToUserId, message.FromUserId, message.Content, message.CreateTime)
	return err
}
