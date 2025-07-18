package utils

import "errors"

var (
	ErrEmailExisted  = errors.New("email đã tồn tại")
	ErrUserNotFound  = errors.New("không tìm thấy người dùng")
	ErrWrongPassword = errors.New("sai mật khẩu")
)