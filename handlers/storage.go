package handlers

import (
	"crypto/sha1"
	"encoding/hex"

	"github.com/ectrc/snow/aid"
	"github.com/gofiber/fiber/v2"
)

var (
	TEMP_STORAGE = map[string][]byte{
		"DefaultEngine.ini": []byte(`
[OnlineSubsystemMcp.Xmpp Prod]
bUseSSL=false
ServerAddr="ws://127.0.0.1:80"
ServerPort="80"

[OnlineSubsystemMcp.Xmpp]
bUseSSL=false
ServerAddr="ws://127.0.0.1:80"
ServerPort="80"

[/Script/Qos.QosRegionManager]
NumTestsPerRegion=5
PingTimeout=3.0

[XMPP]
bEnableWebsockets=true

[OnlineSubsystemMcp]
bUsePartySystemV2=false

[OnlineSubsystemMcp.OnlinePartySystemMcpAdapter]
bUsePartySystemV2=false

[ConsoleVariables]
n.VerifyPeer=0
FortMatchmakingV2.ContentBeaconFailureCancelsMatchmaking=0
Fort.ShutdownWhenContentBeaconFails=0
FortMatchmakingV2.EnableContentBeacon=0
		`),
		"DefaultGame.ini": []byte(`
[/Script/FortniteGame.FortOnlineAccount]
bEnableEulaCheck=false

[/Script/FortniteGame.FortChatManager]
bShouldRequestGeneralChatRooms=false
bShouldJoinGlobalChat=flase
bShouldJoinFounderChat=false
bIsAthenaGlobalChatEnabled=false

[/Script/FortniteGame.FortGameInstance]
!FrontEndPlaylistData=ClearArray
+FrontEndPlaylistData=(PlaylistName=Playlist_DefaultSolo, PlaylistAccess=(bEnabled=True, bIsDefaultPlaylist=True, bVisibleWhenDisabled=True, bDisplayAsNew=False, CategoryIndex=0, bDisplayAsLimitedTime=False, DisplayPriority=0))
+FrontEndPlaylistData=(PlaylistName=Playlist_DefaultDuo, PlaylistAccess=(bEnabled=True, bIsDefaultPlaylist=True, bVisibleWhenDisabled=True, bDisplayAsNew=False, CategoryIndex=0, bDisplayAsLimitedTime=False, DisplayPriority=1))
+FrontEndPlaylistData=(PlaylistName=Playlist_DefaultSquad, PlaylistAccess=(bEnabled=True, bIsDefaultPlaylist=True, bVisibleWhenDisabled=True, bDisplayAsNew=False, CategoryIndex=0, bDisplayAsLimitedTime=False, DisplayPriority=2))
+FrontEndPlaylistData=(PlaylistName=Playlist_Fill_Squads, PlaylistAccess=(bEnabled=True, bIsDefaultPlaylist=False, bVisibleWhenDisabled=True, bDisplayAsNew=False, CategoryIndex=1, bDisplayAsLimitedTime=False, DisplayPriority=0))
+FrontEndPlaylistData=(PlaylistName=Playlist_Blitz_Solo, PlaylistAccess=(bEnabled=True, bIsDefaultPlaylist=False, bVisibleWhenDisabled=True, bDisplayAsNew=True, CategoryIndex=1, bDisplayAsLimitedTime=True, DisplayPriority=1))
`),
		"DefaultRuntimeOptions.ini": []byte(`
[/Script/FortniteGame.FortRuntimeOptions]
bEnableGlobalChat=false
bDisableGifting=false
bDisableGiftingPC=false
bDisableGiftingPS4=false
bDisableGiftingXB=false`),
	}
)

