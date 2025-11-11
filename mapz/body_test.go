package mapz

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/welllog/golib/testz"
)

func TestBody_QueryString(t *testing.T) {
	req := Body{
		"name":     "bob",
		"age":      21,
		"addr":     "wall street",
		"favorite": "football",
	}
	testz.Equal(t, "addr=wall street&age=21&favorite=football&name=bob", req.QueryString(nil))
}

func TestRequest_Read(t *testing.T) {
	req := Body{
		"name": "bob",
		"age":  21,
	}
	b, err := io.ReadAll(req)
	if err != nil {
		t.Fatal(err)
	}
	req.ClearCache()
	b1, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	testz.Equal(t, b1, b)
}
