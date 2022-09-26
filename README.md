# csugo
中南大学校内查询API，提供教务，校车，招聘查询，基于Beego框架

本项目基于https://github.com/csuhan/csugo 改写

## 如何使用

```bash
# 安装Git
https://git-scm.com/
# 安装go运行环境
https://go.dev/doc/install
# 下载项目到本地
git clone git@github.com:SuInk/csugo.git
# 运行项目
cd csugo
go build main.go
go run main.go
```

## 本地访问

打开浏览器输入`http://localhost:9090`

运行成功

## 路由格式

`http://localhost:9090/api/v1/jwc/{{学号}}/{{统一认证密码的base64编码}}/class/2022-2023-1/1?token=csugo-token`

...

## 目前可用

* 课表查询
* 成绩查询
* 排名查询
* 教学周历
* 成绩订阅 https://github.com/SuInk/csu-grade-push
* 导入日历 https://github.com/SuInk/csu-import

## 暂不可用

* 馆藏查询 图书馆只能校园网访问
* 校车查询 学校改了接口
* 校友查询 每课停用了，虽然有绮课代替，但我懒得改
* 四六级查询 小程序审核不通过

## 想写的功能

* GPA 计算

* 查别人课表，类似每课
* 课表可以存储到本地，不用每次请求
* 更改路由格式

要毕业了，不想写了，现在小程序累计访问3000+，平均每天400+



