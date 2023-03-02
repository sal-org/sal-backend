package constant

// for app
const (
	JWTRefreshExpiry = 129600 // jwt refresh token expiry in min // 90 day
	JWTAccessExpiry  = 43200  // jwt access token expiry in min // 30 day
)

// for website

const (
	JWTAccessExpiryForWeb  = 5    // jwt access token expiry in 5 min
	JWTRefreshExpiryForWeb = 1440 // jwt refresh token expiry in min // 1 day
)
