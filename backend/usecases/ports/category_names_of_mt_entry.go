package ports

import (
	"github.com/aereal/coll"
	"github.com/aereal/mt"
)

func CategoryNamesOfMTEntry(entry *mt.Entry) *coll.Set[string] {
	names := coll.NewSet[string]()
	if entry.PrimaryCategory != "" {
		names.Append(entry.PrimaryCategory)
	}
	for _, c := range entry.Category {
		names.Append(c)
	}
	return names
}
