package main

import (
	"bytes"
	"fmt"
	tTemplate "text/template"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/alertmanager/template"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/format"
	"maunium.net/go/mautrix/id"
)

func HandleAlertmanagerPayloadPost(c *fiber.Ctx) error {
	payload := template.Data{}
	if err := c.BodyParser(&payload); err != nil {
		fmt.Println("Failed to parse payload", err)
		return err
	}

	fmt.Println("Got alertmanager payload")

	for _, alert := range payload.Alerts {
		message := ""

		var rendered bytes.Buffer
		at, _ := tTemplate.New("AlertMessage").Parse(alert.Annotations["description"])
		at.Execute(&rendered, alert)

		if alert.Status == "firing" {
			switch alert.Labels["severity"] {
			case "warning":
				message = fmt.Sprintf("⚠️ %s", rendered.String())
			case "notify":
				message = fmt.Sprintf("🔥 @room - %s", rendered.String())
			}
		} else {
			message = fmt.Sprintf("✅ %s", rendered.String())
		}

		_, err := matrixClient.SendMessageEvent(
			id.RoomID(getRoom(c.Query("room", *defaultRoom))),
			event.EventMessage,
			format.RenderMarkdown(message, true, true),
		)
		if err != nil {
			fmt.Println("Failed sending to Matrix", err)
			if httpErr, ok := err.(mautrix.HTTPError); ok {
				return c.Status(httpErr.Response.StatusCode).SendString(httpErr.RespError.Err)
			}
			return fiber.ErrInternalServerError
		}
	}

	return c.SendStatus(fiber.StatusOK)
}
