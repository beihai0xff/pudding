// Package types provides the types of the broker.
package types

import (
	"context"

	"github.com/beihai0xff/pudding/api/gen/pudding/types/v1"
)

// HandleMessage is a function that handles a message.
type HandleMessage func(ctx context.Context, msg *types.Message) error
