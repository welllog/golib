package mapz

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"unsafe"
)

const _HIDDEN_KEY = "xx---.internal.request.payload.---xx"

type Body map[string]any

func (r Body) Read(p []byte) (n int, err error) {
	val, ok := r[_HIDDEN_KEY]
	if !ok {
		var b []byte
		b, err = json.Marshal(r)
		if err != nil {
			return
		}
		reader := bytes.NewReader(b)
		r[_HIDDEN_KEY] = reader
		return reader.Read(p)
	}
	return val.(*bytes.Reader).Read(p)
}

func (r Body) CleanPayload() {
	delete(r, _HIDDEN_KEY)
}

func (r Body) QueryString(valueEncode func(string) string) string {
	b := r.QueryBytes(valueEncode)
	return *(*string)(unsafe.Pointer(&b))
}

func (r Body) QueryBytes(valueEncode func(string) string) []byte {
	r.CleanPayload()

	if len(r) == 0 {
		return nil
	}

	keys := make([]string, 0, len(r))
	var initSize int
	for k := range r {
		keys = append(keys, k)
		initSize += len(k) + 3
	}

	sort.Strings(keys)

	if valueEncode == nil {
		valueEncode = emptyEncode
	}

	buf := bytes.NewBuffer(make([]byte, 0, initSize))
	for _, k := range keys {
		buf.WriteByte('&')
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(valueEncode(toStr(r[k])))
	}
	_, _ = buf.ReadByte()
	return buf.Bytes()
}

func toStr(value any) string {
	switch v := value.(type) {
	case []byte:
		return *(*string)(unsafe.Pointer(&v))
	case string:
		return v
	case nil:
		return ""
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case fmt.Stringer:
		return v.String()
	default:
		b, _ := json.Marshal(v)
		return *(*string)(unsafe.Pointer(&b))
	}
}

var _popEncodeReplacer = strings.NewReplacer("+", "%20", "*", "%2A", "%7E", "~")

func PopEncode(str string) string {
	return _popEncodeReplacer.Replace(url.QueryEscape(str))
}

func emptyEncode(str string) string {
	return str
}
