package constant

import "server/internal/types"

var (
	EnvConfig = new(types.Env)

	ServerConfig = new(types.Server)

	LogConfig = new(types.Log)

	DBConfig = new(types.DB)

	RedisConfig = new(types.Redis)

	JWTConfig = new(types.JWT)

	FileConfig = new(types.File)
)
