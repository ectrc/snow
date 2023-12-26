package handlers

import (
	"math/rand"
	"strconv"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/fortnite"
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
	results := []aid.JSON{}
	for playlist := range fortnite.PlaylistImages {
		results = append(results, createPlaylist(playlist, aid.Config.API.Host + aid.Config.API.Port + "/snow/image/" + playlist + ".png?cache="+strconv.Itoa(rand.Intn(9999))))
	}
	results = append(results, createPlaylist("Playlist_DefaultSolo", "http://bucket.retrac.site/55737fa15677cd57fab9e7f4499d62f89cfde320.png"))

	return c.Status(200).JSON(aid.JSON{
		"Panels": []aid.JSON{
			{
				"PanelName": "1",
				"Pages": []aid.JSON{{
					"results": results,
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

func GetContentPages(c *fiber.Ctx) error {
	seasonString := strconv.Itoa(aid.Config.Fortnite.Season)

	playlists := []aid.JSON{}
	for playlist := range fortnite.PlaylistImages {
		playlists = append(playlists, aid.JSON{
			"image": aid.Config.API.Host + aid.Config.API.Port + "/snow/image/" + playlist + ".png?cache="+strconv.Itoa(rand.Intn(9999)),
			"playlist_name": playlist,
			"hidden": false,
		})
	}

	backgrounds := []aid.JSON{}
	switch aid.Config.Fortnite.Season {
	case 11:
		backgrounds = append(backgrounds, aid.JSON{
			"key": "lobby",
			"stage": "Winter19",
		})
	default:
		backgrounds = append(backgrounds, aid.JSON{
			"key": "lobby",
			"stage": "season" + seasonString,
		})
	}

	return c.Status(fiber.StatusOK).JSON(aid.JSON{
		"subgameselectdata": aid.JSON{
			"saveTheWorldUnowned": aid.JSON{
				"message": aid.JSON{
					"title": "Co-op PvE",
					"body": "Cooperative PvE storm-fighting adventure!",
					"spotlight": false,
					"hidden": true,
					"messagetype": "normal",
				},
			},
			"battleRoyale": aid.JSON{
				"message": aid.JSON{
					"title": "100 Player PvP",
					"body": "100 Player PvP Battle Royale.\n\nPvE progress does not affect Battle Royale.",
					"spotlight": false,
					"hidden": true,
					"messagetype": "normal",
				},
			},
			"creative": aid.JSON{
				"message": aid.JSON{
					"title": "New Featured Islands!",
					"body": "Your Island. Your Friends. Your Rules.\n\nDiscover new ways to play Fortnite, play community made games with friends and build your dream island.",
					"spotlight": false,
					"hidden": true,
					"messagetype": "normal",
				},
			},
			"lastModified": "0000-00-00T00:00:00.000Z",
		},
		"dynamicbackgrounds": aid.JSON{
			"backgrounds": aid.JSON{"backgrounds": backgrounds},
			"lastModified": "0000-00-00T00:00:00.000Z",
		},
		"shopSections": aid.JSON{
			"sectionList": aid.JSON{
				"sections": []aid.JSON{
          {
            "bSortOffersByOwnership": false,
            "bShowIneligibleOffersIfGiftable": false,
            "bEnableToastNotification": true,
            "background":  aid.JSON{
              "stage": "default",
              "_type": "DynamicBackground",
              "key": "vault",
            },
            "_type": "ShopSection",
            "landingPriority": 0,
            "bHidden": false,
            "sectionId": "Featured",
            "bShowTimer": true,
            "sectionDisplayName": "Featured",
            "bShowIneligibleOffers": true,
          },
          {
            "bSortOffersByOwnership": false,
            "bShowIneligibleOffersIfGiftable": false,
            "bEnableToastNotification": true,
            "background":  aid.JSON{
              "stage": "default",
              "_type": "DynamicBackground",
              "key": "vault",
            },
            "_type": "ShopSection",
            "landingPriority": 1,
            "bHidden": false,
            "sectionId": "Daily",
            "bShowTimer": true,
            "sectionDisplayName": "Daily",
            "bShowIneligibleOffers": true,
          },
          {
            "bSortOffersByOwnership": false,
            "bShowIneligibleOffersIfGiftable": false,
            "bEnableToastNotification": false,
            "background":  aid.JSON{
              "stage": "default",
              "_type": "DynamicBackground",
              "key": "vault",
            },
            "_type": "ShopSection",
            "landingPriority": 2,
            "bHidden": false,
            "sectionId": "Battlepass",
            "bShowTimer": false,
            "sectionDisplayName": "Battle Pass",
            "bShowIneligibleOffers": false,
          },
        },
			},
			"lastModified": "0000-00-00T00:00:00.000Z",
		},
		"playlistinformation": aid.JSON{
			"conversion_config": aid.JSON{
				"enableReferences": true,
				"containerName": "playlist_info",
				"contentName": "playlists",
			},
			"playlist_info": aid.JSON{
				"playlists": playlists,
			},
			"is_tile_hidden": false,
			"show_ad_violator": false,
			"frontend_matchmaking_header_style": "Basic",
			"frontend_matchmaking_header_text_description": "Watch @ 3PM EST",
			"frontend_matchmaking_header_text": "ECS Qualifiers",
			"lastModified": "0000-00-00T00:00:00.000Z",
		},
		"lastModified": "0000-00-00T00:00:00.000Z",
	})
}