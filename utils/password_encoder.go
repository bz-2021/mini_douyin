package utils

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// HashAndSalt 对密码进行哈希和加盐
func HashAndSalt(password string) (string, error) {
	// 生成密码的哈希值
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("HasAndSalt 失败%s", err)
	}
	return string(hashedPassword), nil
}

// ComparePasswords 对比明文密码和哈希值是否匹配
func ComparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
