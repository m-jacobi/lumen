package vault

import "sort"

type Vaults struct {
	m map[string]string
}

func DefaultVaults() Vaults {
	return Vaults{m: map[string]string{
		"develop": "~/Documents/Obsidian/develop",
		"work":    "~/Documents/Obsidian/work",
		"private": "~/Documents/Obsidian/private",
	}}
}

func (v Vaults) Has(name string) bool {
	_, ok := v.m[name]
	return ok
}

func (v Vaults) Path(name string) (string, bool) {
	p, ok := v.m[name]
	return p, ok
}

func (v Vaults) Names() []string {
	out := make([]string, 0, len(v.m))
	for k := range v.m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
