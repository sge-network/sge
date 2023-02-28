package types

import "gopkg.in/yaml.v2"

// String implements the Stringer interface.
func (p Phase) String() string {
	out, err := yaml.Marshal(p)
	if err != nil {
		panic(err)
	}
	return string(out)
}
