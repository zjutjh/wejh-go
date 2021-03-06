package userController

import (
	"github.com/gin-gonic/gin"
	"wejh-go/app/apiExpection"
	"wejh-go/app/services/sessionServices"
	"wejh-go/app/services/userServices"
	"wejh-go/app/utils"
	"wejh-go/config/wechat"
)

type autoLoginForm struct {
	Code      string `json:"code" binding:"required"`
	LoginType string `json:"type"`
}
type passwordLoginForm struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	LoginType string `json:"type"`
}

func AuthByPassword(c *gin.Context) {
	var postForm passwordLoginForm
	err := c.ShouldBindJSON(&postForm)

	if err != nil {
		_ = c.AbortWithError(200, apiExpection.ParamError)
		return
	}
	user, err := userServices.GetUserByUsernameAndPassword(postForm.Username, postForm.Password)
	if err != nil || user == nil {
		_ = c.AbortWithError(200, apiExpection.NoThatPasswordOrWrong)
		return
	}

	err = sessionServices.SetUserSession(c, user)
	if err != nil {
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, gin.H{
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"studentID": user.StudentID,
			"bind": gin.H{
				"zf":   user.ZFPassword != "",
				"lib":  user.LibPassword != "",
				"card": user.CardPassword != "",
			},
			"userType":   user.Type,
			"createTime": user.CreateTime,
		},
	})

}

func WeChatLogin(c *gin.Context) {
	var postForm autoLoginForm
	err := c.ShouldBindJSON(&postForm)

	if err != nil {
		_ = c.AbortWithError(200, apiExpection.ParamError)
		return
	}

	session, err := wechat.MiniProgram.GetAuth().Code2Session(postForm.Code)

	if err != nil {
		_ = c.AbortWithError(200, apiExpection.OpenIDError)
		return
	}

	user := userServices.GetUserByWechatOpenID(session.OpenID)

	if user == nil {
		_ = c.AbortWithError(200, apiExpection.UserNotFind)
		return
	}

	err = sessionServices.SetUserSession(c, user)
	if err != nil {
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, gin.H{
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"studentID": user.StudentID,
			"userType":  user.Type,
			"bind": gin.H{
				"zf":   user.ZFPassword != "",
				"lib":  user.LibPassword != "",
				"card": user.CardPassword != "",
			},
			"createTime": user.CreateTime,
		},
	})
}
