package constant

// // for app
// const (
// 	JWTRefreshExpiry = 129600 // jwt refresh token expiry in min // 90 day
// 	JWTAccessExpiry  = 43200  // jwt access token expiry in min // 30 day
// )

// corporate app
const (
	JWTRefreshExpiry = 4320 // (3 day = 4320) jwt refresh token expiry in min // 90 day
	JWTAccessExpiry  = 4320 // (3 day = 4320 e.g. 24*3*60) jwt access token expiry in min // 30 day
)

const (
	JWTAccessExpiryForNormalClient  = 10 // jwt access token expiry in 10 min
	JWTRefreshExpiryForNormalClient = 15 // jwt refresh token expiry in 15 min
)

// for website

const (
	JWTAccessExpiryForWeb  = 5    // jwt access token expiry in 5 min
	JWTRefreshExpiryForWeb = 1440 // jwt refresh token expiry in min // 1 day
)
