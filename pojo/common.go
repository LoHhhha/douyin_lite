package pojo

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Video struct {
	Vid           int64  `json:"vid,omitempty"`
	Uid           int64  `json:"uid"`
	PlayURL       string `json:"play_url"`
	CoverURL      string `json:"cover_url"`
	FavoriteCount int64  `json:"favorite_count"`
	CommentCount  int64  `json:"comment_count"`
	Uploadtime    string `json:"uploadtime,omitempty"`
	Title         string `json:"title,omitempty"`
}

type User struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	Password        string `json:"password,omitempty"`
	FollowCount     int64  `json:"follow_count,omitempty"`
	FollowerCount   int64  `json:"follower_count,omitempty"`
	IsFollow        bool   `json:"is_follow,omitempty"`
	TotalFavorited  int64  `json:"total_favorited,omitempty"`
	Avatar          string `json:"avatar,omitempty"`
	BackgroundImage string `json:"background_image,omitempty"`
	Signature       string `json:"signature,omitempty"`
	WorkCount       int    `json:"work_count,omitempty"`
	FavoriteCount   int    `json:"favorite_count,omitempty"`
}
type Message struct {
	Id         int64  `json:"id"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	FromUserId int64  `json:"from_user_id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateTime string `json:"create_time"`
}

type MessageSendEvent struct {
	UserId     int64  `json:"user_id,omitempty"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

type MessagePushEvent struct {
	FromUserId int64  `json:"user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

type VideoWithInfo struct {
	ID            int64  `json:"id"`
	Author        User   `json:"author"`
	PlayURL       string `json:"play_url"`
	CoverURL      string `json:"cover_url"`
	FavoriteCount int64  `json:"favorite_count"`
	CommentCount  int64  `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
	Title         string `json:"title"`
}

type Comment struct {
	ID         int64  `json:"id"`
	VID        int64  `json:"vid"`
	UID        int64  `json:"uid"`
	Content    string `json:"content"`
	CreateDate string `json:"create_date"`
}

type CommentWithUser struct {
	ID         int64  `json:"id"`
	User       User   `json:"user"`
	Content    string `json:"content"`
	CreateDate string `json:"create_date"`
}
