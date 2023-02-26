package utils

import "errors"

var (
	ERROR_UNIFIED           = errors.New("统一认证服务出现问题，请稍后再试")
	ERROR_SERVER            = errors.New("服务器出了问题, 请稍后再试")
	ERROR_ID_PWD            = errors.New("账号或者密码错误,请确认能通过统一认证服务登录教务系统")
	ERROR_CAPTCHA           = errors.New("验证码识别错误, 重试一下？")
	ERROR_UNKOWN            = errors.New("未知错误")
	ERROR_JWC               = errors.New("教务系统出了点问题, 请稍后再试")
	ERROR_DATA              = errors.New("数据错误")
	ERROR_NO_USER           = errors.New("用户不存在")
	ERROR_INPUT             = errors.New("参数有误")
	ERROR_STUDENT_NOT_FOUND = errors.New("查无此人")
	/*
		驼峰命名新变量, 太懒了, 旧项目混着用, 新项目统一用驼峰
	*/
	ErrorUnified          = errors.New("统一认证服务出了点问题，请稍后再试")
	ErrorFailLogin        = errors.New("统一认证登录失败，请确认您能正常登录my.csu.edu.cn")
	ErrorIdPwd            = errors.New("您提供的统一认证账号密码错误")
	ErrorIdPwdWithCaptcha = errors.New("您已触发验证码，提供的统一认证账号密码有误")
	ErrorLocked           = errors.New("密码错误次数过多，您的统一认证账号已被暂时冻结，请5-10分钟后再试")
	ErrorJwc              = errors.New("教务系统出了点问题,请稍后再试")
	ErrorServer           = errors.New("服务器出了点问题,请稍后再试")
	ErrorInput            = errors.New("参数有误")
	ErrorRegister         = errors.New("用户未托管")
	ErrorReRegister       = errors.New("重复托管，账号密码已覆盖")
	ErrorPwdChanged       = errors.New("用户密码已变更")
	ErrorDeleteUser       = errors.New("用户注销失败")
	ErrorCaptcha          = errors.New("验证码自动识别失败, 请重试")
	ErrorOpenid           = errors.New("Openid获取失败")
)
