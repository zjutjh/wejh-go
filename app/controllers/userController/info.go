package userController

import (
	"github.com/gin-gonic/gin"
	"wejh-go/app/apiExpection"
	"wejh-go/app/services/sessionServices"
	"wejh-go/app/utils"
)

func GetUserInfo(c *gin.Context) {
	user, err := sessionServices.GetUserSession(c)
	if err != nil {
		_ = c.AbortWithError(200, apiExpection.NotLogin)
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
