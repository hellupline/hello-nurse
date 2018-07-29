package main

type (
	PostKeySet map[PostKey]struct{} // nolint
)

func NewPostKeySet(values ...PostKey) PostKeySet { // nolint
	set := make(PostKeySet)
	for _, value := range values {
		set[value] = struct{}{}
	}
	return set
}

func NewPostKeySetFromSlice(values []PostKey) PostKeySet { // nolint
	return NewPostKeySet(values...)
}

func (set PostKeySet) Iter() <-chan PostKey { // nolint
	ch := make(chan PostKey)
	go func() {
		for elem := range set {
			ch <- elem
		}
		close(ch)
	}()

	return ch
}

func (set PostKeySet) ToSlice() []PostKey { // nolint
	keys := make([]PostKey, 0, len(set))
	for elem := range set {
		keys = append(keys, elem)
	}

	return keys
}

func (set PostKeySet) Intersect(other PostKeySet) PostKeySet { // nolint
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

func (set PostKeySet) Union(other PostKeySet) PostKeySet { // nolint
	union := make(PostKeySet)

	for elem := range set {
		union.Add(elem)
	}
	for elem := range other {
		union.Add(elem)
	}
	return union
}

func (set PostKeySet) Difference(other PostKeySet) PostKeySet { // nolint
	difference := make(PostKeySet)
	for elem := range set {
		if !other.Contains(elem) {
			difference.Add(elem)
		}
	}
	return difference
}

func (set PostKeySet) SymmetricDifference(other PostKeySet) PostKeySet { // nolint
	aDiff := set.Difference(other)
	bDiff := other.Difference(set)
	return aDiff.Union(bDiff)
}

func (set PostKeySet) Add(i PostKey) bool { // nolint
	_, found := set[i]
	if found {
		return false
	}

	set[i] = struct{}{}
	return true
}

func (set PostKeySet) Remove(i PostKey) { // nolint
	delete(set, i)
}

func (set PostKeySet) Contains(i ...PostKey) bool { // nolint
	for _, val := range i {
		if _, ok := set[val]; !ok {
			return false
		}
	}
	return true
}

func (set PostKeySet) Equal(other PostKeySet) bool { // nolint
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

func (set PostKeySet) IsSubset(other PostKeySet) bool { // nolint
	for elem := range set {
		if !other.Contains(elem) {
			return false
		}
	}
	return true
}

func (set PostKeySet) IsProperSubset(other PostKeySet) bool { // nolint
	return set.IsSubset(other) && !set.Equal(other)
}

func (set PostKeySet) IsSuperset(other PostKeySet) bool { // nolint
	return other.IsSubset(set)
}

func (set PostKeySet) IsProperSuperset(other PostKeySet) bool { // nolint
	return set.IsSuperset(other) && !set.Equal(other)
}
