// Package set is a template Set type
//
// Tries to be similar to Python's set type
package gen

// An A is the element of the set
//
// template type Set(A)

// SetNothing is used as a zero sized member in the map
type StringSetNothing struct{}

// Set provides a general purpose set modeled on Python's set type.
type StringSet struct {
	m map[string]StringSetNothing
}

// NewSizedSet returns a new empty set with the given capacity
func NewSizedStringSet(capacity int) *StringSet {
	return &StringSet{
		m: make(map[string]StringSetNothing, capacity),
	}
}

// NewSet returns a new empty set
func NewStringSet() *StringSet {
	return NewSizedStringSet(0)
}

// Len returns the number of elements in the set
func (s *StringSet) Len() int {
	return len(s.m)
}

// Contains returns whether elem is in the set or not
func (s *StringSet) Contains(elem string) bool {
	_, found := s.m[elem]
	return found
}

// Add adds elem to the set, returning the set
//
// If the element already exists then it has no effect
func (s *StringSet) Add(elem string) *StringSet {
	s.m[elem] = StringSetNothing{}
	return s
}

// AddList adds a list of elems to the set
//
// If the elements already exists then it has no effect
func (s *StringSet) AddList(elems []string) *StringSet {
	for _, elem := range elems {
		s.m[elem] = StringSetNothing{}
	}
	return s
}

// Discard removes elem from the set
//
// If it wasn't in the set it does nothing
//
// It returns the set
func (s *StringSet) Discard(elem string) *StringSet {
	delete(s.m, elem)
	return s
}

// Remove removes elem from the set
//
// It returns whether the elem was in the set or not
func (s *StringSet) Remove(elem string) bool {
	_, found := s.m[elem]
	if found {
		delete(s.m, elem)
	}
	return found
}

// Pop removes elem from the set and returns it
//
// It also returns whether the elem was found or not
func (s *StringSet) Pop(elem string) (string, bool) {
	_, found := s.m[elem]
	if found {
		delete(s.m, elem)
	}
	return elem, found
}

// AsList returns all the elements as a slice
func (s *StringSet) AsList() []string {
	elems := make([]string, len(s.m))
	i := 0
	for elem := range s.m {
		elems[i] = elem
		i++
	}
	return elems
}

// Clear removes all the elements
func (s *StringSet) Clear() *StringSet {
	s.m = make(map[string]StringSetNothing)
	return s
}

// Copy returns a shallow copy of the Set
func (s *StringSet) Copy() *StringSet {
	newSet := NewSizedStringSet(len(s.m))
	for elem := range s.m {
		newSet.m[elem] = StringSetNothing{}
	}
	return newSet
}

// Difference returns a new set with all the elements that are in this
// set but not in the other
func (s *StringSet) Difference(other *StringSet) *StringSet {
	newSet := NewSizedStringSet(len(s.m))
	for elem := range s.m {
		if _, found := other.m[elem]; !found {
			newSet.m[elem] = StringSetNothing{}
		}
	}
	return newSet
}

// DifferenceUpdate removes all the elements that are in the other set
// from this set.  It returns the set.
func (s *StringSet) DifferenceUpdate(other *StringSet) *StringSet {
	m := s.m
	for elem := range other.m {
		delete(m, elem)
	}
	return s
}

// Intersection returns a new set with all the elements that are only in this
// set and the other set. It returns the new set.
func (s *StringSet) Intersection(other *StringSet) *StringSet {
	newSet := NewSizedStringSet(len(s.m) + len(other.m))
	for elem := range s.m {
		if _, found := other.m[elem]; found {
			newSet.m[elem] = StringSetNothing{}
		}
	}
	for elem := range other.m {
		if _, found := s.m[elem]; found {
			newSet.m[elem] = StringSetNothing{}
		}
	}
	return newSet
}

// IntersectionUpdate changes this set so that it only contains
// elements that are in both this set and the other set.  It returns
// the set.
func (s *StringSet) IntersectionUpdate(other *StringSet) *StringSet {
	for elem := range s.m {
		if _, found := other.m[elem]; !found {
			delete(s.m, elem)
		}
	}
	return s
}

// Union returns a new set with all the elements that are in either
// set. It returns the new set.
func (s *StringSet) Union(other *StringSet) *StringSet {
	newSet := NewSizedStringSet(len(s.m) + len(other.m))
	for elem := range s.m {
		newSet.m[elem] = StringSetNothing{}
	}
	for elem := range other.m {
		newSet.m[elem] = StringSetNothing{}
	}
	return newSet
}

// Update adds all the elements from the other set to this set.
// It returns the set.
func (s *StringSet) Update(other *StringSet) *StringSet {
	for elem := range other.m {
		s.m[elem] = StringSetNothing{}
	}
	return s
}

// IsSuperset returns a bool indicating whether this set is a superset of other set.
func (s *StringSet) IsSuperset(strict bool, other *StringSet) bool {
	if strict && len(other.m) >= len(s.m) {
		return false
	}
string:
	for v := range other.m {
		for i := range s.m {
			if v == i {
				continue string
			}
		}
		return false
	}
	return true
}

// IsSubset returns a bool indicating whether this set is a subset of other set.
func (s *StringSet) IsSubset(strict bool, other *StringSet) bool {
	if strict && len(s.m) >= len(other.m) {
		return false
	}
string:
	for v := range s.m {
		for i := range other.m {
			if v == i {
				continue string
			}
		}
		return false
	}
	return true
}

// IsDisjoint returns a bool indicating whether this set and other set have no elements in common.
func (s *StringSet) IsDisjoint(other *StringSet) bool {
	for v := range s.m {
		if other.Contains(v) {
			return false
		}
	}
	return true
}

// SymmetricDifference returns a new set of all elements that are a member of exactly
// one of this set and other set(elements which are in one of the sets, but not in both).
func (s *StringSet) SymmetricDifference(other *StringSet) *StringSet {
	work1 := s.Union(other)
	work2 := s.Intersection(other)
	for v := range work2.m {
		delete(work1.m, v)
	}
	return work1
}

// SymmetricDifferenceUpdate modifies this set to be a set of all elements that are a member
// of exactly one of this set and other set(elements which are in one of the sets,
// but not in both) and returns this set.
func (s *StringSet) SymmetricDifferenceUpdate(other *StringSet) *StringSet {
	work := s.SymmetricDifference(other)
	*s = *work
	return s
}
