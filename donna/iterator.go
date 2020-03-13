package donna

import "os"

type paramIterator struct {
	currIdx int
}

// Resets the iterator
func (pi *paramIterator) Reset() {
	pi.currIdx = 0
}

// Rewinds the iterator one step.
func (pi *paramIterator) Rewind() {
	if pi.currIdx > 0 {
		pi.currIdx--
	}
}

// Returns the current parameter, and a flag indicating if this was
// a valid request.
func (pi *paramIterator) Curr() string {
	return os.Args[pi.currIdx]
}

// Returns the next parameter, and a flag indicating if there was
// an argument to return.
func (pi *paramIterator) Next() (string, bool) {
	if pi.currIdx + 1 == len(os.Args) {
		return "", false
	}
	pi.currIdx++
	return os.Args[pi.currIdx], true
}
