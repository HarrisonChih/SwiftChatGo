package service

import (
	"fmt"
	"ginchat/models"
	"github.com/gin-gonic/gin"
	"strconv"
	"text/template"
)

// GetIndex
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	index, err := template.ParseFiles("index.html", "views/chat/head.html")
	if err != nil {
		panic(err)
	}
	index.Execute(c.Writer, "index")
	//c.JSON(200, gin.H{
	//	"message": "welcome !",
	//})
}

// ToRegister
// @Tags 用户注册
// @Router /toRegister [get]
func ToRegister(c *gin.Context) {
	index, err := template.ParseFiles("views/user/register.html")
	if err != nil {
		panic(err)
	}
	index.Execute(c.Writer, "register")
}

// ToChat
// @Tags IM用户界面
// @Router /toChat [get]
func ToChat(c *gin.Context) {
	index, err := template.ParseFiles("views/chat/index.html",
		"views/chat/head.html",
		"views/chat/foot.html",
		"views/chat/concat.html",
		"views/chat/group.html",
		"views/chat/profile.html",
		"views/chat/tabmenu.html",
		"views/chat/userinfo.html",
		"views/chat/createcom.html",
		"views/chat/main.html")
	if err != nil {
		panic(err)
	}
	user := models.UserBasic{}
	userid, _ := strconv.Atoi(c.Query("userId"))
	token := c.Query("token")
	user.ID = uint(userid)
	user.Identity = token
	fmt.Println("current user >>>> userid:", userid, " token:", token)
	index.Execute(c.Writer, "chat")
}

func Chat(c *gin.Context) {
	models.Chat(c.Writer, c.Request)

}
