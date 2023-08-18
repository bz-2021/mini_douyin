package pojo

import "time"

type Video struct {
	Id        int64     `gorm:"primaryKey"`
	UserId    string    `gorm:"column:user_id"`
	PlayUrl   string    `gorm:"play_url"`
	CoverUrl  string    `gorm:"cover_url"`
	Title     string    `gorm:"title"`
	CrateDate time.Time `gorm:"create_date"`
}
