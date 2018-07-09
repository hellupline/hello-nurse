package nursetags

type Set map[PostKey]struct{} // nolint

func NewSet(values ...PostKey) Set { // nolint
	set := make(Set)
	for _, value := range values {
		set[value] = struct{}{}
	}
	return set
}

func NewSetFromSlice(values []PostKey) Set { // nolint
	return NewSet(values...)
}

func (set Set) Iter() <-chan PostKey { // nolint
	ch := make(chan PostKey)
	go func() {
		for elem := range set {
			ch <- elem
		}
		close(ch)
	}()

	return ch
}

func (set Set) ToSlice() []PostKey { // nolint
	keys := make([]PostKey, 0, len(set))
	for elem := range set {
		keys = append(keys, elem)
	}

	return keys
}

func (set Set) Intersect(other Set) Set { // nolint
	intersection := make(Set)
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

func (set Set) Union(other Set) Set { // nolint
	union := make(Set)

	for elem := range set {
		union.Add(elem)
	}
	for elem := range other {
		union.Add(elem)
	}
	return union
}

func (set Set) Difference(other Set) Set { // nolint
	difference := make(Set)
	for elem := range set {
		if !other.Contains(elem) {
			difference.Add(elem)
		}
	}
	return difference
}

func (set Set) SymmetricDifference(other Set) Set { // nolint
	aDiff := set.Difference(other)
	bDiff := other.Difference(set)
	return aDiff.Union(bDiff)
}

func (set Set) Add(i PostKey) bool { // nolint
	_, found := set[i]
	if found {
		return false
	}

	set[i] = struct{}{}
	return true
}

func (set Set) Remove(i PostKey) { // nolint
	delete(set, i)
}

func (set Set) Contains(i ...PostKey) bool { // nolint
	for _, val := range i {
		if _, ok := set[val]; !ok {
			return false
		}
	}
	return true
}

func (set Set) Equal(other Set) bool { // nolint
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

func (set Set) IsSubset(other Set) bool { // nolint
	for elem := range set {
		if !other.Contains(elem) {
			return false
		}
	}
	return true
}

func (set Set) IsProperSubset(other Set) bool { // nolint
	return set.IsSubset(other) && !set.Equal(other)
}

func (set Set) IsSuperset(other Set) bool { // nolint
	return other.IsSubset(set)
}

func (set Set) IsProperSuperset(other Set) bool { // nolint
	return set.IsSuperset(other) && !set.Equal(other)
}
