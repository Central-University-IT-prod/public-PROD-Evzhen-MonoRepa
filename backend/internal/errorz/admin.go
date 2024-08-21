package errorz

import "errors"

var (
	ErrorInvalidLogin = errors.New("invalid login")
	AdminNotFound     = errors.New("admin not found")
	PasswordIncorrect = errors.New("password incorrect")

	ErrorWrongTimeInterval = errors.New("wrong time interval")

	ErrorContestNotFound = errors.New("contest not found")

	ErrorTeamLimitAlreadyExists = errors.New("team limit already exists")
	ErrorMissingAccess          = errors.New("missing access")
)
