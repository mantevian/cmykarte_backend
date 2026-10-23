package types

import "database/sql"

type FunctionCallCtx struct {
	DB     *sql.DB
	UserId Uuid
}
