package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/alertmanager/template"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"
)

func HandleAlertmanagerPayloadPost(c *fiber.Ctx) error {
	payload := template.Data{}
	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	for _, alert := range payload.Alerts {
		message := ""
		if alert.Status == "firing" {
			switch alert.Labels["severity"] {
			case "warning":
				message = fmt.Sprintf("⚠️ %s", alert.Annotations["description"])
			case "notify":
				message = fmt.Sprintf("@room - %s", alert.Annotations["description"])
			}
		} else {
			message = fmt.Sprintf("☑️ %s", alert.Annotations["description"])
		}

		_, err := matrixClient.SendText(id.RoomID(*defaultRoom), message)
		if err != nil {
			if httpErr, ok := err.(mautrix.HTTPError); ok {
				return c.Status(httpErr.Response.StatusCode).SendString(httpErr.RespError.Err)
			}
			return fiber.ErrInternalServerError
		}
	}

	return c.SendStatus(fiber.StatusOK)
}
