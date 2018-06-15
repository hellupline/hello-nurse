package nursetags

import (
	"testing"
)

func TestSetAdd(t *testing.T) {
	set := NewMapSet()
	if new := set.Add("1"); !new {
		t.Error("Add returned false, expected true")
	}
	if exists := set.Contains("1"); !exists {
		t.Error("Set does not contains '1' key")
	}
	if len(set) != 1 {
		t.Error("Set sould have only 1 key")
	}
}

func TestSetIntersect(t *testing.T) {
	setA := NewMapSet("1", "2")
	setB := NewMapSet("2", "3")
	intersection := setA.Intersect(setB)
	expect := NewMapSet("2")

	if ok := intersection.Equal(expect); !ok {
		t.Error("Intersect failed to produce a Intersection set")
	}
}

func TestSetUnion(t *testing.T) {
	setA := NewMapSet("1", "2")
	setB := NewMapSet("2", "3")
	union := setA.Union(setB)
	expect := NewMapSet("1", "2", "3")

	if ok := union.Equal(expect); !ok {
		t.Error("Union failed to produce a Union set")
	}
}

func TestSetDifference(t *testing.T) {
	setA := NewMapSet("1", "2")
	setB := NewMapSet("2", "3")
	difference := setA.Difference(setB)
	expect := NewMapSet("1")

	if ok := difference.Equal(expect); !ok {
		t.Error("Difference failed to produce a Diff set")
	}
}

func TestSetSymmetricDifference(t *testing.T) {
	setA := NewMapSet("1", "2", "3")
	setB := NewMapSet("3", "4", "5")
	difference := setA.SymmetricDifference(setB)
	expect := NewMapSet("1", "2", "4", "5")

	if ok := difference.Equal(expect); !ok {
		t.Errorf("Difference failed to produce a Symmetric Diff set: %#+v", difference)
	}
}

func TestSetRemove(t *testing.T) {
	set := NewMapSet("1")
	set.Remove("1")
	if exists := set.Contains("1"); exists {
		t.Error("Set does contains '1' key")
	}
	if len(set) != 0 {
		t.Error("Set sould have only 0 keys")
	}
}

func TestSetContains(t *testing.T) {
	set := NewMapSet("1")
	if exists := set.Contains("1"); !exists {
		t.Error("Set does not contains '1' key")
	}
}

func TestSetEqual(t *testing.T) {
	setA := NewMapSet("1")
	setB := NewMapSet("1")
	if equal := setA.Equal(setB); !equal {
		t.Error("Equal failed to compare 2 equal sets")
	}
}

func TestSetIsSubset(t *testing.T) {
	setA := NewMapSet("1", "2")
	setB := NewMapSet("1", "2")
	if is := setB.IsSubset(setA); !is {
		t.Error("IsSubset failed to compare 2 equal sets")
	}
}

func TestSetIsProperSubset(t *testing.T) {
	setA := NewMapSet("1", "2", "3")
	setB := NewMapSet("1", "2")
	setC := NewMapSet("1", "2")
	if is := setC.IsProperSubset(setB); is {
		t.Error("IsProperSubset failed to compare 2 equal sets")
	}
	if is := setC.IsProperSubset(setA); !is {
		t.Error("IsProperSubset failed to compare 2 sets")
	}
}

func TestSetIsSuperset(t *testing.T) {
	setA := NewMapSet("1", "2")
	setB := NewMapSet("1", "2")
	if is := setA.IsSuperset(setB); !is {
		t.Error("IsSuperset failed to compare 2 equal sets")
	}
}

func TestSetIsProperSuperset(t *testing.T) {
	setA := NewMapSet("1", "2", "3")
	setB := NewMapSet("1", "2", "3")
	setC := NewMapSet("1", "2")
	if is := setA.IsProperSuperset(setB); is {
		t.Error("IsProperSuperset failed to compare 2 equal sets")
	}
	if is := setA.IsProperSuperset(setC); !is {
		t.Error("IsProperSuperset failed to compare 2 sets")
	}
}
