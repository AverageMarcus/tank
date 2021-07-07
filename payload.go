package main

import (
	"errors"
	"fmt"
	"strings"

	"maunium.net/go/mautrix/id"
)

type PayloadType string

const (
	PayloadTypeText   PayloadType = "text"
	PayloadTypeNotice PayloadType = "notice"
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

	room := parseRoom(p.RoomID)
	if strings.HasPrefix(room, "#") || !strings.HasPrefix("!", room) {
		if !strings.HasPrefix(room, "#") {
			room = fmt.Sprintf("#%s", room)
		}

		resp, err := matrixClient.ResolveAlias(id.RoomAlias(room))
		if err == nil {
			p.RoomID = resp.RoomID.String()
		} else {
			p.RoomID = strings.Replace(room, "#", "!", 1)
		}
	}

	return nil
}

func parseRoom(room string) string {
	prefix := ""
	local := ""
	domain := ""

	parts := strings.Split(room, ":")
	if len(parts) == 2 {
		domain = parts[1]
	} else {
		domain = matrixClient.HomeserverURL.Host
	}

	if strings.HasPrefix(parts[0], "!") {
		prefix = "!"
		parts[0] = strings.TrimPrefix(parts[0], "!")
	}
	if strings.HasPrefix(parts[0], "#") {
		prefix = "#"
		parts[0] = strings.TrimPrefix(parts[0], "#")
	}

	local = parts[0]

	return fmt.Sprintf("%s%s:%s", prefix, local, domain)
}
