package handlers

import (
	"strings"
	"time"

	"github.com/ectrc/snow/aid"
	p "github.com/ectrc/snow/person"
	"github.com/gofiber/fiber/v2"
)

var (
	oauthTokenGrantTypes = map[string]func(c *fiber.Ctx, body *FortniteTokenBody) error{
		"client_credentials": PostTokenClientCredentials, // spams the api?? like wtf
		"password": PostTokenPassword,
    "exchange_code": PostTokenExchangeCode,
	}
)

type FortniteTokenBody struct {
	GrantType string `form:"grant_type" binding:"required"`
  ExchangeCode string `form:"exchange_code"`
	Username string `form:"username"`
	Password string `form:"password"`
	TokenType string `form:"token_type"`
}

func PostFortniteToken(c *fiber.Ctx) error {
	var body FortniteTokenBody

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(aid.ErrorBadRequest("Invalid Request Body"))	
	}

	if action, ok := oauthTokenGrantTypes[body.GrantType]; ok {
		return action(c, &body)
	}

	return c.Status(400).JSON(aid.ErrorBadRequest("Invalid Grant Type"))
}

func PostTokenClientCredentials(c *fiber.Ctx, body *FortniteTokenBody) error {
	if aid.Config.Fortnite.DisableClientCredentials {
		return c.Status(400).JSON(aid.ErrorBadRequest("Client Credentials is disabled."))
	}

	clientCredentials, err := aid.JWTSign(aid.JSON{
    "creation_date": time.Now().Format("2006-01-02T15:04:05.999Z"),
		"clsvc": "prod-fn",
		"t": "s",
		"mver": false,
		"clid": aid.Hash([]byte(c.IP())),
		"ic": true,
		"exp": 1707772234,
		"iat": 1707757834,
		"jti": "snow-revoke",
		"pfpid": "prod-fn",
		"am": "client_credentials",
  })
  if err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(aid.ErrorInternalServer)
  }

	return c.Status(200).JSON(aid.JSON{
		"access_token": clientCredentials,
		"application_id": "fghi4567FNFBKFz3E4TROb0bmPS8h1GW",
		"token_type": "bearer",
		"client_id": aid.Hash([]byte(c.IP())),
		"client_service": "prod-fn",
		"internal_client": true,
		"product_id": "prod-fn",
		"expires_in": 3600,
		"expires_at": time.Now().Add(time.Hour).Format("2006-01-02T15:04:05.999Z"),
	})
}

func PostTokenExchangeCode(c *fiber.Ctx, body *FortniteTokenBody) error {
  if body.ExchangeCode == "" {
    return c.Status(400).JSON(aid.ErrorBadRequest("Exchange Code is empty"))
  }

  codeParts := strings.Split(body.ExchangeCode, ".")
  if len(codeParts) != 2 {
    return c.Status(400).JSON(aid.ErrorBadRequest("Invalid Exchange Code"))
  }

  code, failed := aid.KeyPair.DecryptAndVerifyB64(codeParts[0], codeParts[1])
  if failed {
    return c.Status(400).JSON(aid.ErrorBadRequest("Invalid Exchange Code"))
  }

  personParts := strings.Split(string(code), "=")
  if len(personParts) != 2 {
    return c.Status(400).JSON(aid.ErrorBadRequest("Invalid Exchange Code"))
  }

  personId := personParts[0]
  expire, err := time.Parse("2006-01-02T15:04:05.999Z", personParts[1])
  if err != nil {
    return c.Status(400).JSON(aid.ErrorBadRequest("Invalid Exchange Code"))
  }

  if expire.Add(time.Hour).Before(time.Now()) {
    return c.Status(400).JSON(aid.ErrorBadRequest("Invalid Exchange Code"))
  }

  person := p.Find(personId)
  if person == nil {
    return c.Status(400).JSON(aid.ErrorBadRequest("Invalid Exchange Code"))
  }

  access, err := aid.JWTSign(aid.JSON{
    "snow_id": person.ID, // custom
    "creation_date": time.Now().Format("2006-01-02T15:04:05.999Z"),
		"am": "exchange_code",
  })
  if err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(aid.ErrorInternalServer)
  }

  refresh, err := aid.JWTSign(aid.JSON{
    "snow_id": person.ID,
    "creation_date": time.Now().Format("2006-01-02T15:04:05.999Z"),
		"am": "exchange_code",
  })
  if err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(aid.ErrorInternalServer)
  }

  return c.Status(200).JSON(aid.JSON{
    "access_token": "eg1~" + access,
    "account_id": person.ID,
    "client_id": c.IP(),
    "client_service": "fortnite",
    "app": "fortnite",
    "device_id": "default",
    "display_name": person.DisplayName,
    "expires_at": time.Now().Add(time.Hour * 24).Format("2006-01-02T15:04:05.999Z"),
    "expires_in": 86200,
    "internal_client": true,
    "refresh_expires": 86200,
    "refresh_expires_at": time.Now().Add(time.Hour * 24).Format("2006-01-02T15:04:05.999Z"),
    "refresh_token": "eg1~" + refresh,
    "token_type": "bearer",
    "product_id": "prod-fn",
    "sandbox_id": "fn",
  })
}

