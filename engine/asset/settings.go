package asset

import (
	"encoding/json"
	"os"
)

type Settings struct {
	settings any
	loader   *Loader
	key      string
}

func NewSettings(l *Loader, k string) *Settings {
	return &Settings{
		loader: l,
		key:    k,
	}
}

func (s *Settings) Load(sc any) {
	err := json.Unmarshal(s.loader.GetRaw(s.key), sc)
	if err != nil {
		panic(err)
	}

	s.settings = sc
}

func (s *Settings) Persist() {
	path := s.loader.GetPath(s.key)

	raw, err := json.Marshal(s.settings)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(path, raw, 0600)
	if err != nil {
		panic(err)
	}
}
