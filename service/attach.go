package service

import (
	"fmt"
	"ginchat/utils"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

func Upload(c *gin.Context) {
	w := c.Writer
	req := c.Request
	srcFile, head, err := req.FormFile("file")
	if err != nil {
		utils.RespFail(w, err.Error())
		return
	}
	suffix := ""
	ofilName := head.Filename
	tem := strings.Split(ofilName, ".")
	if len(tem) > 1 {
		suffix = "." + tem[len(tem)-1]
	}
	// 若没有后缀（如无后缀文件），可默认设为.bin（可选，确保文件名合法）
	if suffix == "" {
		suffix = ".bin"
	}
	fileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)
	// 3. 确保上传目录存在（避免目录不存在导致创建文件失败）
	if err := os.MkdirAll("./asset/upload/", 0755); err != nil {
		utils.RespFail(w, "创建上传目录失败："+err.Error())
		return
	}
	dstFile, err := os.Create("./asset/upload/" + fileName)
	if err != nil {
		utils.RespFail(w, err.Error())
	}
	defer dstFile.Close() // 补充关闭文件句柄，避免资源泄露
	defer srcFile.Close() // 补充关闭源文件句柄
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		utils.RespFail(w, err.Error())
	}
	url := "./asset/upload/" + fileName
	fmt.Println("file url:", url)
	utils.RespOK(w, url, "发送成功")
}
