package mapper

import "strings"

type SlashMapper struct {
}

func NewSlashMapper() *SlashMapper {
	return &SlashMapper{}
}

func (m *SlashMapper) Encode(name string) string {
	return strings.ReplaceAll(name, "/", "_")
	
}

func (m *SlashMapper) Decode(hash string) (string, error) {
	return strings.ReplaceAll(hash, "_", "/"), nil
}