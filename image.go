package brtool

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
)

// ImageType 探测图片的类型
// imgBytes 图片字节数组
func ImageType(imgBytes []byte) string {
	return http.DetectContentType(imgBytes)
}

func decodeImgBytesAndCreatFile(imgBytes []byte, filePath string) (image.Image, *os.File, error) {
	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return nil, nil, err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return nil, nil, err
	}
	return img, file, nil
}

// ImgBytesToImage []byte生成图片
// imgbytes 图片[]byte数组
// filepath 文件路径名称
func ImgBytesToImage(imgBytes []byte, filePath string) error {
	switch ImageType(imgBytes) {
	case "image/png":
		img, file, err := decodeImgBytesAndCreatFile(imgBytes, filePath)
		if err != nil {
			return err
		}
		defer file.Close()
		err = png.Encode(file, img)
		if err != nil {
			return err
		}

	case "image/jpeg":
		img, file, err := decodeImgBytesAndCreatFile(imgBytes, filePath)
		if err != nil {
			return err
		}
		defer file.Close()
		err = jpeg.Encode(file, img, nil)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown image type")
	}
	return nil
}

// ImageToBytes 图片转换为字节数组
// filepath 图片的路径
func ImageToImgBytes(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	ImgBytes := make([]byte, fileInfo.Size())
	buffer := bufio.NewReader(file)
	_, err = buffer.Read(ImgBytes)
	if err != nil {
		return nil, err
	}
	return ImgBytes, nil

}
