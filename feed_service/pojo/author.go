package pojo

type Author struct {
	ID              int64  `gorm:"primary_key"`
	Name            string `gorm:"name"`
	FollowCount     int64  `gorm:"follow_count"`
	FollowerCount   int64  `gorm:"follower_count"`
	IsFollow        bool   `gorm:"is_follow"`
	Avatar          string `gorm:"avatar"`
	BackgroundImage string `gorm:"background_image"`
	Signature       string `gorm:"signature"`
	TotalFavorite   int64  `gorm:"total_favorite"`
	WorkCount       int64  `gorm:"work_count"`
	FavoriteCount   int64  `gorm:"favorite_count"`
}

func (Author) TableName() string {
	return "user"
}
