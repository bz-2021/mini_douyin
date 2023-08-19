package pojo

type Favorite struct {
	Id      int64 `gorm:"primaryKey"`
	VideoId int64 `gorm:"column:video_id"`
	UserId  int64 `gorm:"column:user_id"`
}

func (Favorite) TableName() string {
	return "favorite"
}
