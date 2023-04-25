package main

import (
	"errors"
)

type PayloadType string

const (
	PayloadTypeText     PayloadType = "text"
	PayloadTypeMarkdown PayloadType = "markdown"
	PayloadTypeNotice   PayloadType = "notice"
)

type Payload struct {
	Type    PayloadType `json:"type"`
	Message string      `json:"message"`
	RoomID  string      `json:"roomID"`
}

func (p *Payload) Validate(defaultRoom *string) error {
	if p.Type == "" {
		p.Type = PayloadTypeText
	}

	if p.RoomID == "" && defaultRoom != nil {
		p.RoomID = *defaultRoom
	}

	if p.Message == "" {
		return errors.New("'message' is required")
	}

	if p.RoomID == "" {
		return errors.New("'roomID' is required")
	}

	p.RoomID = getRoom(p.RoomID)

	return nil
}
