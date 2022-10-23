package errno

import (
	"fmt"
)

var ErrDuplicateMessage = fmt.Errorf("duplicate message key")
