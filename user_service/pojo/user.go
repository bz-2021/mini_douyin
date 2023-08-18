package pojo

type User struct {
	Id              int64  `gorm:"primaryKey"`
	Name            string `gorm:"column:name"`
	Avatar          string `gorm:"column:avatar"`
	Signature       string `gorm:"column:signature"`
	Password        string `gorm:"column:password"`
	FollowCount     int    `gorm:"column:follow_count;default:0"`
	FollowerCount   int    `gorm:"column:follower_count;default:0"`
	BackgroundImage string `gorm:"column:background_image;default:'https://upload-images.jianshu.io/upload_images/5809200-a99419bb94924e6d.jpg?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240'"`
	TotalFavorite   int    `gorm:"column:total_favorite;default:0"`
	WorkCount       int    `gorm:"column:work_count;default:0"`
	FavoriteCount   int    `gorm:"column:favorite_count;default:0"`
}
