package ucerr

type UCErr struct {
	msg string
}

func (e *UCErr) Error() string {
	return e.msg
}

var (
	ErrNoUser           = &UCErr{msg: "no such user"}
	ErrNoChat           = &UCErr{msg: "no such chat"}
	ErrNameAlreadyInUse = &UCErr{msg: "name already in use"}
	ErrUserInChatTwice  = &UCErr{msg: "user can't be in one chat twice"}
	ErrUserNotInChat    = &UCErr{msg: "user not in this chat"}
)
