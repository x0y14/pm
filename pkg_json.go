package pm

import "encoding/json"

type PkgJson struct {
	Deps []*Dependencies `json:"deps,omitempty"`
}

type Dependencies struct {
	Url     string `json:"url"`
	Version string `json:"version"`
}

func UnmarshalPkgJson(bytes []byte) (*PkgJson, error) {
	var pkgJson PkgJson
	err := json.Unmarshal(bytes, &pkgJson)
	if err != nil {
		return nil, err
	}
	return &pkgJson, nil
}
