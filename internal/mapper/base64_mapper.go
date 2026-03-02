package mapper

import "encoding/base64"

type Base64Mapper struct{}

func NewBase64Mapper() *Base64Mapper {
	return &Base64Mapper{}
}

func (m *Base64Mapper) Encode(name string) string {
	return base64.URLEncoding.EncodeToString([]byte(name))
}

func (m *Base64Mapper) Decode(key string) (string, error) {
	b, err := base64.URLEncoding.DecodeString(key)
	return string(b), err
}
