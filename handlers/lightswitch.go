package handlers

import (
	"strconv"
	"time"

	"github.com/ectrc/snow/aid"
	"github.com/gofiber/fiber/v2"
)

func GetLightswitchBulkStatus(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON([]aid.JSON{{
		"serviceInstanceId": "fortnite",
		"status": "UP",
		"banned": false,
	}})
}

func GetTimelineCalendar(c *fiber.Ctx) error {
	events := []aid.JSON{
		{
			"activeUntil": aid.TimeEndOfWeekString(),
			"activeSince": "0001-01-01T00:00:00Z",
			"activeEventId": "EventFlag.Season" + strconv.Itoa(aid.Config.Fortnite.Season),
		},
		{
			"activeUntil": aid.TimeEndOfWeekString(),
			"activeSince": "0001-01-01T00:00:00Z",
			"activeEventId": "EventFlag.LobbySeason" + strconv.Itoa(aid.Config.Fortnite.Season),
		},
	}

	state := aid.JSON{
		"seasonNumber": aid.Config.Fortnite.Season,
		"seasonTemplateId": "AthenaSeason:AthenaSeason" + strconv.Itoa(aid.Config.Fortnite.Season),
		"seasonBegin": time.Now().Add(-time.Hour * 24 * 7).Format("2006-01-02T15:04:05.000Z"),
		"seasonEnd": time.Now().Add(time.Hour * 24 * 7).Format("2006-01-02T15:04:05.000Z"),
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
		},
		"currentTime": time.Now().Format("2006-01-02T15:04:05.000Z"),
		"cacheIntervalMins": 5,
		"eventsTimeOffsetHrs": 0,
	})
}