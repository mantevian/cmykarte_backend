package functions

import (
	"mantevian.xyz/cmykarte_backend/common/types"
)

func Queue(req types.QueueRequest, ctx types.FunctionCallCtx) types.QueueResponse {

	return types.QueueResponse{
		Success: true,
		Message: "hi",
	}
}
