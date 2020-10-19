package main

import (
	"time"

	"gorm.io/gorm"
)

type pageEvent struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	PageID          int64     `json:"page_id"`
	EventTimestamp  int64     `json:"event_timestamp"`
	LiveVideoID     int64     `json:"live_video_id"`
	LiveVideoStatus string    `json:"live_video_status"`
}

func (e *pageEvent) getEvent(db *gorm.DB) error {
	return db.First(&e, e.ID).Error
}

func (e *pageEvent) deleteEvent(db *gorm.DB) error {
	return db.Delete(&e, e.ID).Error
}

func (e *pageEvent) createEvent(db *gorm.DB) error {
	return db.Create(&e).Error
}

func (e *pageEvent) createFromFbEvent(db *gorm.DB, fbE *fbPageEvent) error {
	e.PageID = fbE.Entry[0].ID
	e.EventTimestamp = fbE.Entry[0].Time
	e.LiveVideoID = fbE.Entry[0].Changes[0].Value.ID
	e.LiveVideoStatus = fbE.Entry[0].Changes[0].Value.Status
	return e.createEvent(db)
}

func (e *pageEvent) getEvents(db *gorm.DB, start, count int) ([]pageEvent, error) {
	events := []pageEvent{}
	err := db.Limit(count).Offset(start).Find(&events).Error
	return events, err
}
