package functions

import (
	"mantevian.xyz/cmykarte_backend/common/types"
	"mantevian.xyz/cmykarte_backend/game"
)

func Unqueue(req types.QueueRequest, ctx types.FunctionCallCtx) types.QueueResponse {
	res := game.RemoveQueueEntry(ctx.UserId)

	return types.QueueResponse{
		Success: res,
		Message: "",
	}
}
