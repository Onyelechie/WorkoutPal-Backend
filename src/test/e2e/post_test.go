package e2e

import (
	"net/http"
	"testing"
)

type post struct {
	ID       int64      `json:"id"`
	Title    string     `json:"title"`
	Caption  string     `json:"caption"`
	Body     string     `json:"body"`
	PostedBy string     `json:"postedBy"`
	Status   string     `json:"status"`
	Date     string     `json:"date"`
	Likes    int        `json:"likes"`
	Comments []*comment `json:"comments"`
	IsLiked  bool       `json:"isLiked"`
}

type comment struct {
	ID       int64      `json:"id"`
	Username string     `json:"username"`
	Comment  string     `json:"comment"`
	Date     string     `json:"date"`
	Replies  []*comment `json:"replies"`
}

type createPostReq struct {
	Title   string `json:"title"`
	Caption string `json:"caption"`
	Body    string `json:"body"`
	Status  string `json:"status"`
}

type commentOnPostReq struct {
	PostID  int64  `json:"postId"`
	Comment string `json:"comment"`
}

type commentOnCommentReq struct {
	CommentID int64  `json:"commentId"`
	PostID    int64  `json:"postId"`
	Comment   string `json:"comment"`
}

type likePostReq struct {
	PostID int64 `json:"postId"`
}

type unlikePostReq struct {
	PostID int64 `json:"postId"`
}

func testEndToEnd_Posts_Create(t *testing.T) {
	body := createPostReq{
		Title:   "E2E Post " + randStringAlphaNum(6),
		Caption: "caption",
		Body:    "body text",
		Status:  "active",
	}

	resp := doRequest(t, http.MethodPost, "/posts", body, nil)
	defer resp.Body.Close()

	mustStatus(t, resp, http.StatusCreated)

	created := mustDecode[post](t, resp)
	if created.ID == 0 {
		t.Fatalf("expected non-zero post id")
	}
	if created.Title != body.Title {
		t.Fatalf("expected title=%q got=%q", body.Title, created.Title)
	}
}

func testEndToEnd_Posts_List(t *testing.T) {
	resp := doRequest(t, http.MethodGet, "/posts", nil, nil)
	defer resp.Body.Close()

	mustStatus(t, resp, http.StatusOK)

	list := mustDecode[[]post](t, resp)
	if len(list) == 0 {
		t.Fatalf("expected at least one post")
	}
	if list[0].ID == 0 {
		t.Fatalf("expected post to have ID")
	}
	if list[0].Title == "" {
		t.Fatalf("expected post to have Title")
	}
}

func testEndToEnd_Posts_CommentOnPost(t *testing.T) {
	createBody := createPostReq{
		Title:   "E2E Comment " + randStringAlphaNum(6),
		Caption: "caption",
		Body:    "body",
		Status:  "active",
	}

	createResp := doRequest(t, http.MethodPost, "/posts", createBody, nil)
	defer createResp.Body.Close()
	mustStatus(t, createResp, http.StatusCreated)

	created := mustDecode[post](t, createResp)
	if created.ID == 0 {
		t.Fatalf("expected created post id != 0")
	}

	commentBody := commentOnPostReq{
		PostID:  created.ID,
		Comment: "nice post",
	}

	resp := doRequest(t, http.MethodPost, "/posts/comment", commentBody, nil)
	defer resp.Body.Close()
	mustStatus(t, resp, http.StatusOK)

	msg := mustDecode[basicResponse](t, resp)
	if msg.Message == "" {
		t.Fatalf("expected success message")
	}
}

func testEndToEnd_Posts_CommentOnComment(t *testing.T) {
	createBody := createPostReq{
		Title:   "E2E Reply " + randStringAlphaNum(6),
		Caption: "caption",
		Body:    "body",
		Status:  "active",
	}

	createResp := doRequest(t, http.MethodPost, "/posts", createBody, nil)
	defer createResp.Body.Close()
	mustStatus(t, createResp, http.StatusCreated)

	created := mustDecode[post](t, createResp)
	if created.ID == 0 {
		t.Fatalf("expected created post id != 0")
	}

	commentBody := commentOnPostReq{
		PostID:  created.ID,
		Comment: "first comment",
	}
	commentResp := doRequest(t, http.MethodPost, "/posts/comment", commentBody, nil)
	defer commentResp.Body.Close()
	mustStatus(t, commentResp, http.StatusOK)

	replyBody := commentOnCommentReq{
		CommentID: 1,
		PostID:    created.ID,
		Comment:   "reply comment",
	}
	replyResp := doRequest(t, http.MethodPost, "/posts/comment/reply", replyBody, nil)
	defer replyResp.Body.Close()
	mustStatus(t, replyResp, http.StatusOK)

	msg := mustDecode[basicResponse](t, replyResp)
	if msg.Message == "" {
		t.Fatalf("expected success message")
	}
}

func testEndToEnd_Posts_Delete(t *testing.T) {
	createBody := createPostReq{
		Title:   "E2E Delete " + randStringAlphaNum(6),
		Caption: "caption",
		Body:    "body",
		Status:  "active",
	}

	createResp := doRequest(t, http.MethodPost, "/posts", createBody, nil)
	defer createResp.Body.Close()
	mustStatus(t, createResp, http.StatusCreated)

	created := mustDecode[post](t, createResp)
	if created.ID == 0 {
		t.Fatalf("expected created post id != 0")
	}

	delResp := doRequest(t, http.MethodDelete, "/posts/"+int64ToStr(created.ID), created, nil)
	defer delResp.Body.Close()
	mustStatus(t, delResp, http.StatusOK)

	msg := mustDecode[basicResponse](t, delResp)
	if msg.Message == "" {
		t.Fatalf("expected success message")
	}
}

func testEndToEnd_Posts_Like_Unlike(t *testing.T) {
	createBody := createPostReq{
		Title:   "E2E Like " + randStringAlphaNum(6),
		Caption: "caption",
		Body:    "body",
		Status:  "active",
	}

	createResp := doRequest(t, http.MethodPost, "/posts", createBody, nil)
	defer createResp.Body.Close()
	mustStatus(t, createResp, http.StatusCreated)

	created := mustDecode[post](t, createResp)
	if created.ID == 0 {
		t.Fatalf("expected created post id != 0")
	}

	likeBody := likePostReq{PostID: created.ID}
	likeResp := doRequest(t, http.MethodPost, "/posts/like", likeBody, nil)
	defer likeResp.Body.Close()
	mustStatus(t, likeResp, http.StatusOK)

	likeMsg := mustDecode[basicResponse](t, likeResp)
	if likeMsg.Message == "" {
		t.Fatalf("expected success message on like")
	}

	listResp := doRequest(t, http.MethodGet, "/posts", nil, nil)
	defer listResp.Body.Close()
	mustStatus(t, listResp, http.StatusOK)
	list := mustDecode[[]post](t, listResp)
	if len(list) == 0 {
		t.Fatalf("expected at least one post after like")
	}

	for _, p := range list {
		if p.ID == created.ID {
			if !p.IsLiked {
				t.Fatalf("expected post %d to be liked", created.ID)
			}
			break
		}
	}

	unlikeBody := unlikePostReq{PostID: created.ID}
	unlikeResp := doRequest(t, http.MethodPost, "/posts/unlike", unlikeBody, nil)
	defer unlikeResp.Body.Close()
	mustStatus(t, unlikeResp, http.StatusOK)

	unlikeMsg := mustDecode[basicResponse](t, unlikeResp)
	if unlikeMsg.Message == "" {
		t.Fatalf("expected success message on unlike")
	}
}