func PostTokenPassword(c *fiber.Ctx, body *FortniteTokenBody) error {
	if aid.Config.Fortnite.Password {
		return c.Status(400).JSON(aid.ErrorBadRequest("Username and password authentication is disabled for security reasons. Please use an exchange code given by the discord bot."))
	}

	if body.Username == "" || body.Password == "" {
		return c.Status(400).JSON(aid.ErrorBadRequest("Username/Password is empty"))
	}

	person := p.FindByDisplay(strings.Split(body.Username, "@")[0])
	if person == nil {
		return c.Status(400).JSON(aid.ErrorBadRequest("No Account Found"))
	}

	access, err := aid.JWTSign(aid.JSON{
		"snow_id": person.ID, // custom
		"creation_date": time.Now().Format("2006-01-02T15:04:05.999Z"),
		"am": "password",
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(aid.ErrorInternalServer)
	}

	refresh, err := aid.JWTSign(aid.JSON{
		"snow_id": person.ID,
		"creation_date": time.Now().Format("2006-01-02T15:04:05.999Z"),
		"am": "password",
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(aid.ErrorInternalServer)
	}

	return c.Status(200).JSON(aid.JSON{
		"access_token": "eg1~" + access,
		"account_id": person.ID,
		"client_id": c.IP(),
		"client_service": "fortnite",
		"app": "fortnite",
		"device_id": "default",
		"display_name": person.DisplayName,
		"expires_at": time.Now().Add(time.Hour * 24).Format("2006-01-02T15:04:05.999Z"),
		"expires_in": 86200,
		"internal_client": true,
		"refresh_expires": 86200,
		"refresh_expires_at": time.Now().Add(time.Hour * 24).Format("2006-01-02T15:04:05.999Z"),
		"refresh_token": "eg1~" + refresh,
		"token_type": "bearer",
		"product_id": "prod-fn",
		"sandbox_id": "fn",
	})
}

func GetTokenVerify(c *fiber.Ctx) error {
	snowId, err := aid.GetSnowFromToken(c.Get("Authorization"))
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(aid.ErrorBadRequest("Invalid Access Token"))
	}

	person := p.Find(snowId)
	if person == nil {
		return c.Status(fiber.StatusForbidden).JSON(aid.ErrorBadRequest("Invalid Access Token"))
	}

	return c.Status(200).JSON(aid.JSON{
		"app": "fortnite",
		"token": strings.ReplaceAll(c.Get("Authorization"), "bearer eg1~", ""),
		"token_type": "bearer",
		"expires_at": time.Now().Add(time.Hour * 24).Format("2006-01-02T15:04:05.999Z"),
		"expires_in": 86200,
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

func DeleteToken(c *fiber.Ctx) error {
	return c.Status(200).JSON(aid.JSON{})
}

func MiddlewareFortnite(c *fiber.Ctx) error {
	snowId, err := aid.GetSnowFromToken(c.Get("Authorization"))
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(aid.ErrorBadRequest("Invalid Access Token"))
	}

	person := p.Find(snowId)
	if person == nil {
		return c.Status(fiber.StatusForbidden).JSON(aid.ErrorBadRequest("Invalid Access Token"))
	}

	if person.GetLatestActiveBan() != nil {
		return c.Status(fiber.StatusForbidden).JSON(aid.ErrorBadRequest("Account is banned"))
	}

	c.Locals("person", person)
	return c.Next()
}

func MiddlewareWeb(c *fiber.Ctx) error {
	snowId, err := aid.GetSnowFromToken(c.Get("Authorization"))
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(aid.JSON{"error":"Invalid Access Token"})
	}

	person := p.Find(snowId)
	if person == nil {
		return c.Status(fiber.StatusForbidden).JSON(aid.JSON{"error":"Invalid Access Token"})
	}

	c.Locals("person", person)
	return c.Next()
}

func GetPublicAccount(c *fiber.Ctx) error {
	person := p.Find(c.Params("accountId"))
	if person == nil {
		return c.Status(400).JSON(aid.ErrorBadRequest("No Account Found"))
	}

	return c.Status(200).JSON(aid.JSON{
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

	return c.Status(200).JSON(response)
}

func GetPublicAccountExternalAuths(c *fiber.Ctx) error {
	person := p.Find(c.Params("accountId"))
	if person == nil {
		return c.Status(400).JSON(aid.ErrorBadRequest("No Account Found"))
	}

	return c.Status(200).JSON([]aid.JSON{})
}

func GetPublicAccountByDisplayName(c *fiber.Ctx) error {
	person := p.FindByDisplay(c.Params("displayName"))
	if person == nil {
		return c.Status(400).JSON(aid.ErrorBadRequest("No Account Found"))
	}

	return c.Status(200).JSON(aid.JSON{
		"id": person.ID,
		"displayName": person.DisplayName,
		"externalAuths": []aid.JSON{},
	})
}

func GetPrivacySettings(c *fiber.Ctx) error {
	return c.Status(200).JSON(aid.JSON{
		"privacySettings": aid.JSON{
			"playRegion": "PUBLIC",
			"badges": "PUBLIC",
			"languages": "PUBLIC",
		},
	})
}