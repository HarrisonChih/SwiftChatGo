package models

import (
	"fmt"
	"ginchat/utils"
	"gorm.io/gorm"
	"time"
)

type UserBasic struct {
	gorm.Model
	Name          string
	PassWord      string
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email         string `valid:"email"`
	Identity      string
	ClientIp      string
	ClientPort    string
	Salt          string
	LoginTime     time.Time `gorm:"autoCreateTime"`
	HeartbeatTime time.Time `gorm:"autoCreateTime"`
	LoginOutTime  time.Time `gorm:"column:login_out_time; autoCreateTime" json:"login_out_time"`
	IsLogOut      bool
	DeviceInfo    string
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

func FindUserByNameAndPwd(name, password string) *UserBasic {
	user := &UserBasic{}
	utils.DB.Where("name = ? AND pass_word = ?", name, password).First(user)

	//token 加密
	temp := fmt.Sprintf("%d", time.Now().Unix())
	tokenEncrypt := utils.Sha256Encode(temp)

	utils.DB.Model(user).Where("id = ?", user.ID).Update("identity", tokenEncrypt)
	return user
}

func FindUserByID(userId uint) *UserBasic {
	user := &UserBasic{}
	utils.DB.Where("id = ?", userId).First(user)
	return user
}

func FindUserByName(name string) *UserBasic {
	user := &UserBasic{}
	utils.DB.Where("name = ?", name).First(user)
	return user
}

func FindUserByEmail(email string) *UserBasic {
	user := &UserBasic{}
	utils.DB.Where("email = ?", email).First(user)
	return user
}

func FindUserByPhone(phone string) *UserBasic {
	user := &UserBasic{}
	utils.DB.Where("phone = ?", phone).First(user)
	return user
}

func CreateUser(user UserBasic) *gorm.DB {
	return utils.DB.Create(&user)
}

func DeleteUser(user UserBasic) *gorm.DB {
	return utils.DB.Delete(&user)
}

func UpdateUser(user UserBasic) *gorm.DB {
	return utils.DB.Model(&user).Updates(UserBasic{Name: user.Name, PassWord: user.PassWord, Phone: user.Phone, Email: user.Email})
}

// 查找某个用户
func FindByID(id uint) UserBasic {
	user := UserBasic{}
	utils.DB.Where("id = ?", id).First(&user)
	return user
}
