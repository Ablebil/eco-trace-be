package response

// Auth Domain
const (
	EmailAlreadyExists  = "Email already exists"
	UserNotFound        = "User not found"
	InvalidOTP          = "Invalid or expired OTP"
	InvalidCredentials  = "Invalid email or password"
	InvalidRefreshToken = "Invalid or expired refresh token"
	UserNotVerified     = "User not verified"
	OAuthStateNotFound  = "OAuth state not found"
	OAuthStateInvalid   = "OAuth state invalid"

	FailedFindUser           = "Failed to find user"
	FailedCreateUser         = "Failed to create user"
	FailedUpdateUser         = "Failed to update user"
	FailedAddRefreshToken    = "Failed to add refresh token"
	FailedGetRefreshTokens   = "Failed to get refresh tokens"
	FailedRemoveRefreshToken = "Failed to remove refresh token"
	FailedExchangeOAuthToken = "Failed to exchange OAuth token"
	FailedGetOAuthProfile    = "Failed to get OAuth profile"
	FailedGoogleLogin        = "Failed to initiate Google login"

	RegisterSuccess     = "Registration successful. OTP has been sent to email"
	VerifyOTPSuccess    = "Verification successful"
	LoginSuccess        = "Login successful"
	RefreshTokenSuccess = "Token refresh successful"
	LogoutSuccess       = "Logout successful"
)

// challenge Domain
const (
	ChallengeNotFound         = "Challenge not found"
	ChallengeAlreadyTaken     = "Challenge already taken"
	ChallengeNotTaken         = "Challenge not taken by user"
	ChallengeAlreadyCompleted = "Challenge already completed"
	ChallengeNotActive        = "Challenge is not active"

	FailedGetChallenges     = "Failed to get challenges"
	FailedGetUserChallenges = "Failed to get user challenges"
	FailedTakeChallenge     = "Failed to take challenge"
	FailedCompleteChallenge = "Failed to complete challenge"
	FailedUpdateUserExp     = "Failed to update user experience"
	FailedGetBadges         = "Failed to get badges"
	FailedGetUserBadges     = "Failed to get user badges"
	FailedUnlockBadge       = "Failed to unlock badge"

	TakeChallengeSuccess     = "Challenge taken successfully"
	CompleteChallengeSuccess = "Challenge completed successfully"
	BadgeUnlockedSuccess     = "New badge unlocked!"
)

// Others
const (
	FailedHashPassword         = "Failed to hash password"
	FailedGenerateOTP          = "Failed to generate OTP"
	FailedStoreOTP             = "Failed to store OTP"
	FailedDeleteOTP            = "Failed to delete OTP"
	FailedSendOTPEmail         = "Failed to send OTP email"
	FailedGenerateRefreshToken = "Failed to generate refresh token"
	FailedGenerateAccessToken  = "Failed to generate access token"
	FailedGenerateOAuthState   = "Failed to generate OAuth state"
	FailedStoreOAuthState      = "Failed to store OAuth state"
	FailedDeleteOAuthState     = "Failed to delete OAuth state"
	FailedGenerateOAuthLink    = "Failed to generate OAuth link"
	FailedOAuthCallback        = "Failed to handle OAuth callback"
)

// Handler
const (
	FailedParsingRequestBody    = "Failed parsing request body"
	FailedParsingRequestParams  = "Failed parsing request params"
	FailedValidateRequest       = "Failed to validate request"
	MissingAccessToken          = "Missing access token"
	InvalidAccessToken          = "Invalid access token"
	InvalidOrMissingBearerToken = "Invalid or missing bearer token"
)
