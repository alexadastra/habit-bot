package service

const (
	// internal errors
	checkinFailedErrorMessage       = "Error while saving 'checkin' action. Please try again."
	gratitudeFailedErrorMessage     = "Error while saving 'gratitude' action. Please try again."
	statsFetchingFailedErrorMessage = "Error while fetching stats. Please try again."
	stateFetchingFailedErrorMessage = "Error while fetching the user state. Please try again."
	stateSettingFailedErrorMessage  = "Error while setting the user state. Please try again."

	// user-ish errors
	invalidStateErrorMessage = "Invalid state. Please try again."

	// successful messages
	welcomeMessage = "Welcome to Habit Bot by @alex_ad_astra!"
)
