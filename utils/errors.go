package utils

import "errors"

var (
	ERROR_UNIFIED    = errors.New("统一认证服务出现问题，请稍后再试")
	ERROR_SERVER     = errors.New("服务器出了问题，请稍后再试")
	ERROR_ID_PWD     = errors.New("账号或者密码错误,请确认能通过统一认证服务登录教务系统")
	ERROR_CAPTCHA    = errors.New("验证码识别错误，,重试一下？")
	ERROR_UNKOWN     = errors.New("未知错误")
	ERROR_JWC        = errors.New("教务系统出了点问题,请稍后再试")
	ERROR_DATA       = errors.New("数据错误")
	ERROR_NO_USER    = errors.New("用户不存在")
	ERROR_INPUT      = errors.New("参数有误")
	ERROR_NO_STUDENT = errors.New("查无此人")
)
