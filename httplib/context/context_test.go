package context

import (
	"bytes"
	"net/http"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func getNormalContext() *gin.Context {
	c := gin.Context{}
	request, _ := http.NewRequest("GET", "http://127.0.0.1/testgo/v0/sum?num1=1&num2=2", nil)
	c.Request = request
	rwiter := TestResponseWriter{}
	rwiter.Body = *bytes.NewBuffer([]byte{})
	c.Writer = &rwiter
	c.Writer.Write([]byte{'1'})
	return &c
}

func TestGetTestContextBody(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "normal",
			args: args{
				getNormalContext(),
			},
			want: []byte{'1'},
		},
	}
	for _, tt := range tests {
		if got := GetTestContextBody(tt.args.c); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. GetTestContextBody() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
