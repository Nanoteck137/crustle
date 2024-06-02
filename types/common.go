package types

import "path"

type WorkDir string

func (d WorkDir) String() string {
	return string(d)
}

func (d WorkDir) DataFile() string {
	return path.Join(d.String(), "data.json")
}

type DataFile struct {
	Token string `json:"token"`
}
