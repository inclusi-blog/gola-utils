package redis_util

type RedisStoreConfig struct {
	Host                  string `json:"host"`
	Port                  string `json:"port"`
	Db                    int    `json:"db"`
	ReadTimeoutInSeconds  int    `json:"read_timeout_in_seconds"`
	WriteTimeoutInSeconds int    `json:"write_timeout_in_seconds"`
	DialTimeoutInSeconds  int    `json:"dial_timeout_in_seconds"`
	Mode                  string `json:"mode"`
	Password              string `json:"password"`
}
