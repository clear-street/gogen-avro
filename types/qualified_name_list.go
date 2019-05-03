package types

// QualifiedNameList implements the Sortable interface for QualifiedNames
type QualifiedNameList []QualifiedName

func (q QualifiedNameList) Len() int {
	return len(q)
}

func (q QualifiedNameList) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

// Sort functions by the struct to which they're attached first, then the name of the method itself. If the function isn't attached to a struct, put it at the bottom
func (q QualifiedNameList) Less(i, j int) bool {
	if q[i].Name == "" && q[j].Name != "" {
		return true
	}
	if q[i].Name != "" && q[j].Name == "" {
		return false
	}
	if q[i].Name != "" && q[j].Name != "" {
		if q[i].Name > q[j].Name {
			return true
		} else if q[i].Name < q[j].Name {
			return false
		}
	}
	return q[i].Namespace < q[j].Namespace
}
