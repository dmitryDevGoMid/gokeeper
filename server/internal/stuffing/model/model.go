package model

import "github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/repository/user"

type RequestBody struct {
	Body []byte    `json:"body"`
	User user.User `json:"user"`
}
