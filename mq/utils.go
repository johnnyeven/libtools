package mq

func Prefix(key string, prefix string) string {
	return prefix + key
}
