package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"
)

var (
	matrixClient *mautrix.Client
	defaultRoom  *string
	port         = "3000"
)

func init() {
	homeserverURL, ok := os.LookupEnv("HOMESERVER_URL")
	if !ok {
		panic("HOMESERVER_URL env var required")
	}
	userID, ok := os.LookupEnv("USER_ID")
	if !ok {
		panic("USER_ID env var required")
	} else {
		userID = strings.Replace(userID, homeserverURL, "", -1)
		userID = strings.TrimPrefix(userID, "@")
	}
	accessToken, ok := os.LookupEnv("ACCESS_TOKEN")
	if !ok {
		panic("ACCESS_TOKEN env var required")
	}

	if p, ok := os.LookupEnv("PORT"); ok {
		port = p
	}

	room, ok := os.LookupEnv("DEFAULT_ROOM")
	if ok {
		defaultRoom = &room
	}

	client, err := mautrix.NewClient(homeserverURL, id.NewUserID(userID, homeserverURL), accessToken)
	if err != nil {
		panic(err)
	}
	matrixClient = client
}

func main() {
	app := fiber.New()
	app.Post("/", HandlePayloadPost)
	app.Listen(fmt.Sprintf(":%s", port))
}

func HandlePayloadPost(c *fiber.Ctx) error {
	payload := Payload{}

	if err := c.BodyParser(&payload); err != nil {
		return err
	}
	if err := payload.Validate(defaultRoom); err != nil {
		fmt.Println("Invalid payload - ", err)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	switch payload.Type {
	case PayloadTypeText:
		_, err := matrixClient.SendText(id.RoomID(payload.RoomID), payload.Message)
		if err != nil {
			if httpErr, ok := err.(mautrix.HTTPError); ok {
				return c.Status(httpErr.Response.StatusCode).SendString(httpErr.RespError.Err)
			}
			return fiber.ErrInternalServerError
		}
	case PayloadTypeNotice:
		_, err := matrixClient.SendNotice(id.RoomID(payload.RoomID), payload.Message)
		if err != nil {
			if httpErr, ok := err.(mautrix.HTTPError); ok {
				return c.Status(httpErr.Response.StatusCode).SendString(httpErr.RespError.Err)
			}
			return fiber.ErrInternalServerError
		}
	}

	return c.SendStatus(fiber.StatusOK)
}
