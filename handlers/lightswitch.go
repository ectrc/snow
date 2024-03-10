package handlers

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/person"
	"github.com/gofiber/fiber/v2"
)

func GetLightswitchBulkStatus(c *fiber.Ctx) error {
	person := c.Locals("person").(*person.Person)

	ban := person.GetLatestActiveBan()

	return c.Status(fiber.StatusOK).JSON([]aid.JSON{{
		"serviceInstanceId": "fortnite",
		"status" :"UP",
		"message": "fortnite is up.",
		"maintenanceUri": nil,
		"allowedActions": []string{"PLAY","DOWNLOAD"},
		"banned": ban != nil && time.Now().Before(ban.Expiry),
		"launcherInfoDTO": aid.JSON{
			"appName":"Fortnite",
			"namespace":"fn",
		},
	}})
}

func GetFortniteTimeline(c *fiber.Ctx) error {
	userAgent := c.Get("User-Agent")
	if !strings.Contains(userAgent, "++Fortnite") {
		return c.Status(fiber.StatusBadRequest).JSON(aid.ErrorBadRequest("No User Agent"))
	}

	build := regexp.MustCompile(`\d+\.\d+`).FindString(userAgent)
	if len(strings.Split(build, ".")) != 2 {
		return c.Status(fiber.StatusBadRequest).JSON(aid.ErrorBadRequest("Invalid Build"))
	}
	
	season, err := strconv.Atoi(strings.Split(build, ".")[0])
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(aid.ErrorBadRequest("Failed to parse Build"))
	}

	events := []aid.JSON{
		{
			"activeUntil": "9999-12-31T23:59:59.999Z",
			"activeSince": "0001-01-01T00:00:00Z",
			"eventType": "EventFlag.Season" + strings.Split(build, ".")[0],
		},
	}

	switch season {
	case 2:
		events = append(events, aid.JSON{
			"activeUntil": "9999-01-01T00:00:00.000Z",
			"activeSince": "9999-01-01T00:00:00.000Z",
			"eventType": "EventFlag.LobbyWinterDecor",
		})
	case 6:
		events = append(events, aid.JSON{
			"activeUntil": "9999-01-01T00:00:00.000Z",
			"activeSince": "9999-01-01T00:00:00.000Z",
			"eventType": "EventFlag.LobbySeason6Halloween",
		})
	case 11:
		events = append(events, aid.JSON{
			"activeUntil": "9999-01-01T00:00:00.000Z",
			"activeSince": "9999-01-01T00:00:00.000Z",
			"eventType": "EventFlag.LTE_WinterFest2019",
		}, aid.JSON{
			"activeUntil": "9999-01-01T00:00:00.000Z",
			"activeSince": "9999-01-01T00:00:00.000Z",
			"eventType": "EventFlag.LTE_WinterFest",
		}, aid.JSON{
			"activeUntil": "9999-01-01T00:00:00.000Z",
			"activeSince": "9999-01-01T00:00:00.000Z",
			"eventType": "EventFlag.Winterfest.Tree",
		})
	default:
		events = append(events, aid.JSON{
			"activeUntil": "9999-01-01T00:00:00.000Z",
			"activeSince": "9999-01-01T00:00:00.000Z",
			"eventType": "EventFlag.LobbySeason" + strings.Split(build, ".")[0],
		})
	}

	state := aid.JSON{
		"eventNamedWeights": aid.JSON{},
		"seasonNumber": season,
		"seasonTemplateId": "AthenaSeason:AthenaSeason" + strings.Split(build, ".")[0],
		"seasonBegin": time.Now().Add(-time.Hour * 24 * 7).Format("2006-01-02T15:04:05.000Z"),
		"seasonEnd": time.Now().Add(time.Hour * 24 * 65).Format("2006-01-02T15:04:05.000Z"),
		"seasonDisplayedEnd": time.Now().Add(time.Hour * 24 * 7).Format("2006-01-02T15:04:05.000Z"),
		"activeStorefronts": []aid.JSON{},
		"dailyStoreEnd": aid.TimeEndOfDay(),
		"weeklyStoreEnd": aid.TimeEndOfWeekString(),
		"sectionStoreEnds": aid.JSON{},
		"stwEventStoreEnd": aid.TimeEndOfWeekString(),
		"stwWeeklyStoreEnd": aid.TimeEndOfWeekString(),
	}

	client := aid.JSON{
		"states": []aid.JSON{{
			"activeEvents": events,
			"state": state,
			"validFrom": "0001-01-01T00:00:00Z",
		}},
		"cacheExpire": "9999-12-31T23:59:59.999Z",
	}

	return c.Status(fiber.StatusOK).JSON(aid.JSON{
		"channels": aid.JSON{
			"client-events": client,
			"client-matchmaking": aid.JSON{
				"states": []aid.JSON{},
				"cacheExpire": "9999-12-31T23:59:59.999Z",	
			},
		},
		"currentTime": time.Now().Format("2006-01-02T15:04:05.000Z"),
		"cacheIntervalMins": 5,
		"eventsTimeOffsetHrs": 0,
	})
}