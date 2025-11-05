package util

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
	"workoutpal/src/internal/model"
	"workoutpal/src/util/constants"

	"github.com/go-chi/render"
	"github.com/lib/pq"
)

func Error(err error, instance string) *model.Error {
	if err == nil {
		return nil
	}

	userMsg, status := msgAndStatus(err)

	return &model.Error{
		Type:     "ERROR",
		Status:   status,
		Detail:   userMsg,
		Instance: instance,
		Error:    err.Error(),
	}
}

func extractFields(err error) map[string]any {
	m := make(map[string]any)
	return m
}

func msgAndStatus(err error) (constants.UserMessage, int) {
	if err == nil {
		return constants.DEFAULT, http.StatusOK
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return fromPqError(pqErr)
	}

	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		field := typeErr.Field
		if field == "" {
			field = "(unknown field)"
		}
		return constants.UserMessage(
			"Invalid value for field '" + field +
				"' — expected " + typeErr.Type.String() +
				" but got " + typeErr.Value + ".",
		), http.StatusUnprocessableEntity
	}

	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return constants.UserMessage(
			"Malformed JSON near position " + fmt.Sprint(syntaxErr.Offset) + ".",
		), http.StatusBadRequest
	}

	if strings.HasPrefix(err.Error(), "json: unknown field ") {
		name := strings.TrimPrefix(err.Error(), "json: unknown field ")
		return constants.UserMessage("Unknown field " + name + "."), http.StatusBadRequest
	}

	if errors.Is(err, io.ErrUnexpectedEOF) {
		return constants.UserMessage("Incomplete or truncated JSON input."), http.StatusBadRequest
	}

	var tperr *time.ParseError
	if errors.As(err, &tperr) {
		return constants.INVALID_FORMAT, http.StatusUnprocessableEntity
	}

	// Check for specific error messages from repository layer
	if strings.Contains(err.Error(), "user already exists") {
		return constants.DUPLICATE, http.StatusBadRequest
	}

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return constants.NOT_FOUND, http.StatusInternalServerError
	case errors.Is(err, context.DeadlineExceeded):
		return constants.TIMEOUT, http.StatusGatewayTimeout
	default:
		return constants.UNKNOWN, http.StatusInternalServerError
	}
}

func fromPqError(e *pq.Error) (constants.UserMessage, int) {
	if e == nil {
		return constants.UNKNOWN_DB, http.StatusInternalServerError
	}

	code := string(e.Code)
	switch code {
	case "23503":
		text := strings.ToLower(e.Detail + " " + e.Message)

		if strings.Contains(text, "is not present in table") ||
			strings.HasPrefix(strings.ToLower(e.Message), "insert or update") {
			col, ref := extractFkParts(e.Detail)
			if col != "" && ref != "" {
				return constants.UserMessage("This record references another that doesn’t exist."), http.StatusUnprocessableEntity
			}
			return constants.UserMessage("This record references another that doesn’t exist."), http.StatusUnprocessableEntity
		}

		if strings.Contains(text, "is still referenced from table") ||
			strings.HasPrefix(strings.ToLower(e.Message), "update or delete") {
			return constants.FOREIGN_KEY, http.StatusConflict
		}

		return constants.FOREIGN_KEY, http.StatusConflict
	default:
		return fromPgCode(code)
	}
}

func fromPgCode(code string) (constants.UserMessage, int) {
	switch code {
	case "23505":
		return constants.DUPLICATE, http.StatusBadRequest
	case "23502":
		return constants.MISSING_FIELD, http.StatusUnprocessableEntity
	case "23514":
		return constants.CHECK_VIOLATION, http.StatusUnprocessableEntity
	case "23P01":
		return constants.CONFLICT, http.StatusConflict
	case "23503":
		return constants.FOREIGN_KEY, http.StatusConflict
	case "22P02":
		return constants.INVALID_FORMAT, http.StatusUnprocessableEntity
	case "22001":
		return constants.TOO_LONG, http.StatusUnprocessableEntity
	case "22003":
		return constants.OUT_OF_RANGE, http.StatusUnprocessableEntity
	case "22007":
		return constants.INVALID_FORMAT, http.StatusUnprocessableEntity
	case "40001", "40P01", "55P03":
		return constants.CONFLICT, http.StatusConflict
	case "08006", "08001", "53300", "57P01", "57P02":
		return constants.UNAVAILABLE, http.StatusServiceUnavailable

	default:
		return constants.UNKNOWN_DB, http.StatusInternalServerError
	}
}

var reFkNotPresent = regexp.MustCompile(`Key \((?P<col>[^)]+)\)=\([^)]+\) is not present in table "(?P<reftable>[^"]+)"`)

func extractFkParts(detail string) (col, ref string) {
	if m := reFkNotPresent.FindStringSubmatch(detail); m != nil {
		for i, name := range reFkNotPresent.SubexpNames() {
			if name == "col" {
				col = m[i]
			}
			if name == "reftable" {
				ref = m[i]
			}
		}
	}
	return
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, err *model.Error) {
	w.WriteHeader(err.Status)
	render.JSON(w, r, err)
}

func ErrorResponseWithStatus(w http.ResponseWriter, r *http.Request, err *model.Error, status int) {
	w.WriteHeader(status)
	render.JSON(w, r, err)
}
