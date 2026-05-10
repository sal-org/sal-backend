package util

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	LOGGER "salbackend/logger"
	"strings"
	"time"

	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const randomIDdigits = "abcdefghijklmnopqrstuvwxyz0123456789"

const randomIDdigitsWithCap = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

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

func saveToDiskFile(file []byte, extension string) (string, bool) {
	LOGGER.Log("Saving to disk...")

	newfile := bytes.NewReader(file)

	fileName := "/tmp/" + generateRandomID(10) + extension

	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println("saveToDisk", err)
		return "", false
	}
	defer f.Close()
	_, err = io.Copy(f, newfile)
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

	// conf := aws.Config{
	// 	Credentials: credentials.NewStaticCredentials(s3AccessKey, s3SecretKey, ""),
	// 	Region:      aws.String(s3Region),
	// }
	// sess, _ := session.NewSession(&conf)

	// svc := s3manager.NewUploader(sess)

	// LOGGER.Log("Uploading file to S3...")

	// openedFile, err := os.Open(savedFileName)
	// if err != nil {
	// 	fmt.Println("uploadToS3", err)
	// 	return "", false
	// }
	// s, _ := openedFile.Stat()
	// fmt.Println(s.Size())
	// defer openedFile.Close()

	// fileName := path + "/" + getFileMD5Hash(savedFileName) + extension

	// _, err = svc.Upload(&s3manager.UploadInput{
	// 	Bucket:      aws.String(s3Bucket),
	// 	Key:         aws.String(fileName),
	// 	Body:        openedFile,
	// 	ContentType: aws.String(getFileMIMEType(strings.ToLower(extension))),
	// 	// ACL:         aws.String(acl),
	// })

	// os.Remove(savedFileName)
	// if err != nil {
	// 	fmt.Println("uploadToS3 "+path, err)
	// 	return "", false
	// }
	// return fileName, true

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Load AWS config (replaces aws.Config + session.NewSession)
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(s3Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				s3AccessKey,
				s3SecretKey,
				"",
			),
		),
		config.WithRetryMaxAttempts(3),
	)
	if err != nil {
		fmt.Println("failed to load aws config:", err)
		return "", false
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg)

	LOGGER.Log("Uploading file to S3...")

	// Open file
	openedFile, err := os.Open(savedFileName)
	if err != nil {
		fmt.Println("uploadToS3 open file error:", err)
		return "", false
	}
	defer openedFile.Close()

	// File info
	fileStat, err := openedFile.Stat()
	if err == nil {
		fmt.Println("file size:", fileStat.Size())
	}

	// Generate file name
	fileName := path + "/" + getFileMD5Hash(savedFileName) + extension

	// Upload to S3 using PutObject
	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(fileName),
		Body:   openedFile,
		ContentType: aws.String(
			getFileMIMEType(strings.ToLower(extension)),
		),

		// Optional:
		// ServerSideEncryption: types.ServerSideEncryptionAes256,
	})

	if err != nil {
		fmt.Println("uploadToS3 upload error:", err)
		return "", false
	}

	// Remove local file after successful upload
	err = os.Remove(savedFileName)
	if err != nil {
		fmt.Println("warning: failed to remove local file:", err)
	}

	fmt.Println("file uploaded successfully:", fileName)

	return fileName, true

}

