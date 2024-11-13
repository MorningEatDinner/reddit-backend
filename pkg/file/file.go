package file

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/xiaorui/reddit-async/reddit-backend/pkg/helpers"
)

// SaveUploadAvatar：保存用户头像
func SaveUploadAvatar(ctx *gin.Context, file *multipart.FileHeader, userID int64) (path string, err error) {
	publicPath := "public"
	dirName := fmt.Sprintf("/uploads/avatar/%s/%s/",
		time.Now().In(time.Local), fmt.Sprintf("%v", userID))
	os.MkdirAll(publicPath+dirName, 0755)

	// 保存文件
	filename := helpers.RandomString(16) + filepath.Ext(file.Filename)
	avatarPath := publicPath + dirName + filename
	if err := ctx.SaveUploadedFile(file, avatarPath); err != nil {
		return "", err
	}

	// 修剪图片
	img, err := imaging.Open(avatarPath, imaging.AutoOrientation(true))
	if err != nil {
		return "", err
	}
	resizeAvatar := imaging.Thumbnail(img, 256, 256, imaging.Lanczos)
	resizeFilename := helpers.RandomString(16) + filepath.Ext(file.Filename)
	resizeAvatarPath := publicPath + dirName + resizeFilename
	err = imaging.Save(resizeAvatar, resizeAvatarPath)
	if err != nil {
		return "", err
	}

	// 如果剪切成功则删除旧的文件
	err = os.Remove(avatarPath)
	if err != nil {
		return "", err
	}

	return avatarPath, nil
}
