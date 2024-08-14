package constant

// // for app
// const (
// 	JWTRefreshExpiry = 129600 // jwt refresh token expiry in min // 90 day
// 	JWTAccessExpiry  = 43200  // jwt access token expiry in min // 30 day
// )

// corporate app
const (
	JWTRefreshExpiry = 43200 // (30 day = 43200) jwt refresh token expiry in min // 30 day
	JWTAccessExpiry  = 43200 // (30 day = 43200 e.g. 24*30*60) jwt access token expiry in min // 30 day
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