func UploadToS3File(s3Bucket, path, s3AccessKey, s3SecretKey, s3Region, extension, acl string, file []byte) (string, bool) {

	savedFileName, saved := saveToDiskFile(file, extension)
	if !saved {
		return "", false
	}
	fmt.Println(savedFileName)

	// conf := aws.Config{
	// 	Credentials: credentials.NewStaticCredentials(s3AccessKey, s3SecretKey, ""),
	// 	Region:      aws.String(s3Region),
	// }
	// sess, _ := session.NewSession(&conf)

	// svc := s3manager.NewUploader(sess)

	// LOGGER.Log("Uploading file to S3...")

	// openedFile, err := os.Open(savedFileName)
	// if err != nil {
	// 	fmt.Println("uploadToS3", err)
	// 	return "", false
	// }
	// s, _ := openedFile.Stat()
	// fmt.Println(s.Size())
	// defer openedFile.Close()

	// fileName := path + "/" + getFileMD5Hash(savedFileName) + extension

	// _, err = svc.Upload(&s3manager.UploadInput{
	// 	Bucket:      aws.String(s3Bucket),
	// 	Key:         aws.String(fileName),
	// 	Body:        openedFile,
	// 	ContentType: aws.String(getFileMIMEType(strings.ToLower(extension))),
	// 	ACL:         aws.String(acl),
	// })

	// os.Remove(savedFileName)
	// if err != nil {
	// 	fmt.Println("uploadToS3 "+path, err)
	// 	return "", false
	// }
	// return fileName, true

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Load AWS config (replaces aws.Config + session.NewSession)
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(s3Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				s3AccessKey,
				s3SecretKey,
				"",
			),
		),
		config.WithRetryMaxAttempts(3),
	)
	if err != nil {
		fmt.Println("failed to load aws config:", err)
		return "", false
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg)

	LOGGER.Log("Uploading file to S3...")

	// Open file
	openedFile, err := os.Open(savedFileName)
	if err != nil {
		fmt.Println("uploadToS3 open file error:", err)
		return "", false
	}
	defer openedFile.Close()

	// File info
	fileStat, err := openedFile.Stat()
	if err == nil {
		fmt.Println("file size:", fileStat.Size())
	}

	// Generate file name
	fileName := path + "/" + getFileMD5Hash(savedFileName) + extension

	// Upload to S3 using PutObject
	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(fileName),
		Body:   openedFile,
		ContentType: aws.String(
			getFileMIMEType(strings.ToLower(extension)),
		),

		// Optional:
		// ServerSideEncryption: types.ServerSideEncryptionAes256,
	})

	if err != nil {
		fmt.Println("uploadToS3 upload error:", err)
		return "", false
	}

	// Remove local file after successful upload
	err = os.Remove(savedFileName)
	if err != nil {
		fmt.Println("warning: failed to remove local file:", err)
	}

	fmt.Println("file uploaded successfully:", fileName)

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

func generateRandomIDForSignedURl(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = randomIDdigitsWithCap[rand.Intn(len(randomIDdigitsWithCap))]
	}
	return string(b)
}

func PreSignedS3URLToUploadPut(s3Bucket, path, s3AccessKey, s3SecretKey, s3Region, extension string) (string, string) {

	// // Initialize a session in us-west-2 that the SDK will use to load
	// // credentials from the shared credentials file ~/.aws/credentials.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(s3Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				s3AccessKey,
				s3SecretKey,
				"",
			),
		),
	)
	if err != nil {
		fmt.Println("failed to load aws config:", err)
		return "", ""
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg)

	// Generate unique file key
	fileNameKey := path + "/" + generateRandomIDForSignedURl(30) + extension

	// Create presign client
	presignClient := s3.NewPresignClient(client)

	// Generate presigned PUT URL
	request, err := presignClient.PresignPutObject(
		ctx,
		&s3.PutObjectInput{
			Bucket: aws.String(s3Bucket),
			Key:    aws.String(fileNameKey),
			ContentType: aws.String(
				getFileMIMEType(strings.ToLower(extension)),
			),

			// Optional:
			// ACL: types.ObjectCannedACLPublicRead,
		},
		s3.WithPresignExpires(5*time.Minute),
	)

	if err != nil {
		fmt.Println("failed to generate presigned url:", err)
		return "", ""
	}

	// request.URL contains the signed upload URL
	return request.URL, fileNameKey

}

func PreSignedS3URLToGetTheData(s3Bucket, path, s3AccessKey, s3SecretKey, s3Region string) string {

	// // credentials from the shared credentials file ~/.aws/credentials.
	// Create AWS session
	if path == "" {
		return ""
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(s3Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				s3AccessKey,
				s3SecretKey,
				"",
			),
		),
	)
	if err != nil {
		fmt.Println("config error:", err)
		return ""
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg)

	// Create presign client
	presignClient := s3.NewPresignClient(client)

	// Generate presigned GET URL
	request, err := presignClient.PresignGetObject(
		ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(s3Bucket),
			Key:    aws.String(path),
		},
		s3.WithPresignExpires(15*time.Minute),
	)

	if err != nil {
		fmt.Println("presign error:", err)
		return ""
	}

	// Return signed URL
	return request.URL
}
