package model

type Post struct {
	ID       int64      `json:"id"`
	Title    string     `json:"title"`
	Caption  string     `json:"caption"`
	Body     string     `json:"body"`
	PostedBy string     `json:"postedBy"`
	Status   string     `json:"status"`
	Date     string     `json:"date"`
	Likes    int        `json:"likes"`
	Comments []*Comment `json:"comments"`
	IsLiked  bool       `json:"isLiked"`
}

type Comment struct {
	ID       int64      `json:"id"`
	Username string     `json:"username"`
	Comment  string     `json:"comment"`
	Date     string     `json:"date"`
	Replies  []*Comment `json:"replies"`
}

type ReadPostRequest struct {
	Followings bool `json:"followings"`
}

type CreatePostRequest struct {
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	Body     string `json:"body"`
	PostedBy int64  `json:"postedBy"`
	Status   string `json:"status"`
}

type UpdatePostRequest struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	Body     string `json:"body"`
	PostedBy int64  `json:"postedBy"`
	Status   string `json:"status"`
}

type DeletePostRequest struct {
	ID int64 `json:"id"`
}

type CommentOnPostRequest struct {
	PostID  int64  `json:"postId"`
	UserID  int64  `json:"userId"`
	Comment string `json:"comment"`
}

type CommentOnCommentRequest struct {
	CommentID int64  `json:"commentId"`
	PostID    int64  `json:"postId"`
	UserID    int64  `json:"userId"`
	Comment   string `json:"comment"`
}

type LikePostRequest struct {
	PostID int64 `json:"postId"`
	UserID int64 `json:"userId"`
}

type UnikePostRequest struct {
	PostID int64 `json:"postId"`
	UserID int64 `json:"userId"`
}
