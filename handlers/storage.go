package handlers

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/storage"
	"github.com/gofiber/fiber/v2"
)

func GetCloudStorageFiles(c *fiber.Ctx) error {
	files := map[string][]byte {
		"DefaultEngine.ini": storage.GetDefaultEngine(),
		"DefaultGame.ini": storage.GetDefaultGame(),
		"DefaultRuntimeOptions.ini": storage.GetDefaultRuntime(),
	}
	result := []aid.JSON{}

	for name, data := range files {
		sumation1 := sha1.Sum(data)
		sumation256 := sha256.Sum256(data)

		result = append(result, aid.JSON{
			"uniqueFilename": name,
			"filename": name,
			"hash": hex.EncodeToString(sumation1[:]),
			"hash256": hex.EncodeToString(sumation256[:]),
			"length": len(data),
			"contentType": "application/octet-stream",
			"uploaded": aid.TimeStartOfDay(),
			"storageType": "S3",
			"storageIds": aid.JSON{
				"primary": "primary",
			},
			"doNotCache": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func GetCloudStorageConfig(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(aid.JSON{
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
	switch c.Params("fileName") {
	case "DefaultEngine.ini":
		c.Set("Content-Type", "application/octet-stream")
		c.Status(fiber.StatusOK)
		c.Send(storage.GetDefaultEngine())
		return nil
	}

	return c.Status(400).JSON(aid.ErrorBadRequest)
}

func GetUserStorageFiles(c *fiber.Ctx) error {
	return c.Status(200).JSON([]aid.JSON{})
}

func GetUserStorageFile(c *fiber.Ctx) error {
	return c.Status(200).JSON(aid.JSON{})
}

func PutUserStorageFile(c *fiber.Ctx) error {
	return c.Status(200).JSON(aid.JSON{})
}
