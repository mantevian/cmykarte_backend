package functions

import (
	"mantevian.xyz/cmykarte_backend/auth"
	"mantevian.xyz/cmykarte_backend/common/types"
)

func Register(req types.RegisterRequest, ctx types.FunctionCallCtx) types.RegisterResponse {
	if req.Password != req.PasswordConfirm {
		return types.RegisterError("Passwords don't match")
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return types.RegisterError(err.Error())
	}

	rows, err := ctx.DB.Query(
		"insert into users (id, name, password, created_at) values (gen_random_uuid(), $1, $2, now()) returning *",
		req.Name,
		hash,
	)

	if err != nil {
		return types.RegisterError(err.Error())
	}

	var user types.User

	rows.Next()
	rows.Scan(&user.Id, &user.Name, &user.Password, &user.CreatedAt)
	rows.Close()

	return types.RegisterResponse{
		Success:   true,
		UserId:    user.Id,
		Message:   "User registered successfully",
		CreatedAt: user.CreatedAt,
	}
}
