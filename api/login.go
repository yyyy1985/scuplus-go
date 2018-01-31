package api

import (
	"log"
	"strconv"

	"github.com/mohuishou/scu"

	"github.com/kataras/iris"
	scujwc "github.com/mohuishou/scujwc-go"
	"github.com/mohuishou/scuplus-go/middleware"
	"github.com/mohuishou/scuplus-go/model"
)

// Login 用户登录
func Login(ctx iris.Context) {
	code := ctx.FormValue("code")

	if code == "" {
		Error(ctx, 10400, "code不能为空", nil)
		return
	}

	// 获取openid
	user := &model.User{}
	if err := user.Wechat.GetOpenid(code); err != nil {
		Error(ctx, 10401, "获取用户信息失败", nil)
		return
	}

	// 登录
	token, err := user.Login()
	if err != nil {
		log.Println(err)
		Error(ctx, 10401, "登录失败", nil)
		return
	}

	userLibrary := model.UserLibrary{}
	model.DB().Model(&user).Related(&userLibrary)

	Success(ctx, "登录成功！", map[string]interface{}{
		"token":  token,
		"verify": user.Verify,
	})
}

// Bind 绑定统一认证账号
func Bind(ctx iris.Context) {
	studentID := ctx.FormValue("student_id")
	password := ctx.FormValue("password")

	if studentID == "" || password == "" {
		Error(ctx, 10400, "参数错误", nil)
		return
	}

	// 验证统一账号是否可以登录
	if _, err := scu.NewCollector(studentID, password); err != nil {
		Error(ctx, 10401, err.Error(), nil)
		return
	}

	user := model.User{
		StudentID: studentID,
		Password:  password,
	}

	uid := middleware.GetUserID(ctx)
	user.Verify = 1
	if err := model.DB().Model(&user).Where("id = ?", uid).Updates(&user).Error; err != nil {
		log.Println("用户绑定账号失败: ", err)
		Error(ctx, 30004, "系统错误！", nil)
		return
	}
}

// BindJwc 绑定教务处
func BindJwc(ctx iris.Context) {
	studentID := ctx.FormValue("student_id")
	passsword := ctx.FormValue("password")
	user := model.User{
		StudentID: studentID,
		Password:  passsword,
	}

	if user.StudentID == "" || user.Password == "" {
		Error(ctx, 30002, "参数错误", nil)
		return
	}

	sid, err := strconv.Atoi(user.StudentID)
	if err != nil {
		Error(ctx, 30400, "学号必须是数字", nil)
		return
	}

	// 验证教务处是否可以登录
	if _, err := scujwc.NewJwc(sid, user.Password); err != nil {
		Error(ctx, 30003, err.Error(), nil)
		return
	}

	uid := middleware.GetUserID(ctx)
	user.JwcVerify = 1
	if err := model.DB().Model(&user).Where("id = ?", uid).Updates(&user).Error; err != nil {
		log.Println("[Error]: ", err)
		Error(ctx, 30004, "系统错误！", nil)
		return
	}
	Success(ctx, "绑定成功！", nil)
}
