package util

import (
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	LOGGER "salbackend/logger"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const randomIDdigits = "abcdefghijklmnopqrstuvwxyz0123456789"

func generateRandomID(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = randomIDdigits[rand.Intn(len(randomIDdigits))]
	}
	return string(b)
}

func saveToDisk(file multipart.File, extension string) (string, bool) {
	LOGGER.Log("Saving to disk...")

	fileName := "/tmp/" + generateRandomID(10) + extension

	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println("saveToDisk", err)
		return "", false
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		fmt.Println("saveToDisk", err)
		return "", false
	}

	return fileName, true
}

func UploadToS3(s3Bucket, path, s3AccessKey, s3SecretKey, s3Region, extension, acl string, file multipart.File) (string, bool) {

	savedFileName, saved := saveToDisk(file, extension)
	if !saved {
		return "", false
	}
	fmt.Println(savedFileName)

	conf := aws.Config{
		Credentials: credentials.NewStaticCredentials(s3AccessKey, s3SecretKey, ""),
		Region:      aws.String(s3Region),
	}
	sess := session.New(&conf)

	svc := s3manager.NewUploader(sess)

	LOGGER.Log("Uploading file to S3...")

	openedFile, err := os.Open(savedFileName)
	if err != nil {
		fmt.Println("uploadToS3", err)
		return "", false
	}
	s, _ := openedFile.Stat()
	fmt.Println(s.Size())
	defer openedFile.Close()

	fileName := path + "/" + getMD5Hash(savedFileName) + extension

	_, err = svc.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(s3Bucket),
		Key:         aws.String(fileName),
		Body:        openedFile,
		ContentType: aws.String(getFileMIMEType(strings.ToLower(extension))),
		ACL:         aws.String(acl),
	})

	os.Remove(savedFileName)
	if err != nil {
		fmt.Println("uploadToS3 "+path, err)
		return "", false
	}
	return fileName, true
}

func getFileMIMEType(extension string) string {
	switch extension {
	// video
	case ".mp2":
		return "video/mpeg"
	case ".mpa":
		return "video/mpeg"
	case ".mpe":
		return "video/mpeg"
	case ".mpeg":
		return "video/mpeg"
	case ".mpg":
		return "video/mpeg"
	case ".mpv2":
		return "video/mpeg"
	case ".mp4":
		return "video/mp4"
	case ".mov":
		return "video/quicktime"
	case ".qt":
		return "video/quicktime"
	case ".lsf":
		return "video/x-la-asf"
	case ".lsx":
		return "video/x-la-asf"
	case ".asf":
		return "video/x-ms-asf"
	case ".asr":
		return "video/x-ms-asf"
	case ".asx":
		return "video/x-ms-asf"
	case ".avi":
		return "video/x-msvideo"
	case ".movie":
		return "video/x-sgi-movie"
	case ".3gp":
		return "video/3gpp"
	case ".3gpp":
		return "video/3gpp"
	case ".3gpp2":
		return "video/3gpp2"
	case ".3g2":
		return "video/3gpp2"
	// image
	case ".bmp":
		return "image/bmp"
	case ".cod":
		return "image/cis-cod"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".ief":
		return "image/ief"
	case ".jpe":
		return "image/jpeg"
	case ".jpeg":
		return "image/jpeg"
	case ".jpg":
		return "image/jpeg"
	case ".jfif":
		return "image/pipeg"
	case ".svg":
		return "image/svg+xml"
	case ".tif":
		return "image/tiff"
	case ".tiff":
		return "image/tiff"
	case ".ras":
		return "image/x-cmu-raster"
	case ".cmx":
		return "image/x-cmx"
	case ".ico":
		return "image/x-icon"
	case ".pnm":
		return "image/x-portable-anymap"
	case ".pbm":
		return "image/x-portable-bitmap"
	case ".pgm":
		return "image/x-portable-graymap"
	case ".ppm":
		return "image/x-portable-pixmap"
	case ".rgb":
		return "image/x-rgb"
	case ".xbm":
		return "image/x-xbitmap"
	case ".xpm":
		return "image/x-xpixmap"
	case ".xwd":
		return "image/x-xwindowdump"
	default:
		return ""
	}
}
