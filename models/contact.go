package models

import (
	"fmt"
	"ginchat/utils"
	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	OwnerId  uint //谁的通讯录
	TargetId uint //通讯录里有谁
	Type     int  //类型 1好友 2群
	Desc     string
}

func (table *Contact) TableName() string {
	return "contacts"
}

// 查找好友
func SearchFriend(userId uint) []UserBasic {
	contacts := make([]Contact, 0)
	utils.DB.Where("owner_id = ? and type = 1", userId).Find(&contacts)
	objIds := make([]uint, 0)
	for _, v := range contacts {
		fmt.Println(v)
		objIds = append(objIds, v.TargetId)
	}
	users := make([]UserBasic, 0)
	utils.DB.Where("id IN ?", objIds).Find(&users)
	return users
}

func AddFriend(userId uint, targetName string) (int, string) {
	user := &UserBasic{}
	if targetName != "" {
		user = FindUserByName(targetName)
		if user.ID == userId {
			return -1, "不能添加自己"
		}
		contact0 := Contact{}
		utils.DB.Where("owner_id = ? and target_id = ? and type = 1", userId, user.ID).First(&contact0)
		if contact0.ID != 0 {
			return -1, "好友已存在"
		}
		if user.Salt != "" {
			tx := utils.DB.Begin()
			//事务中出现任何异常均Rollback
			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
				}
			}()
			contact := Contact{}
			contact.OwnerId = userId
			contact.TargetId = user.ID
			contact.Type = 1
			if err := utils.DB.Create(&contact).Error; err != nil {
				tx.Rollback()
				return -1, "添加失败"
			}

			contact1 := Contact{}
			contact1.OwnerId = user.ID
			contact1.TargetId = userId
			contact1.Type = 1
			if err := utils.DB.Create(&contact1).Error; err != nil {
				tx.Rollback()
				return -1, "添加失败"
			}

			tx.Commit()
			return 0, "添加好友成功"
		}
	}
	return -1, "用户不存在"
}

func SearchUserByGroupId(communityId uint) []uint {
	contacts := make([]Contact, 0)
	objIds := make([]uint, 0)
	utils.DB.Where("target_id = ? and type=2", communityId).Find(&contacts)
	for _, v := range contacts {
		objIds = append(objIds, uint(v.OwnerId))
	}
	return objIds
}
