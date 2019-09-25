package httpx

import (
	"fmt"
	"net/http"
	"strings"
)

type IRedirect interface {
	Redirect(host string) string
	Status() int
}

func RedirectWithStatusMovedPermanently(uri string) *StatusMovedPermanently {
	return &StatusMovedPermanently{
		location: location(uri),
	}
}

type StatusMovedPermanently struct {
	location
}

func (s *StatusMovedPermanently) Status() int {
	return http.StatusMovedPermanently
}

func RedirectWithStatusFound(uri string) *StatusFound {
	return &StatusFound{
		location: location(uri),
	}
}

type StatusFound struct {
	location
}

func (s *StatusFound) Status() int {
	return http.StatusFound
}

func RedirectWithStatusSeeOther(uri string) *StatusSeeOther {
	return &StatusSeeOther{
		location: location(uri),
	}
}

type StatusSeeOther struct {
	location
}

func (s *StatusSeeOther) Status() int {
	return http.StatusSeeOther
}

func RedirectWithStatusTemporaryRedirect(uri string) *StatusTemporaryRedirect {
	return &StatusTemporaryRedirect{
		location: location(uri),
	}
}

type StatusTemporaryRedirect struct {
	location
}

func (s *StatusTemporaryRedirect) Status() int {
	return http.StatusTemporaryRedirect
}

func RedirectWithStatusPermanentRedirect(uri string) *StatusPermanentRedirect {
	return &StatusPermanentRedirect{
		location: location(uri),
	}
}

type StatusPermanentRedirect struct {
	location
}

func (s *StatusPermanentRedirect) Status() int {
	return http.StatusPermanentRedirect
}

type location string

func (r *location) Error() string {
	return fmt.Sprintf("location : %s", r)
}

func (r *location) Redirect(host string) string {
	u := string(*r)
	if strings.HasPrefix(u, "/") {
		u = host + u
		if !strings.HasPrefix(u, "http") {
			u = "http://" + u
		}
	}
	return u
}
