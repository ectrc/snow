package handlers

import (
	"github.com/ectrc/snow/aid"
	"github.com/gofiber/fiber/v2"
)

func createContentPanel(title string, id string) aid.JSON {
	return aid.JSON{
		"NumPages": 1,
		"AnalyticsId": id,
		"PanelType": "AnalyticsList",
		"AnalyticsListName": title,
		"CuratedListOfLinkCodes": []aid.JSON{},
		"ModelName": "",
		"PageSize": 7,
		"PlatformBlacklist": []aid.JSON{},
		"PanelName": id,
		"MetricInterval": "",
		"SkippedEntriesCount": 0,
		"SkippedEntriesPercent": 0,
		"SplicedEntries": []aid.JSON{},
		"PlatformWhitelist": []aid.JSON{},
		"EntrySkippingMethod": "None",
		"PanelDisplayName": aid.JSON{
			"Category": "Game",
			"NativeCulture": "",
			"Namespace": "CreativeDiscoverySurface_Frontend",
			"LocalizedStrings": []aid.JSON{{
				"key": "en",
				"value": title,
			}},
			"bIsMinimalPatch": false,
			"NativeString": title,
			"Key": "",
		},
		"PlayHistoryType": "RecentlyPlayed",
		"bLowestToHighest": false,
		"PanelLinkCodeBlacklist": []aid.JSON{},
		"PanelLinkCodeWhitelist": []aid.JSON{},
		"FeatureTags": []aid.JSON{},
		"MetricName": "",
	}
}

func createPlaylist(mnemonic string, image string) aid.JSON {
	return aid.JSON{
		"linkData": aid.JSON{
			"namespace": "fn",
			"mnemonic": mnemonic,
			"linkType": "BR:Playlist",
			"active": true,
			"disabled": false,
			"version": 1,
			"moderationStatus": "Unmoderated",
			"accountId": "epic",
			"creatorName": "Epic",
			"descriptionTags": []string{},
			"metadata": aid.JSON{
				"image_url": image,
				"matchmaking": aid.JSON{
					"override_playlist": mnemonic,
				},
			},
		},
		"lastVisited": nil,
		"linkCode": mnemonic,
		"isFavorite": false,
	}
}

func PostDiscovery(c *fiber.Ctx) error {
	return c.Status(200).JSON(aid.JSON{
		"Panels": []aid.JSON{
			{
				"PanelName": "1",
				"Pages": []aid.JSON{{
					"results": []aid.JSON{
						createPlaylist("playlist_defaultsolo", "https://cdn2.unrealengine.com/solo-1920x1080-1920x1080-bc0a5455ce20.jpg"),
						createPlaylist("playlist_defaultduo", "https://cdn2.unrealengine.com/duos-1920x1080-1920x1080-5a411fe07b21.jpg"),
						createPlaylist("playlist_trios", "https://cdn2.unrealengine.com/trios-1920x1080-1920x1080-d5054bb9691a.jpg"),
						createPlaylist("playlist_defaultsquad", "https://cdn2.unrealengine.com/squads-1920x1080-1920x1080-095c0732502e.jpg"),
					},
					"hasMore": false,
				}},
			},
		},
		"TestCohorts": []string{
			"playlists",
		},
		"ModeSets": aid.JSON{},
	})
}

func PostAssets(c *fiber.Ctx) error {
	var body struct {
		DAD_CosmeticItemUserOptions int `json:"DAD_CosmeticItemUserOptions"`
		FortCreativeDiscoverySurface int `json:"FortCreativeDiscoverySurface"`
		FortPlaylistAthena int `json:"FortPlaylistAthena"`
	}

	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(400).JSON(aid.JSON{"error":err.Error()})
	}

	testCohort := aid.JSON{
		"AnalyticsId": "0",
		"CohortSelector": "PlayerDeterministic",
		"PlatformBlacklist": []aid.JSON{},
		"ContentPanels": []aid.JSON{
			createContentPanel("Featured", "1"),
		},
		"PlatformWhitelist": []aid.JSON{},
		"SelectionChance": 0.1,
		"TestName": "playlists",
	}

	return c.Status(200).JSON(aid.JSON{
		"FortCreativeDiscoverySurface": aid.JSON{
			"meta": aid.JSON{
				"promotion": 1,
			},
			"assets": aid.JSON{
				"CreativeDiscoverySurface_Frontend": aid.JSON{
					"meta": aid.JSON{
						"revision": 1,
						"headRevision": 1,
						"promotion": 1,
						"revisedAt": "0000-00-00T00:00:00.000Z",
						"promotedAt": "0000-00-00T00:00:00.000Z",
					},
					"assetData": aid.JSON{
						"AnalyticsId": "t412",
						"TestCohorts": []aid.JSON{
							testCohort,
						},
						"GlobalLinkCodeBlacklist": []aid.JSON{},
						"SurfaceName": "CreativeDiscoverySurface_Frontend",
						"TestName": "20.10_4/11/2022_hero_combat_popularConsole",
						"primaryAssetId": "FortCreativeDiscoverySurface:CreativeDiscoverySurface_Frontend",
						"GlobalLinkCodeWhitelist": []aid.JSON{},
					},
				},
			},
		},
	})
}