package types

import (
	"context"

	"github.com/beihai0xff/pudding/api/gen/pudding/types/v1"
)

type HandleMessage func(ctx context.Context, msg *types.Message) error
