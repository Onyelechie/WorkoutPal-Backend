package model

import "workoutpal/src/util/constants"

type Error struct {
	Type     string                `json:"type"`
	Status   int                   `json:"status"`
	Detail   constants.UserMessage `json:"detail"`
	Instance string                `json:"instance"`
	Error    string                `json:"error"`
}
