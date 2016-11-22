package services

import "fmt"

var (
	ErrCreatePassword         = fmt.Errorf("couldn't create a proper password")
	ErrUpdateUser             = fmt.Errorf("you don't have access to this user")
	ErrAppNotFound            = fmt.Errorf("app not found")
	ErrReleaseLocked          = fmt.Errorf("released is locked")
	ErrReleaseAlreadyLocked   = fmt.Errorf("release is already locked")
	ErrReleaseAlreadyUnlocked = fmt.Errorf("release is already unlocked")
)
