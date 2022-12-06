package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// LoadConfig - load .env file from given path for local, else will be getting from env var
func LoadConfig() {
	// load .env file from given path for local, else will be getting from env var
	if !strings.EqualFold(os.Getenv("prod"), "true") {
		err := godotenv.Load(".test-env")
		if err != nil {
			panic("Error loading .env file")
		}
	}

	DBConfig = os.Getenv("DB_CONFIG")
	DBConnectionPool, _ = strconv.Atoi(os.Getenv("DB_CONNECTION_POOL"))
	Log, _ = strconv.ParseBool(os.Getenv("LOG"))
	Migrate, _ = strconv.ParseBool(os.Getenv("MIGRATE"))
	RazorpayAuth = os.Getenv("RAZORPAY_AUTH")
	RazorpaySecret = os.Getenv("RAZORPAY_SECRET")
	OneSignalAppIDForClient = os.Getenv("ONESIGNAL_APP_ID_FOR_CLIENT")
	OneSignalAppIDForTherapist = os.Getenv("ONESIGNAL_APP_ID_FOR_THERAPIST")
	S3Bucket = os.Getenv("S3_BUCKET")
	MediaURL = os.Getenv("MEDIA_URL")
	AWSAccesKey = os.Getenv("AWSACCESSKEY")
	AWSSecretKey = os.Getenv("AWSSECRETKEY")
	AWSRegion = os.Getenv("AWSREGION")
	JWTSecret = []byte(os.Getenv("JWT_SECRET"))
	S3BUCKETAGORA = os.Getenv("S3_BUCKET_AGORA")
	MediaURLFORAGORA = os.Getenv("Media_URL_FOR_AGORA")
	AGORA_APP_ID = os.Getenv("AGORA_APP_ID")
	AGORA_APP_CERTIFICATE = os.Getenv("AGORA_APP_CERTIFICATE")
	AGORA_Customer_Key = os.Getenv("AGORA_CUSTOMER_KEY")
	AGORA_Customer_Secret = os.Getenv("AGORA_CUSTOMER_SECRET")
	FromEmailID = os.Getenv("FROM_EMAIL_ID")
}
