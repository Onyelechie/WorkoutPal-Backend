package model

type Post struct {
	ID       int64     `json:"id"`
	PostedBy string    `json:"postedBy"`
	Title    string    `json:"title"`
	Caption  string    `json:"caption"`
	Date     string    `json:"date"`
	Content  string    `json:"content"`
	Likes    int       `json:"likes"`
	Comments []Comment `json:"comments"`
}

type Comment struct {
	CommentedBy string `json:"commentedBy"`
	Comment     string `json:"comment"`
	Date        string `json:"date"`
}

type ReadPostRequest struct {
	Followings bool `json:"followings"`
}

type CreatePostRequest struct {
	Title   string `json:"title"`
	Caption string `json:"caption"`
	Content string `json:"content"`
}

type CommentOnPostRequest struct {
	PostID  int64  `json:"postId"`
	Comment string `json:"comment"`
}

type LikePostRequest struct {
	PostID int `json:"postID"`
}
