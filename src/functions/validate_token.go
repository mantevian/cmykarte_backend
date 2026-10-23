package functions

import (
	"mantevian.xyz/cmykarte_backend/auth"
	"mantevian.xyz/cmykarte_backend/common/types"
)

func ValidateToken(req types.ValidateTokenRequest, ctx types.FunctionCallCtx) types.ValidateTokenResponse {
	claims, err := auth.ValidateToken(req.Token)

	if err != nil {
		return types.ValidateTokenFalse()
	}

	res := types.ValidateTokenResponse{
		Success: true,
		Claims:  *claims,
	}

	userId := res.Claims.UserId

	rows, err := ctx.DB.Query(
		"select id from users where id = $1",
		userId,
	)

	if err != nil {
		return types.ValidateTokenFalse()
	}

	if !rows.Next() {
		return types.ValidateTokenFalse()
	}

	rows.Close()

	return res
}
