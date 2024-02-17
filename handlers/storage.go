package handlers

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/person"
	"github.com/ectrc/snow/storage"
	"github.com/gofiber/fiber/v2"
)

type cloudstorage struct {
	f []aid.JSON
}

func (c *cloudstorage) Add(name string, bytes []byte) error {
	sumation1 := sha1.Sum(bytes)
	sumation256 := sha256.Sum256(bytes)

	c.f = append(c.f, aid.JSON{
		"uniqueFilename": name,
		"filename": name,
		"hash": hex.EncodeToString(sumation1[:]),
		"hash256": hex.EncodeToString(sumation256[:]),
		"length": len(bytes),
		"contentType": "application/octet-stream",
		"uploaded": time.Now().Format(time.RFC3339),
		"storageType": "S3",
		"storageIds": aid.JSON{
			"primary": "primary",
		},
		"doNotCache": true,
	})

	return nil
}

func (c *cloudstorage) Get() []aid.JSON {
	if c.f == nil {
		c.f = []aid.JSON{}
	}
	return c.f
}

func GetCloudStorageFiles(c *fiber.Ctx) error {
	lookup := map[string][]byte {
		"DefaultEngine.ini": storage.GetDefaultEngine(),
		"DefaultGame.ini": storage.GetDefaultGame(),
		"DefaultRuntimeOptions.ini": storage.GetDefaultRuntime(),
	}
	
	files := cloudstorage{}
	for name, bytes := range lookup {
		files.Add(name, bytes)
	}

	return c.Status(200).JSON(files.Get())
}

func GetCloudStorageConfig(c *fiber.Ctx) error {
	return c.Status(200).JSON(aid.JSON{
		"enumerateFilesPath": "/api/cloudstorage/system",
		"enableMigration": true,
		"enableWrites": true,
		"epicAppName": "Live",
		"isAuthenticated": true,
		"disableV2": true,
		"lastUpdated": aid.TimeStartOfDay(),
		"transports": []string{},
	})
}

func GetCloudStorageFile(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/octet-stream")
	switch c.Params("fileName") {
	case "DefaultEngine.ini":
		return c.Status(200).Send(storage.GetDefaultEngine())
	case "DefaultGame.ini":
		return c.Status(200).Send(storage.GetDefaultGame())
	case "DefaultRuntimeOptions.ini":
		return c.Status(200).Send(storage.GetDefaultRuntime())
	}

	return c.Status(404).JSON(aid.ErrorBadRequest("File not found"))
}

func GetUserStorageFiles(c *fiber.Ctx) error {
	if !aid.Config.Amazon.Enabled {
		return c.Status(200).JSON([]aid.JSON{})
	}
	person := c.Locals("person").(*person.Person)
	files := cloudstorage{}

	file, err := storage.Repo.Amazon.GetUserFile(person.ID)
	if err == nil {
		files.Add("ClientSettings.sav", file)
	}
	
	return c.Status(200).JSON(files.Get())
}

func GetUserStorageFile(c *fiber.Ctx) error {
	if !aid.Config.Amazon.Enabled {
		return c.SendStatus(204)
	}
	person := c.Locals("person").(*person.Person)

	file, err := storage.Repo.Amazon.GetUserFile(person.ID)
	if err != nil {
		return c.Status(500).JSON(aid.ErrorBadRequest("Failed to retrieve user file"))
	}

	c.Set("Content-Type", "application/octet-stream")
	return c.Status(200).Send(file)
}

func PutUserStorageFile(c *fiber.Ctx) error {
	if !aid.Config.Amazon.Enabled {
		return c.SendStatus(204)
	}
	person := c.Locals("person").(*person.Person)
	body := c.Body()

	err := storage.Repo.Amazon.CreateUserFile(person.ID, body)
	if err != nil {
		return c.Status(500).JSON(aid.ErrorBadRequest("Failed to create user file"))
	}

	return c.SendStatus(204)
}
