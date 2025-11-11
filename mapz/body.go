package mapz

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/welllog/golib/strz"
)

const _HIDDEN_KEY = "xx---.internal.request.payload.---xx"

// Body is a map[string]any
type Body map[string]any

// Read reads the body as bytes.Reader
func (b Body) Read(p []byte) (int, error) {
	r, err := b.JsonReader()
	if err != nil {
		return 0, err
	}

	return r.Read(p)
}

// CleanPayload cleans the payload
// Deprecated: use ClearCache instead.
func (b Body) CleanPayload() {
	delete(b, _HIDDEN_KEY)
}

// ClearCache clears the cached reader
func (b Body) ClearCache() {
	delete(b, _HIDDEN_KEY)
}

func (b Body) Close() error {
	b.ClearCache()
	return nil
}

func (b Body) Seek(offset int64, whence int) (int64, error) {
	r, err := b.JsonReader()
	if err != nil {
		return 0, err
	}

	return r.Seek(offset, whence)
}

func (b Body) CloneJsonReader() (*bytes.Reader, error) {
	r, err := b.JsonReader()
	if err != nil {
		return nil, err
	}

	copied := *r
	_, _ = copied.Seek(0, io.SeekStart)
	return &copied, nil
}

func (b Body) MustJsonReader() *bytes.Reader {
	r, err := b.JsonReader()
	if err != nil {
		panic(err)
	}
	return r
}

// JsonReader will return a bytes.Reader for the JSON representation of the Body
// It will cache the reader for subsequent calls
func (b Body) JsonReader() (*bytes.Reader, error) {
	val, ok := b[_HIDDEN_KEY]
	if !ok {
		bs, err := json.Marshal(b)
		if err != nil {
			return nil, err
		}

		reader := bytes.NewReader(bs)
		b[_HIDDEN_KEY] = reader
		return reader, nil
	}

	r, ok := val.(*bytes.Reader)
	if !ok {
		return nil, fmt.Errorf("invalid reader type")
	}
	return r, nil
}

// QueryString returns the query string
func (b Body) QueryString(valueEncode func(string) string) string {
	bs := b.QueryBytes(valueEncode)
	return strz.UnsafeString(bs)
}

// QueryBytes returns the query bytes
func (b Body) QueryBytes(valueEncode func(string) string) []byte {
	b.ClearCache()

	if len(b) == 0 {
		return nil
	}

	keys := make([]string, 0, len(b))
	var initSize int
	for k := range b {
		keys = append(keys, k)
		initSize += 2 * len(k)
	}

	sort.Strings(keys)

	buf := bytes.NewBuffer(make([]byte, 0, initSize))
	for _, k := range keys {
		buf.WriteByte('&')
		buf.WriteString(k)
		buf.WriteByte('=')
		value := toStr(b[k])
		if valueEncode != nil {
			value = valueEncode(value)
		}
		buf.WriteString(value)
	}
	_, _ = buf.ReadByte()
	return buf.Bytes()
}

func toStr(value any) string {
	switch v := value.(type) {
	case []byte:
		return strz.UnsafeString(v)
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
		b, _ := json.Marshal(value)
		return strz.UnsafeString(b)
	}
}

var _popEncodeReplacer = strings.NewReplacer("+", "%20", "*", "%2A", "%7E", "~")

// PopEncode encodes the string for pop
func PopEncode(str string) string {
	return _popEncodeReplacer.Replace(url.QueryEscape(str))
}
