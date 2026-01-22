package utils

func EnsureValidAppID(appID int) int {
	if appID > 0 {
		return appID
	}
	return 0
}
