package main

import (
	"ginchat/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	//"gorm.io/driver/sqlite"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:123456@tcp(127.0.0.1:3306)/ginchat?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	//db.AutoMigrate(&models.UserBasic{})
	//db.AutoMigrate(&models.Message{})
	db.AutoMigrate(&models.Community{})
	//db.AutoMigrate(&models.GroupBasic{})

	//// Create
	//user := &models.UserBasic{}
	//user.Name = "test"
	//db.Create(user)
	//
	//// Read
	//fmt.Println(db.First(user, 1))
	////db.First(&product, 1) // find product with integer primary key
	////db.First(&product, "code = ?", "D42") // find product with code D42
	//
	//// Update - update product's price to 200
	//db.Model(user).Update("PassWord", "1234")
	//// Update - update multiple fields
	////db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
	////db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})
	//
	// Delete - delete product
	//db.Delete(&product, 1)
}
