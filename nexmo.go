package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/format"
	"maunium.net/go/mautrix/id"
)

type NexmoSMS struct {
	ID   string `json:"messageId"`
	From string `json:"msisdn"`
	Text string `json:"test"`
	Type string `json:"type"`
}

func HandleNexmoSMSPost(c *fiber.Ctx) error {
	payload := NexmoSMS{}
	if err := c.BodyParser(&payload); err != nil {
		fmt.Println("Failed to parse payload", err)
		return err
	}

	msg := fmt.Sprintf(`
New Message from: %s

> %s`, payload.From, payload.Text)

	_, err := matrixClient.SendMessageEvent(
		id.RoomID(getRoom("SMS")),
		event.EventMessage,
		format.RenderMarkdown(msg, true, true),
	)
	if err != nil {
		fmt.Println("Failed sending to Matrix", err)
		if httpErr, ok := err.(mautrix.HTTPError); ok {
			return c.Status(httpErr.Response.StatusCode).SendString(httpErr.RespError.Err)
		}
		return fiber.ErrInternalServerError
	}

	return c.SendStatus(fiber.StatusOK)
}
