package nursedatabase

type (
	PostKeySet map[PostKey]struct{} // nolint: golint
)

func NewPostKeySet(values ...PostKey) PostKeySet { // nolint: golint
	set := make(PostKeySet)
	for _, value := range values {
		set[value] = struct{}{}
	}
	return set
}

func NewPostKeySetFromSlice(values []PostKey) PostKeySet { // nolint: golint
	return NewPostKeySet(values...)
}

func (set PostKeySet) Iter() <-chan PostKey { // nolint: golint
	ch := make(chan PostKey)
	go func() {
		for elem := range set {
			ch <- elem
		}
		close(ch)
	}()

	return ch
}

func (set PostKeySet) ToSlice() []PostKey { // nolint: golint
	keys := make([]PostKey, 0, len(set))
	for elem := range set {
		keys = append(keys, elem)
	}

	return keys
}

func (set PostKeySet) Intersect(other PostKeySet) PostKeySet { // nolint: golint
	intersection := make(PostKeySet)
	// loop over smaller set
	if len(set) < len(other) {
		for elem := range set {
			if other.Contains(elem) {
				intersection.Add(elem)
			}
		}
	} else {
		for elem := range other {
			if set.Contains(elem) {
				intersection.Add(elem)
			}
		}
	}
	return intersection
}

func (set PostKeySet) Union(other PostKeySet) PostKeySet { // nolint: golint
	union := make(PostKeySet)

	for elem := range set {
		union.Add(elem)
	}
	for elem := range other {
		union.Add(elem)
	}
	return union
}

func (set PostKeySet) Difference(other PostKeySet) PostKeySet { // nolint: golint
	difference := make(PostKeySet)
	for elem := range set {
		if !other.Contains(elem) {
			difference.Add(elem)
		}
	}
	return difference
}

func (set PostKeySet) SymmetricDifference(other PostKeySet) PostKeySet { // nolint: golint
	aDiff := set.Difference(other)
	bDiff := other.Difference(set)
	return aDiff.Union(bDiff)
}

func (set PostKeySet) Add(i PostKey) bool { // nolint: golint
	_, found := set[i]
	if found {
		return false
	}

	set[i] = struct{}{}
	return true
}

func (set PostKeySet) Remove(i PostKey) { // nolint: golint
	delete(set, i)
}

func (set PostKeySet) Contains(i ...PostKey) bool { // nolint: golint
	for _, val := range i {
		if _, ok := set[val]; !ok {
			return false
		}
	}
	return true
}

func (set PostKeySet) Equal(other PostKeySet) bool { // nolint: golint
	if len(set) != len(other) {
		return false
	}
	for elem := range set {
		if !other.Contains(elem) {
			return false
		}
	}
	return true
}

func (set PostKeySet) IsSubset(other PostKeySet) bool { // nolint: golint
	for elem := range set {
		if !other.Contains(elem) {
			return false
		}
	}
	return true
}

func (set PostKeySet) IsProperSubset(other PostKeySet) bool { // nolint: golint
	return set.IsSubset(other) && !set.Equal(other)
}

func (set PostKeySet) IsSuperset(other PostKeySet) bool { // nolint: golint
	return other.IsSubset(set)
}

func (set PostKeySet) IsProperSuperset(other PostKeySet) bool { // nolint: golint
	return set.IsSuperset(other) && !set.Equal(other)
}
