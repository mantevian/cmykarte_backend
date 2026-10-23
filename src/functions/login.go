package functions

import (
	"mantevian.xyz/cmykarte_backend/auth"
	"mantevian.xyz/cmykarte_backend/common/types"
)

func Login(req types.LoginRequest, ctx types.FunctionCallCtx) types.LoginResponse {
	rows, err := ctx.DB.Query(
		"select id, name, password, created_at from users where name = $1",
		req.Name,
	)

	if err != nil {
		return types.LoginError(err.Error())
	}

	var user types.User

	rows.Next()
	rows.Scan(&user.Id, &user.Name, &user.Password, &user.CreatedAt)
	rows.Close()

	if !auth.CheckPassword(req.Password, user.Password) {
		return types.LoginError("Incorrect credentials")
	}

	token, err := auth.GenerateToken(user.Id, user.Name)
	if err != nil {
		return types.LoginError(err.Error())
	}

	return types.LoginResponse{
		Success: true,
		UserId:  user.Id,
		Message: "User logged in successfully",
		Token:   token,
	}
}
