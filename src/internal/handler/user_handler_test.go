package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"workoutpal/src/internal/model"
	"workoutpal/src/util/constants"

	mock_service "workoutpal/src/mock_internal/domain/service"

	"github.com/golang/mock/gomock"
)

func withIDCtx(r *http.Request, id int64) *http.Request {
	ctx := context.WithValue(r.Context(), constants.ID_KEY, id)
	return r.WithContext(ctx)
}

func TestUserHandler_CreateNewUser_BadJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockUserService(ctrl)
	h := &userHandler{userService: svc}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString("{"))
	r.Header.Set("Content-Type", "application/json")

	h.CreateNewUser(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestUserHandler_CreateNewUser_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockUserService(ctrl)
	h := &userHandler{userService: svc}

	req := model.CreateUserRequest{
		Username: "max",
		Email:    "max@example.com",
		Name:     "Max",
		Password: "Str0ng!Pass",
	}
	svc.EXPECT().
		CreateUser(gomock.AssignableToTypeOf(model.CreateUserRequest{})).
		Return((*model.User)(nil), errors.New("username taken"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users", mustJSON(t, req))
	r.Header.Set("Content-Type", "application/json")

	h.CreateNewUser(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
}

func TestUserHandler_CreateNewUser_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockUserService(ctrl)
	h := &userHandler{userService: svc}

	req := model.CreateUserRequest{
		Username: "max",
		Email:    "max@example.com",
		Name:     "Max",
		Password: "Str0ng!Pass",
	}
	want := &model.User{ID: 1, Username: req.Username, Email: req.Email, Name: req.Name}

	svc.EXPECT().
		CreateUser(gomock.AssignableToTypeOf(model.CreateUserRequest{})).
		Return(want, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users", mustJSON(t, req))
	r.Header.Set("Content-Type", "application/json")

	h.CreateNewUser(w, r)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201", w.Code)
	}
	var got model.User
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID != want.ID || got.Username != want.Username || got.Email != want.Email {
		t.Fatalf("unexpected user: %+v", got)
	}
}

func TestUserHandler_ReadAllUsers_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockUserService(ctrl)
	h := &userHandler{userService: svc}

	svc.EXPECT().ReadUsers().Return(nil, errors.New("db down"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users", nil)

	h.ReadAllUsers(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
}

func TestUserHandler_ReadAllUsers_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockUserService(ctrl)
	h := &userHandler{userService: svc}

	want := []*model.User{
		{ID: 1, Username: "max"},
		{ID: 2, Username: "sam"},
	}
	svc.EXPECT().ReadUsers().Return(want, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users", nil)

	h.ReadAllUsers(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var got []model.User
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got) != len(want) || got[0].ID != want[0].ID {
		t.Fatalf("unexpected users: %+v", got)
	}
}

func TestUserHandler_ReadUserByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockUserService(ctrl)
	h := &userHandler{userService: svc}

	const id int64 = 42
	svc.EXPECT().ReadUserByID(id).Return((*model.User)(nil), errors.New("not found"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users/42", nil)
	r = withIDCtx(r, id)

	h.ReadUserByID(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 404", w.Code)
	}
}

func TestUserHandler_ReadUserByID_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockUserService(ctrl)
	h := &userHandler{userService: svc}

	const id int64 = 7
	want := &model.User{ID: id, Username: "max"}
	svc.EXPECT().ReadUserByID(id).Return(want, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users/7", nil)
	r = withIDCtx(r, id)

	h.ReadUserByID(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var got model.User
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID != want.ID || got.Username != want.Username {
		t.Fatalf("unexpected user: %+v", got)
	}
}

func TestUserHandler_UpdateUser_BadJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockUserService(ctrl)
	h := &userHandler{userService: svc}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/users/9", bytes.NewBufferString("{"))
	r = withIDCtx(r, 9)
	r.Header.Set("Content-Type", "application/json")

	h.UpdateUser(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestUserHandler_UpdateUser_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockUserService(ctrl)
	h := &userHandler{userService: svc}

	const id int64 = 9
	req := model.UpdateUserRequest{Username: "newname"}
	svc.EXPECT().
		UpdateUser(gomock.AssignableToTypeOf(model.UpdateUserRequest{})).
		Return((*model.User)(nil), errors.New("bad update"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/users/9", mustJSON(t, req))
	r = withIDCtx(r, id)
	r.Header.Set("Content-Type", "application/json")

	h.UpdateUser(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
}

func TestUserHandler_UpdateUser_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockUserService(ctrl)
	h := &userHandler{userService: svc}

	const id int64 = 9
	req := model.UpdateUserRequest{Username: "newname"}
	want := &model.User{ID: id, Username: req.Username}

	svc.EXPECT().
		UpdateUser(gomock.AssignableToTypeOf(model.UpdateUserRequest{})).
		DoAndReturn(func(got model.UpdateUserRequest) (*model.User, error) {
			if got.ID != id {
				return nil, errors.New("missing id propagation")
			}
			return want, nil
		})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/users/9", mustJSON(t, req))
	r = withIDCtx(r, id)
	r.Header.Set("Content-Type", "application/json")

	h.UpdateUser(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var got model.User
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID != want.ID || got.Username != want.Username {
		t.Fatalf("unexpected user: %+v", got)
	}
}

func TestUserHandler_DeleteUser_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockUserService(ctrl)
	h := &userHandler{userService: svc}

	const id int64 = 13
	svc.EXPECT().DeleteUser(model.DeleteUserRequest{ID: id}).Return(errors.New("not found"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/users/13", nil)
	r = withIDCtx(r, id)

	h.DeleteUser(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 404", w.Code)
	}
}

func TestUserHandler_DeleteUser_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockUserService(ctrl)
	h := &userHandler{userService: svc}

	const id int64 = 13
	svc.EXPECT().DeleteUser(model.DeleteUserRequest{ID: id}).Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/users/13", nil)
	r = withIDCtx(r, id)

	h.DeleteUser(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var br model.BasicResponse
	if err := json.NewDecoder(w.Body).Decode(&br); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if br.Message != "User deleted successfully" {
		t.Fatalf("message = %q, want %q", br.Message, "User deleted successfully")
	}
}
