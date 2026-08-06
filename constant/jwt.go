package constant

// // for app
// const (
// 	JWTRefreshExpiry = 129600 // jwt refresh token expiry in min // 90 day
// 	JWTAccessExpiry  = 43200  // jwt access token expiry in min // 30 day
// )

// corporate app
const (
	JWTRefreshExpiry = 43200 // (30 day = 43200) jwt refresh token expiry in min // 30 day
	JWTAccessExpiry  = 1440  // (1 day = 1440 e.g. 24*60) and (15 day = 21600 e.g. 24*60*15) jwt access token expiry in min // 30 day
)

const (
	JWTAccessExpiryForNormalClient  = 43200 // jwt access token expiry in 30 day
	JWTRefreshExpiryForNormalClient = 43200 // jwt refresh token expiry in 30 day
)

// for website

const (
	JWTAccessExpiryForWeb  = 1440 // jwt access token expiry in 1 day
	JWTRefreshExpiryForWeb = 1440 // jwt refresh token expiry in min // 1 day
)
