package nursetags

type Set map[PostKey]struct{}

func NewSet(values ...PostKey) Set {
	set := make(Set)
	for _, value := range values {
		set[value] = struct{}{}
	}
	return set
}

func NewSetFromSlice(values []PostKey) Set {
	return NewSet(values...)
}

func (set Set) Iter() <-chan PostKey {
	ch := make(chan PostKey)
	go func() {
		for elem := range set {
			ch <- elem
		}
		close(ch)
	}()

	return ch
}

func (set Set) ToSlice() []PostKey {
	keys := make([]PostKey, 0, len(set))
	for elem := range set {
		keys = append(keys, elem)
	}

	return keys
}

func (set Set) Intersect(other Set) Set {
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

func (set Set) Union(other Set) Set {
	union := make(Set)

	for elem := range set {
		union.Add(elem)
	}
	for elem := range other {
		union.Add(elem)
	}
	return union
}

func (set Set) Difference(other Set) Set {
	difference := make(Set)
	for elem := range set {
		if !other.Contains(elem) {
			difference.Add(elem)
		}
	}
	return difference
}

func (set Set) SymmetricDifference(other Set) Set {
	aDiff := set.Difference(other)
	bDiff := other.Difference(set)
	return aDiff.Union(bDiff)
}

func (set Set) Add(i PostKey) bool {
	_, found := set[i]
	if found {
		return false
	}

	set[i] = struct{}{}
	return true
}

func (set Set) Remove(i PostKey) {
	delete(set, i)
}

func (set Set) Contains(i ...PostKey) bool {
	for _, val := range i {
		if _, ok := set[val]; !ok {
			return false
		}
	}
	return true
}

func (set Set) Equal(other Set) bool {
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

func (set Set) IsSubset(other Set) bool {
	for elem := range set {
		if !other.Contains(elem) {
			return false
		}
	}
	return true
}

func (set Set) IsProperSubset(other Set) bool {
	return set.IsSubset(other) && !set.Equal(other)
}

func (set Set) IsSuperset(other Set) bool {
	return other.IsSubset(set)
}

func (set Set) IsProperSuperset(other Set) bool {
	return set.IsSuperset(other) && !set.Equal(other)
}
