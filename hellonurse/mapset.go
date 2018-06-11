package hellonurse

type MapSet map[string]struct{}

func NewMapSet(values ...string) MapSet {
	set := make(MapSet)
	for _, value := range values {
		set[value] = struct{}{}
	}
	return set
}

func NewMapSetFromSlice(values []string) MapSet {
	return NewMapSet(values...)
}

func (set MapSet) Iter() <-chan string {
	ch := make(chan string)
	go func() {
		for elem := range set {
			ch <- elem
		}
		close(ch)
	}()

	return ch
}

func (set MapSet) ToSlice() []string {
	keys := make([]string, 0, len(set))
	for elem := range set {
		keys = append(keys, elem)
	}

	return keys
}

func (set MapSet) Intersect(other MapSet) MapSet {
	intersection := make(MapSet)
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

func (set MapSet) Union(other MapSet) MapSet {
	union := make(MapSet)

	for elem := range set {
		union.Add(elem)
	}
	for elem := range other {
		union.Add(elem)
	}
	return union
}

func (set MapSet) Difference(other MapSet) MapSet {
	difference := make(MapSet)
	for elem := range set {
		if !other.Contains(elem) {
			difference.Add(elem)
		}
	}
	return difference
}

func (set MapSet) SymmetricDifference(other MapSet) MapSet {
	aDiff := set.Difference(other)
	bDiff := other.Difference(set)
	return aDiff.Union(bDiff)
}

func (set MapSet) Add(i string) bool {
	_, found := set[i]
	if found {
		return false
	}

	set[i] = struct{}{}
	return true
}

func (set MapSet) Remove(i string) {
	delete(set, i)
}

func (set MapSet) Contains(i ...string) bool {
	for _, val := range i {
		if _, ok := set[val]; !ok {
			return false
		}
	}
	return true
}

func (set MapSet) Equal(other MapSet) bool {
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

func (set MapSet) IsSubset(other MapSet) bool {
	for elem := range set {
		if !other.Contains(elem) {
			return false
		}
	}
	return true
}

func (set MapSet) IsProperSubset(other MapSet) bool {
	return set.IsSubset(other) && !set.Equal(other)
}

func (set MapSet) IsSuperset(other MapSet) bool {
	return other.IsSubset(set)
}

func (set MapSet) IsProperSuperset(other MapSet) bool {
	return set.IsSuperset(other) && !set.Equal(other)
}
