package dao

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPaginate(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		wantOffset int
		wantLimit  int
		wantOrder  string
		wantErr    string
	}{
		{"No pagination", "", -1, -1, "", "<nil>"},
		{"Third page", "page=3&size=5&order=name DESC", 15, 5, "name DESC", "<nil>"},
		{"Limited", "size=5", -1, 5, "", "<nil>"},
		{"Bad page params", "page=3&size=-1", -1, -1, "", "page size cannot be lower than 0 for nonnegative page number"},
		{"Bad page size", "size=0", -1, -1, "", "page size cannot be 0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := gin.Context{Request: &http.Request{URL: &url.URL{RawQuery: tt.query}}}
			gotOffset, gotLimit, gotOrder, err := Paginate(&ctx)
			if fmt.Sprint(err) != tt.wantErr {
				t.Errorf("error\nactual: %v\nexpect: %v", err, tt.wantErr)
				return
			}
			if gotOffset != tt.wantOffset {
				t.Errorf("offset\nactual: %v\nexpect: %v", gotOffset, tt.wantOffset)
			}
			if gotLimit != tt.wantLimit {
				t.Errorf("limit\nactual: %v\nexpect: %v", gotLimit, tt.wantLimit)
			}
			if gotOrder != tt.wantOrder {
				t.Errorf("order\nactual: %v\nexpect: %v", gotOrder, tt.wantOrder)
			}
		})
	}
}
