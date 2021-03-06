package schema

// QualifiedNameList implements the Sortable interface for QualifiedNames
type QualifiedNameList []QualifiedName

func (q QualifiedNameList) Len() int {
	return len(q)
}

func (q QualifiedNameList) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

// Sort qualified names by their name first, then by its namespace
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
