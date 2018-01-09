package api

import (
	"thaThrowdown/common/database/dgraph"
)

// J is a wrapper of string, interface for serializing JSON
type J map[string]interface{}

// SimpleMessage is the struct encapsulating the message response
type SimpleMessage struct {
	IsSuccess bool   `json:"success"`
	Message   string `json:"message"`
}

// Fail returns a fail JSON document
func Fail(msg string) *J {
	return &J{"success": false, "message": msg}
}

// Success returns a successful JSON document
func Success(msg string) *J {
	return &J{"success": true, "message": msg}
}

// SuccessWithPayload returns a successful JSON document with payload
func SuccessWithPayload(msg string, payload interface{}) *J {
	return &J{"success": true, "message": msg, "data": payload}
}

// SuccessWithID returns a successful JSON document with Id
func SuccessWithID(msg string, id dgraph.UID) *J {
	return &J{"success": true, "message": msg, "id": id.ToHex()}
}
