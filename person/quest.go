package person

import (
	"github.com/ectrc/snow/storage"
	"github.com/google/uuid"
)

type Quest struct {
	ID              string
	TemplateID      string
	Objectives      []string
	ObjectiveCounts []int64
	BundleID        string
	ScheduleID      string
}

func NewQuest(templateID string, bundleID string, scheduleID string) *Quest {
	return &Quest{
		ID:              uuid.New().String(),
		TemplateID:      templateID,
		Objectives:      []string{},
		ObjectiveCounts: []int64{},
		BundleID:        bundleID,
		ScheduleID:      scheduleID,
	}
}

func FromDatabaseQuest(quest *storage.DB_Quest, profileType *string) *Quest {
	return &Quest{
		ID:              quest.ID,
		TemplateID:      quest.TemplateID,
		Objectives:      quest.Objectives,
		ObjectiveCounts: quest.ObjectiveCounts,
		BundleID:        quest.BundleID,
		ScheduleID:      quest.ScheduleID,
	}
}

func (q *Quest) Delete() {
	storage.Repo.DeleteQuest(q.ID)
}

func (q *Quest) AddObjective(objective string, count int64) {
	q.Objectives = append(q.Objectives, objective)
	q.ObjectiveCounts = append(q.ObjectiveCounts, count)
}

func (q *Quest) SetObjectiveCount(objective string, count int64) {
	for i, obj := range q.Objectives {
		if obj == objective {
			q.ObjectiveCounts[i] = count
			return
		}
	}
}

func (q *Quest) UpdateObjectiveCount(objective string, delta int64) {
	for i, obj := range q.Objectives {
		if obj == objective {
			q.ObjectiveCounts[i] += delta
			return
		}
	}
}

func (q *Quest) GetObjectiveCount(objective string) int64 {
	for i, obj := range q.Objectives {
		if obj == objective {
			return q.ObjectiveCounts[i]
		}
	}

	return 0
}

func (q *Quest) GetObjectiveIndex(objective string) int {
	for i, obj := range q.Objectives {
		if obj == objective {
			return i
		}
	}

	return -1
}

func (q *Quest) RemoveObjective(objective string) {
	for i, obj := range q.Objectives {
		if obj == objective {
			q.Objectives = append(q.Objectives[:i], q.Objectives[i+1:]...)
			q.ObjectiveCounts = append(q.ObjectiveCounts[:i], q.ObjectiveCounts[i+1:]...)
			return
		}
	}
}

func (q *Quest) ToDatabase(profileId string) *storage.DB_Quest {
	return &storage.DB_Quest{
		ProfileID:       profileId,
		ID:              q.ID,
		TemplateID:      q.TemplateID,
		Objectives:      q.Objectives,
		ObjectiveCounts: q.ObjectiveCounts,
		BundleID:        q.BundleID,
		ScheduleID:      q.ScheduleID,
	}
}

func (q *Quest) Save() {
	//storage.Repo.SaveQuest(q.ToDatabase())
}