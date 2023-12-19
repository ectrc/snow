package handlers

import (
	"strings"
	"time"

	"github.com/ectrc/snow/aid"
	p "github.com/ectrc/snow/person"
	"github.com/ectrc/snow/storage"
	"github.com/gofiber/fiber/v2"
)

var (
	oauthTokenGrantTypes = map[string]func(c *fiber.Ctx, body *FortniteTokenBody) error{
		"client_credentials": PostTokenClientCredentials,
		"password": PostTokenPassword,
	}
)

type FortniteTokenBody struct {
	GrantType string `form:"grant_type" binding:"required"`
	Username string `form:"username"`
	Password string `form:"password"`
}

func PostFortniteToken(c *fiber.Ctx) error {
	var body FortniteTokenBody

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(aid.ErrorBadRequest("Invalid Request Body"))	
	}

	if action, ok := oauthTokenGrantTypes[body.GrantType]; ok {
		return action(c, &body)
	}

	return c.Status(fiber.StatusBadRequest).JSON(aid.ErrorBadRequest("Invalid Grant Type"))
}

func PostTokenClientCredentials(c *fiber.Ctx, body *FortniteTokenBody) error {
	client, sig := aid.KeyPair.EncryptAndSignB64([]byte(c.IP()))
	hash := aid.Hash([]byte(client + "." + sig))

	return c.Status(fiber.StatusOK).JSON(aid.JSON{
		"access_token": hash,
		"token_type": "bearer",
		"client_id": c.IP(),
		"client_service": "fortnite",
		"internal_client": true,
		"expires_in": 3600,
		"expires_at": time.Now().Add(time.Hour).Format("2006-01-02T15:04:05.999Z"),
		"product_id": "prod-fn",
		"sandbox_id": "fn",
	})
}

func PostTokenPassword(c *fiber.Ctx, body *FortniteTokenBody) error {
	if body.Username == "" || body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(aid.ErrorBadRequest("Username/Password is empty"))
	}

	person := p.FindByDisplay(strings.Split(body.Username, "@")[0])
	if person == nil {
		return c.Status(fiber.StatusBadRequest).JSON(aid.ErrorBadRequest("No Account Found"))
	}

	access, ac_sig := aid.KeyPair.EncryptAndSignB64([]byte(person.ID))
	ac_hash := aid.Hash([]byte(access + "." + ac_sig))

	ac_token := &storage.DB_GameToken{
		ID: ac_hash,
		PersonID: person.ID,
		AccessToken: access + "." + ac_sig,
		Type: "access",
	}
	storage.Repo.SaveToken(ac_token)

	refresh, re_sig := aid.KeyPair.EncryptAndSignB64([]byte(person.ID))
	re_hash := aid.Hash([]byte(refresh + "." + re_sig))

	re_token := &storage.DB_GameToken{
		ID: re_hash,
		PersonID: person.ID,
		AccessToken: refresh + "." + re_sig,
		Type: "refresh",
	}
	storage.Repo.SaveToken(re_token)

	return c.Status(fiber.StatusOK).JSON(aid.JSON{
		// "access_token": access + "." + ac_sig,
		"access_token": ac_hash,
		"account_id": person.ID,
		"client_id": c.IP(),
		"client_service": "fortnite",
		"app": "fortnite",
		"device_id": "default",
		"display_name": person.DisplayName,
		"expires_at": time.Now().Add(time.Hour * 24).Format("2006-01-02T15:04:05.999Z"),
		"expires_in": 86400,
		"internal_client": true,
		"refresh_expires": 86400,
		"refresh_expires_at": time.Now().Add(time.Hour * 24).Format("2006-01-02T15:04:05.999Z"),
		// "refresh_token": refresh + "." + re_sig,
		"refresh_token": re_hash,
		"token_type": "bearer",
		"product_id": "prod-fn",
		"sandbox_id": "fn",
	})
}

func GetOAuthVerify(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	if auth == "" {
		return c.Status(fiber.StatusForbidden).JSON(aid.ErrorBadRequest("Authorization Header is empty"))
	}
	real := strings.ReplaceAll(auth, "bearer ", "")

	found := storage.Repo.GetToken(real)
	if found == nil {
		return c.Status(fiber.StatusForbidden).JSON(aid.ErrorBadRequest("Invalid Access Token"))
	}
	snowId := found.PersonID

	person := p.Find(snowId)
	if person == nil {
		return c.Status(fiber.StatusForbidden).JSON(aid.ErrorBadRequest("Invalid Access Token"))
	}

	return c.Status(fiber.StatusOK).JSON(aid.JSON{
		"app": "fortnite",
		"token": real,
		"token_type": "bearer",
		"expires_at": time.Now().Add(time.Hour * 24).Format("2006-01-02T15:04:05.999Z"),
		"expires_in": 86400,
		"client_id": c.IP(),
		"session_id": "0",
		"device_id": "default",
		"internal_client": true,
		"client_service": "fortnite",
		"in_app_id": person.ID,
		"account_id": person.ID,
		"displayName": person.DisplayName,
		"product_id": "prod-fn",
		"sandbox_id": "fn",	
	})
}

func MiddlewareFortnite(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	if auth == "" {
		return c.Status(fiber.StatusForbidden).JSON(aid.ErrorBadRequest("Authorization Header is empty"))
	}
	real := strings.ReplaceAll(auth, "bearer ", "")

	found := storage.Repo.GetToken(real)
	if found == nil {
		return c.Status(fiber.StatusForbidden).JSON(aid.ErrorBadRequest("Invalid Access Token"))
	}
	snowId := found.PersonID

	person := p.Find(snowId)
	if person == nil {
		return c.Status(fiber.StatusForbidden).JSON(aid.ErrorBadRequest("Invalid Access Token"))
	}

	c.Locals("person", person)
	return c.Next()
}

func MiddlewareWeb(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	if auth == "" {
		return c.Status(fiber.StatusForbidden).JSON(aid.ErrorBadRequest("Authorization Header is empty"))
	}

	found := storage.Repo.GetToken(auth)
	if found == nil {
		return c.Status(fiber.StatusForbidden).JSON(aid.ErrorBadRequest("Invalid Access Token"))
	}
	snowId := found.PersonID

	person := p.Find(snowId)
	if person == nil {
		return c.Status(fiber.StatusForbidden).JSON(aid.ErrorBadRequest("Invalid Access Token"))
	}

	c.Locals("person", person)
	return c.Next()
}

func DeleteToken(c *fiber.Ctx) error {
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