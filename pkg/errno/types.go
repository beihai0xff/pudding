package errno

import (
	"fmt"
)

// ErrDuplicateMessage is the error returned when a message is already in the queue.
var ErrDuplicateMessage = fmt.Errorf("duplicate message key")
