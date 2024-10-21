package cache

// DefaultClone attempts to make a safe copy of the given
// object using interfaces.
//
// - Clone() *T
// - Copy() *T
// - Clone() T
// - Copy() T
func DefaultClone[T any](src *T) (*T, bool) {
	if src == nil {
		return nil, false
	}

	switch v := any(src).(type) {
	case interface{ Clone() *T }:
		return v.Clone(), true
	case interface{ Copy() *T }:
		return v.Copy(), true
	case interface{ Clone() T }:
		p := v.Clone()
		return &p, true
	case interface{ Copy() T }:
		p := v.Copy()
		return &p, true
	default:
		// TODO: safe shallow copy
		return nil, false
	}
}
