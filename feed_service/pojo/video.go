package pojo

type Video struct {
	Id            int64  `gorm:"primaryKey"`
	UserId        string `gorm:"column:user_id"`
	PlayUrl       string `gorm:"column:play_url"`
	CoverUrl      string `gorm:"column:cover_url"`
	Title         string `gorm:"column:title"`
	CrateDate     int64  `gorm:"column:create_date"`
	FavoriteCount int64  `gorm:"column:favorite_count;default:0"`
	CommentCount  int64  `gorm:"column:comment_count;default:0"`
	Author        Author `gorm:"foreignKey:user_id"`
	IsFavorite    bool   `gorm:"-"`
}
