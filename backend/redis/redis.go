package redis

type Redis struct {
	tb map[int]string
}

func New() *Redis {
	return &Redis{
		tb: map[int]string{},
	}
}

func (redis *Redis) Put(key int, value string) {
	redis.tb[key] = value
}

func (redis *Redis) Get(key int) string {
	value, exists := redis.tb[key]
	if !exists {
		return ""
	}
	return value
}
