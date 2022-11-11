package global

import (
	"errors"
	"fmt"
)

var ErrTimeout = errors.New("timeout")
var ErrCtxCancel = errors.New("ctx cancel")
var ErrWaitUserSync = fmt.Errorf("Synchronizing user data, please wait...")
var ErrUserAskTimeout = fmt.Errorf("do ask timeout.....")
var ErrUserAskWriteAllFail = fmt.Errorf("do ask write all fail.....")
var ErrTweetNonceConflict = fmt.Errorf("tweet nonce conflict")
