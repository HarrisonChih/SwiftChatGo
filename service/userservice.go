package service

import (
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// GetUserList
// @Summary 查看所有用户
// @Tags 用户模块
// @Success 200 {string} json{"code", "massage", "data"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()

	c.JSON(200, gin.H{
		"code":    200,
		"message": "查找成功",
		"data":    data,
	})
}

// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repasswd query string false "确认密码"
// @Success 200 {string} json{"code", "massage", "data"}
// @Router /user/createUser [get]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.Request.FormValue("name")
	data := models.FindUserByName(user.Name)
	if data.Name != "" {
		c.JSON(400, gin.H{
			"code":    -1,
			"message": "用户名已注册",
		})
		return
	}

	// 生成 0 ~ 999999 之间的随机整数，确保是 6 位以内的非负数
	source := rand.NewSource(time.Now().UnixNano()) // 初始化随机数种子（只需要在程序启动时执行一次）
	randnum := rand.New(source)
	salt := fmt.Sprintf("%06d", randnum.Intn(1000000)) // rand.Intn(1000000) 生成 0~999999 的随机数
	user.Salt = salt

	password := c.Request.FormValue("password")
	repasswd := c.Request.FormValue("repasswd")

	if user.Name == "" || password == "" {
		c.JSON(400, gin.H{
			"code":    -1,
			"message": "用户名或密码不能为空",
		})
		return
	}

	if password != repasswd {
		c.JSON(400, gin.H{
			"code":    -1,
			"message": "两次输入的密码不一致",
		})
		return
	}
	user.PassWord = utils.MakePassword(password, user.Salt) //使用SHA-256加密
	models.CreateUser(user)
	c.JSON(200, gin.H{
		"code":    0,
		"message": "账号创建成功",
		"data":    user,
	})
}

// DeleteUser
// @Summary 注销用户
// @Tags 用户模块
// @param id query string false "id"
// @Success 200 {string} json{"code", "massage", "data"}
// @Router /user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(200, gin.H{
		"code":    0,
		"message": "注销成功",
		"data":    user,
	})
}

// UpdateUser
// @Summary 修改用户信息
// @Tags 用户模块
// @param id formData string false "id"
// @param name formData string false "name"
// @param password formData string false "password"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Success 200 {string} json{"code", "massage", "data"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")

	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{
			"code":    -1,
			"message": "输入格式有误，请重新检查",
		})
		return
	}

	models.UpdateUser(user)
	c.JSON(200, gin.H{
		"code":    0,
		"message": "修改成功",
		"data":    user,
	})
}

// FindUserByNameAndPwd
// @Summary 用户登录校验
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success 200 {string} json{"code", "massage", "data"}
// @Router /user/findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {
	name := c.Request.FormValue("name")
	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(400, gin.H{
			"code":    -1,
			"message": "用户名错误，请重新输入",
			"data":    nil,
		})
		return
	}
	password := c.Request.FormValue("password")
	flag := utils.ValidPassword(password, user.Salt, user.PassWord)
	if !flag {
		c.JSON(400, gin.H{
			"code":    -1,
			"message": "密码错误，请重新输入",
			"data":    nil,
		})
		return
	}
	pwd := utils.MakePassword(password, user.Salt)
	data := models.FindUserByNameAndPwd(name, pwd)

	c.JSON(200, gin.H{
		"code":    0,
		"message": "登陆成功",
		"data":    data,
	})
}

// 防止跨域站点伪造请求
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	// 1. 升级 HTTP 连接为 WebSocket 连接
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 2. 延迟关闭连接（函数退出时关闭，避免资源泄露）
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	// 3. 消息处理（订阅 Redis 消息并推送给客户端）
	MsgHandler(ws, c)
}

func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	for { // 1. 订阅 Redis 频道（utils.PublishKey 是订阅的频道名称
		msg, err := utils.Subscribe(c, utils.PublishKey)
		if err != nil {
			fmt.Println(err)
		}
		// 2. 格式化消息（添加时间戳）
		tm := time.Now().Format("2006-01-02 15:04:05")
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
		// 3. 通过 WebSocket 向客户端推送消息
		err = ws.WriteMessage(websocket.TextMessage, []byte(m))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}

func RedisMsg(c *gin.Context) {
	userIdA, _ := strconv.Atoi(c.PostForm("userIdA"))
	userIdB, _ := strconv.Atoi(c.PostForm("userIdB"))
	start, _ := strconv.Atoi(c.PostForm("start"))
	end, _ := strconv.Atoi(c.PostForm("end"))
	isRev, _ := strconv.ParseBool(c.PostForm("isRev"))
	res := models.RedisMsg(int64(userIdA), int64(userIdB), int64(start), int64(end), isRev)
	//res := models.RedisMsg(int64(userIdA), int64(userIdB))
	utils.RespOKList(c.Writer, "ok", res)
}

// SearchFriends
// @Summary 查询好友列表
// @Tags 聊天界面
// @param userId query string false "用户ID"
// @Success 200 {string} json{"code, massage, data"}
// @Router /searchFriends [post]
func SearchFriends(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	users := models.SearchFriend(uint(userId))
	utils.RespOKList(c.Writer, users, len(users))
}

// AddFriend
// @Summary 添加好友
// @Tags 聊天界面
// @param userId query string false "用户ID"
// @param targetId query string false "好友ID"
// @Success 200 {string} json{"code, massage, data"}
// @Router /contact/addfriend [post]
func AddFriend(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	targetName := c.Request.FormValue("targetName")
	code, msg := models.AddFriend(uint(userId), targetName)
	if code == 0 {
		utils.RespOK(c.Writer, code, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// CreateCommunity
// @Summary 新建群聊
// @Tags 群聊界面
// @param ownerId query string false "用户ID"
// @param name query string false "群名称"
// @param icon query string false "群头像"
// @param desc query string false "群描述"
// @Success 200 {string} json{"code, massage, data"}
// @Router /contact/createCommunity [post]
func CreateCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId"))
	name := c.Request.FormValue("name")
	icon := c.Request.FormValue("icon")
	desc := c.Request.FormValue("desc")
	community := models.Community{}
	community.OwnerId = uint(ownerId)
	community.Name = name
	community.Img = icon
	community.Desc = desc
	code, msg := models.CreateCommunity(community)
	if code == 0 {
		utils.RespOK(c.Writer, code, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// LoadCommunity
// @Summary 加载群列表
// @Tags 群聊界面
// @param ownerId query string false "用户ID"
// @Success 200 {string} json{"code, massage, data"}
// @Router /contact/loadcommunity [post]
func LoadCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId"))
	data, msg := models.LoadCommunity(uint(ownerId))
	if len(data) != 0 {
		utils.RespList(c.Writer, 0, data, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// JoinGroups
// @Summary 加群
// @Tags 群聊界面
// @param userId query string false "用户ID"
// @param comId query string false "群ID"
// @Success 200 {string} json{"code, massage, data"}
// @Router /contact/joinGroup [post]
func JoinGroups(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	comId := c.Request.FormValue("comId")

	//	name := c.Request.FormValue("name")
	data, msg := models.JoinGroup(uint(userId), comId)
	if data == 0 {
		utils.RespOK(c.Writer, data, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

func FindByID(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	data := models.FindByID(uint(userId))
	utils.RespOK(c.Writer, data, "ok")
}
