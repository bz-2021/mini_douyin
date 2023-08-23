package pojo

type Comment struct {
	Id         int64  `gorm:"primaryKey"`
	UserId     int64  `gorm:"column:user_id"`
	VideoId    int64  `gorm:"column:video_id"`
	Content    string `gorm:"column:content"`
	CreateDate string `gorm:"column:create_date;type:timestamp"`
}

func (Comment) TableName() string {
	return "comment"
}
