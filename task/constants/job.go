package constants

type TaskProcessor func(*Task) (interface{}, error)
