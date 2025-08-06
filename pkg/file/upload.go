package file

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// SaveUploadedFile 保存上传的文件
// 支持自动创建目标文件夹
func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	// 创建目标文件夹
	dir := filepath.Dir(dst)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
