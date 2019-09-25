package courier

import (
	"context"
	"testing"

	"github.com/davecgh/go-spew/spew"

	"github.com/johnnyeven/libtools/courier/httpx"
)

type Group struct {
	EmptyOperator
}

type Auth struct{}

func (auth Auth) Output(ctx context.Context) (result interface{}, err error) {
	return
}

type Get struct {
	httpx.MethodGet
}

func (get Get) Path() string {
	return "/id"
}

func (get Get) Output() (result interface{}, err error) {
	return
}

func TestGroup(t *testing.T) {
	router := NewRouter(Group{})
	router.Register(NewRouter(Auth{}))

	spew.Dump(router)
	spew.Dump(router.Routes())
}
