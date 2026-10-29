package functions

import (
	"mantevian.xyz/cmykarte_backend/common/types"
	"mantevian.xyz/cmykarte_backend/game"
)

func Queue(req types.QueueRequest, ctx types.FunctionCallCtx) types.QueueResponse {
	res := game.AddQueueEntry(ctx.UserId, ctx.Username)

	return types.QueueResponse{
		Success: res,
		Message: "",
	}
}
