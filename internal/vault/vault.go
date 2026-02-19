package vault

import "sort"

type Vaults struct {
	vaults map[string]string
}

func DefaultVaults() Vaults {
	return New(map[string]string{
		"develop": "~/Documents/Obsidian/develop",
		"work":    "~/Documents/Obsidian/work",
		"private": "~/Documents/Obsidian/private",
	})
}

func New(vaults map[string]string) Vaults {
	if vaults == nil {
		vaults = make(map[string]string)
	}
	return Vaults{vaults: vaults}
}

func NewVaults(vaults map[string]string) Vaults {
	return New(vaults)
}

func (v Vaults) Has(name string) bool {
	_, ok := v.vaults[name]
	return ok
}

func (v Vaults) Path(name string) (string, bool) {
	path, ok := v.vaults[name]
	return path, ok
}

func (v Vaults) Names() []string {
	names := make([]string, 0, len(v.vaults))
	for name := range v.vaults {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
