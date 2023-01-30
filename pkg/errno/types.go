// Package errno implements errors returned by gRPC. These errors are
// serialized and transmitted on the wire between server and client, and allow
// for additional data to be transmitted via the Details field in the status
package errno

import (
	"fmt"
)

// ErrDuplicateMessage is the error returned when a message is already in the queue.
var ErrDuplicateMessage = fmt.Errorf("duplicate message key")