func GetCloudStorageFiles(c *fiber.Ctx) error {
	engineHash := sha1.Sum(TEMP_STORAGE["DefaultEngine.ini"])
	engineHash256 := sha1.Sum(TEMP_STORAGE["DefaultEngine.ini"])
	gameHash := sha1.Sum(TEMP_STORAGE["DefaultGame.ini"])
	gameHash256 := sha1.Sum(TEMP_STORAGE["DefaultGame.ini"])
	runtimeHash := sha1.Sum(TEMP_STORAGE["DefaultRuntimeOptions.ini"])
	runtimeHash256 := sha1.Sum(TEMP_STORAGE["DefaultRuntimeOptions.ini"])

	return c.Status(fiber.StatusOK).JSON([]aid.JSON{
		{
			"uniqueFilename": "DefaultEngine.ini",
			"filename": "DefaultEngine.ini",
			"hash": hex.EncodeToString(engineHash[:]),
			"hash256": hex.EncodeToString(engineHash256[:]),
			"length": len(TEMP_STORAGE["DefaultEngine.ini"]),
			"contentType": "application/octet-stream",
			"uploaded": "2021-01-01T00:00:00.000Z",
			"storageType": "S3",
			"doNotCache": false,
			"storageIds": []string{"primary"},
		},
		{
			"uniqueFilename": "DefaultGame.ini",
			"filename": "DefaultGame.ini",
			"hash": hex.EncodeToString(gameHash[:]),
			"hash256": hex.EncodeToString(gameHash256[:]),
			"length": len(TEMP_STORAGE["DefaultGame.ini"]),
			"contentType": "application/octet-stream",
			"uploaded": "2021-01-01T00:00:00.000Z",
			"storageType": "S3",
			"doNotCache": false,
			"storageIds": []string{"primary"},
		},
		{
			"uniqueFilename": "DefaultRuntimeOptions.ini",
			"filename": "DefaultRuntimeOptions.ini",
			"hash": hex.EncodeToString(runtimeHash[:]),
			"hash256": hex.EncodeToString(runtimeHash256[:]),
			"length": len(TEMP_STORAGE["DefaultRuntimeOptions.ini"]),
			"contentType": "application/octet-stream",
			"uploaded": "2021-01-01T00:00:00.000Z",
			"storageType": "S3",
			"doNotCache": false,
			"storageIds": []string{"primary"},
		},
	})
}

func GetCloudStorageConfig(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(aid.JSON{
		"enumerateFilesPath": "/api/cloudstorage/system",
		"enableMigration": true,
		"enableWrites": true,
		"epicAppName": "Live",
		"isAuthenticated": true,
		"disableV2": true,
		"lastUpdated": "0000-00-00T00:00:00.000Z",
		"transports": []string{},
	})
}

func GetCloudStorageFile(c *fiber.Ctx) error {
	fileName := c.Params("fileName")
	if fileName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(aid.ErrorBadRequest)
	}

	file, ok := TEMP_STORAGE[fileName]
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(aid.ErrorNotFound)
	}

	return c.Status(fiber.StatusOK).Send(file)
}

func GetUserStorageFiles(c *fiber.Ctx) error {
	basePath := "UserStorage/" + c.Params("accountId") + "/"
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		if err := os.MkdirAll(basePath, 0755); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(aid.ErrorInternalServer)
		}
	} else {
		filePath := basePath
		fileContent, err := os.ReadFile(filePath + "ClientSettings.Sav")
		if err != nil {
			return c.Status(fiber.StatusOK).JSON(aid.JSON{})
		}
		settingsHash := sha1.Sum(fileContent)
		settingsHash256 := sha1.Sum(fileContent)
		return c.Status(fiber.StatusOK).JSON(aid.JSON{
			"uniqueFilename": "ClientSettings.Sav",
			"filename":       "ClientSettings.Sav",
			"hash":           hex.EncodeToString(settingsHash[:]),
			"hash256":        hex.EncodeToString(settingsHash256[:]),
			"length":         len(fileContent),
			"contentType":    "application/octet-stream",
			"uploaded":       "2021-01-01T00:00:00.000Z",
			"storageType":    "S3",
			"doNotCache":     false,
			"storageIds":     []string{"primary"},
		})
	}
	return c.Status(fiber.StatusOK).JSON(aid.JSON{})
}

func GetUserStorageFile(c *fiber.Ctx) error {
	basePath := "UserStorage/" + c.Params("accountId") + "/"
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		if err := os.MkdirAll(basePath, 0755); err != nil {
			return c.Status(fiber.StatusOK).JSON(aid.JSON{})
		}
	} else {
		filePath := basePath + c.Params("fileName")
		_, err := os.Stat(filePath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(aid.ErrorNotFound)
		} else {
			return c.Status(fiber.StatusOK).SendFile(filePath)
		}

	}
	return c.Status(fiber.StatusInternalServerError).JSON(aid.ErrorInternalServer)
}

func PutUserStorageFile(c *fiber.Ctx) error {
	bytes := string(c.BodyRaw())

	if c.Request().Header.ContentLength() > 400000 || strings.ToLower(c.Params("fileName")) != "clientsettings.sav" {
		return c.Status(fiber.StatusBadRequest).JSON(aid.ErrorBadRequest("Invalid File"))
	}

	filePath := "UserStorage/" + c.Params("accountId") + "/"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := os.MkdirAll(filePath, 0755); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(aid.ErrorInternalServer)
		}
	}

	file, err := os.Create(filePath + c.Params("fileName"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(aid.ErrorInternalServer)
	}
	defer file.Close()

	_, err = io.Copy(file, strings.NewReader(bytes))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(aid.ErrorInternalServer)
	}

	return c.Status(fiber.StatusOK).JSON(aid.JSON{})
}
