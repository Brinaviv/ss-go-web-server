package model

import "encoding/json"

type Following struct {
	Following []UserID `json:"following"` // TODO: find out why JSON array has empty entry
}

func (f *Following) stringSlice() []string {
	res := make([]string, len(f.Following))
	for _, id := range f.Following {
		res = append(res, id.String())
	}
	return res
}

func (u *Following) MarshalJSON() ([]byte, error) {
	type Alias Following
	return json.Marshal(&struct {
		Following []string `json:"following"`
		*Alias
	}{
		Following: u.stringSlice(),
		Alias:     (*Alias)(u),
	})
}
