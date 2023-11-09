package handlers

import (
	"strings"
	"time"

	"github.com/ectrc/snow/aid"
	p "github.com/ectrc/snow/person"
	"github.com/gofiber/fiber/v2"
)

var (
	oatuhTokenGrantTypes = map[string]func(c *fiber.Ctx, body *OAuthTokenBody) error{
		"client_credentials": PostOAuthTokenClientCredentials,
		"password": PostOAuthTokenPassword,
	}
)

type OAuthTokenBody struct {
	GrantType string `form:"grant_type" binding:"required"`
	Username string `form:"username"`
	Password string `form:"password"`
}

func PostOAuthToken(c *fiber.Ctx) error {
	var body OAuthTokenBody

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(aid.ErrorBadRequest("Invalid Request Body"))	
	}

	if action, ok := oatuhTokenGrantTypes[body.GrantType]; ok {
		return action(c, &body)
	}

	return c.Status(fiber.StatusBadRequest).JSON(aid.ErrorBadRequest("Invalid Grant Type"))
}

func PostOAuthTokenClientCredentials(c *fiber.Ctx, body *OAuthTokenBody) error {
	credentials, err := aid.JWTSign(aid.JSON{
		"snow_id": 0, // custom
		"creation_date": time.Now().Format("2006-01-02T15:04:05.999Z"),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(aid.ErrorInternalServer)
	}

	return c.Status(fiber.StatusOK).JSON(aid.JSON{
		"access_token": "eg1~"+credentials,
		"token_type": "bearer",
		"client_id": c.IP(),
		"client_service": "snow",
		"internal_client": true,
		"expires_in": 3600,
		"expires_at": time.Now().Add(time.Hour).Format("2006-01-02T15:04:05.999Z"),
	})
}

func PostOAuthTokenPassword(c *fiber.Ctx, body *OAuthTokenBody) error {
	if body.Username == "" || body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(aid.ErrorBadRequest("Username/Password is empty"))
	}

	person := p.FindByDisplay(strings.Split(body.Username, "@")[0])
	if person == nil {
		return c.Status(fiber.StatusBadRequest).JSON(aid.ErrorBadRequest("No Account Found"))
	}

	if person.AccessKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(aid.ErrorBadRequest("Activation Required"))
	}

	if person.AccessKey != body.Password {
		return c.Status(fiber.StatusBadRequest).JSON(aid.ErrorBadRequest("Invalid Access Key"))
	}

	access, err := aid.JWTSign(aid.JSON{
		"snow_id": person.ID, // custom
		"creation_date": time.Now().Format("2006-01-02T15:04:05.999Z"),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(aid.ErrorInternalServer)
	}

	refresh, err := aid.JWTSign(aid.JSON{
		"snow_id": person.ID,
		"creation_date": time.Now().Format("2006-01-02T15:04:05.999Z"),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(aid.ErrorInternalServer)
	}

	return c.Status(fiber.StatusOK).JSON(aid.JSON{
		"access_token": "eg1~"+access,
		"account_id": person.ID,
		"client_id": c.IP(),
		"client_service": "snow",
		"device_id": "default",
		"display_name": person.DisplayName,
		"expires_at": time.Now().Add(time.Hour * 24).Format("2006-01-02T15:04:05.999Z"),
		"expires_in": 86400,
		"internal_client": true,
		"refresh_expires": 86400,
		"refresh_expires_at": time.Now().Add(time.Hour * 24).Format("2006-01-02T15:04:05.999Z"),
		"refresh_token": "eg1~"+refresh,
		"token_type": "bearer",
	})
}

func GetOAuthVerify(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	if auth == "" {
		return c.Status(fiber.StatusForbidden).JSON(aid.ErrorBadRequest("Authorization Header is empty"))
	}
	real := strings.ReplaceAll(auth, "bearer eg1~", "")

	claims, err := aid.JWTVerify(real)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(aid.ErrorBadRequest("Invalid Access Token"))
	}

	if claims["snow_id"] == nil {
		return c.Status(fiber.StatusForbidden).JSON(aid.ErrorBadRequest("Invalid Access Token"))
	}

	snowId, ok := claims["snow_id"].(string)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(aid.ErrorBadRequest("Invalid Access Token"))
	}

	person := p.Find(snowId)
	if person == nil {
		return c.Status(fiber.StatusForbidden).JSON(aid.ErrorBadRequest("Invalid Access Token"))
	}

	return c.Status(fiber.StatusOK).JSON(aid.JSON{
		"app": "fortnite",
		"token": "eg1~"+real,
		"token_type": "bearer",
		"expires_at": time.Now().Add(time.Hour * 24).Format("2006-01-02T15:04:05.999Z"),
		"expires_in": 86400,
		"client_id": c.IP(),
		"session_id": "0",
		"device_id": "default",
		"internal_client": true,
		"client_service": "snow",
		"in_app_id": person.ID,
		"account_id": person.ID,
		"displayName": person.DisplayName,
	})
}

func OAuthMiddleware(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	if auth == "" {
		return c.Status(fiber.StatusForbidden).JSON(aid.ErrorBadRequest("Authorization Header is empty"))
	}
	real := strings.ReplaceAll(auth, "bearer eg1~", "")

	claims, err := aid.JWTVerify(real)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(aid.ErrorBadRequest("Invalid Access Token"))
	}

	if claims["snow_id"] == nil {
		return c.Status(fiber.StatusForbidden).JSON(aid.ErrorBadRequest("Invalid Access Token"))
	}

	snowId, ok := claims["snow_id"].(string)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(aid.ErrorBadRequest("Invalid Access Token"))
	}

	person := p.Find(snowId)
	if person == nil {
		return c.Status(fiber.StatusForbidden).JSON(aid.ErrorBadRequest("Invalid Access Token"))
	}

	c.Locals("person", person)

	return c.Next()
}

func DeleteOAuthSessions(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(aid.JSON{})
}

func GetPublicAccount(c *fiber.Ctx) error {
	person := p.Find(c.Params("accountId"))
	if person == nil {
		return c.Status(fiber.StatusBadRequest).JSON(aid.ErrorBadRequest("No Account Found"))
	}

	return c.Status(fiber.StatusOK).JSON(aid.JSON{
		"id": person.ID,
		"displayName": person.DisplayName,
		"externalAuths": []aid.JSON{},
	})
}

func GetPublicAccounts(c *fiber.Ctx) error {
	response := []aid.JSON{}

	accountIds := c.Request().URI().QueryArgs().PeekMulti("accountId")
	for _, accountIdSlice := range accountIds {
		person := p.Find(string(accountIdSlice))
		if person == nil {
			continue
		}

		response = append(response, aid.JSON{
			"id": person.ID,
			"displayName": person.DisplayName,
			"externalAuths": []aid.JSON{},
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func GetPublicAccountExternalAuths(c *fiber.Ctx) error {
	person := p.Find(c.Params("accountId"))
	if person == nil {
		return c.Status(fiber.StatusBadRequest).JSON(aid.ErrorBadRequest("No Account Found"))
	}

	return c.Status(fiber.StatusOK).JSON([]aid.JSON{})
}

func GetPublicAccountByDisplayName(c *fiber.Ctx) error {
	person := p.FindByDisplay(c.Params("displayName"))
	if person == nil {
		return c.Status(fiber.StatusBadRequest).JSON(aid.ErrorBadRequest("No Account Found"))
	}

	return c.Status(fiber.StatusOK).JSON(aid.JSON{
		"id": person.ID,
		"displayName": person.DisplayName,
		"externalAuths": []aid.JSON{},
	})
}