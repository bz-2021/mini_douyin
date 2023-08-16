package pojo

type User struct {
	Id        int64  `gorm:"primaryKey"`
	Name      string `gorm:"column:name"`
	Avatar    string `gorm:"column:avatar"`
	Signature string `gorm:"column:signature"`
	Password  string `gorm:"column:password"`
}
