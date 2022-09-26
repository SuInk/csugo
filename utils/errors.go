package utils

import "errors"

var (
	ERROR_SERVER     = errors.New("服务器出了点问题，重试一下？")
	ERROR_ID_PWD     = errors.New("账号或者密码错误,请在my.csu.edu.cn登录后重试")
	ERROR_CAPTCHA    = errors.New("验证码识别错误，,重试一下？")
	ERROR_UNKOWN     = errors.New("未知错误")
	ERROR_JWC        = errors.New("教务系统出了点问题,请重试")
	ERROR_DATA       = errors.New("数据错误")
	ERROR_NO_USER    = errors.New("用户不存在")
	ERROR_INPUT      = errors.New("参数有误")
	ERROR_NO_STUDENT = errors.New("查无此人")
)