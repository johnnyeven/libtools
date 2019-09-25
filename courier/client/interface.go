package client

import (
	"github.com/johnnyeven/libtools/courier"
)

type IRequest interface {
	Do() courier.Result
}
