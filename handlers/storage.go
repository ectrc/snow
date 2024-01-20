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
	sum := sha1.Sum(storage.GetDefaultEngine())
	more := sha256.Sum256(storage.GetDefaultEngine())

	return c.Status(fiber.StatusOK).JSON([]fiber.Map{
		{
			"uniqueFilename": "DefaultEngine.ini",
			"filename": "DefaultEngine.ini",
			"hash": hex.EncodeToString(sum[:]),
			"hash256": hex.EncodeToString(more[:]),
			"length": len(storage.GetDefaultEngine()),
			"contentType": "application/octet-stream",
			"uploaded": aid.TimeStartOfDay(),
			"storageType": "S3",
			"storageIds": fiber.Map{
				"primary": "primary",
			},
			"doNotCache": false,
		},
	})
}

func GetCloudStorageConfig(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
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
