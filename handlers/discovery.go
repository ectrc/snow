package handlers

import (
	"strconv"

	"github.com/ectrc/snow/aid"
	"github.com/gofiber/fiber/v2"
)

func GetContentPages(c *fiber.Ctx) error {
	seasonString := strconv.Itoa(aid.Config.Fortnite.Season)

	playlists := []aid.JSON{
		{
			"image": "https://cdn.snows.rocks/squads.png",
			"playlist_name": "Playlist_DefaultSquad",
			"hidden": false,
		},
		{
			"image": "https://cdn.snows.rocks/duos.png",
			"playlist_name": "Playlist_DefaultDuo",
			"hidden": false,
		},
		{
			"image": "https://cdn.snows.rocks/solo.png",
			"playlist_name": "Playlist_DefaultSolo",
			"hidden": false,
		},
		{
			"image": "https://cdn.snows.rocks/arena_solo.png",
			"playlist_name": "Playlist_ShowdownAlt_Solo",
			"hidden": false,
		},
		{
			"image": "https://cdn.snows.rocks/arena_duos.png",
			"playlist_name": "Playlist_ShowdownAlt_Duos",
			"hidden": false,
		},
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
					"image": "https://cdn.snows.rocks/loading_stw.png",
				},
			},
			"saveTheWorld": aid.JSON{
				"message": aid.JSON{
					"title": "Co-op PvE",
					"body": "Cooperative PvE storm-fighting adventure!",
					"spotlight": false,
					"hidden": true,
					"messagetype": "normal",
					"image": "https://cdn.snows.rocks/loading_stw.png",
				},
			},
			"battleRoyale": aid.JSON{
				"message": aid.JSON{
					"title": "100 Player PvP",
					"body": "100 Player PvP Battle Royale.\n\nPvE progress does not affect Battle Royale.",
					"spotlight": false,
					"hidden": true,
					"messagetype": "normal",
					"image": "https://cdn.snows.rocks/loading_br.png",
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
		"tournamentinformation": aid.JSON{
			"tournament_info": aid.JSON{
				"tournaments": []aid.JSON{
					{
						"tournament_display_id": "SnowArenaSolo",
						"playlist_tile_image": "https://cdn.snows.rocks/arena_solo.png",
						"title_line_2" : "ARENA",
					},
					{
						"tournament_display_id": "SnowArenaDuos",
						"playlist_tile_image": "https://cdn.snows.rocks/arena_duos.png",
						"title_line_2" : "ARENA",
					},
				},
			},
			"lastModified": "0000-00-00T00:00:00.000Z",
		},
		"lastModified": "0000-00-00T00:00:00.000Z",
	})
}