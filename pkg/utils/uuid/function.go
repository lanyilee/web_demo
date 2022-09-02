package uuid

import "strings"

// NewUuidString new V4 a uuid string without '-'
func NewUuidString() string {
	id := Must(NewV4())
	return strings.Replace(id.String(), "-", "", -1)
}

// NewUuidV4String  new V4 a uuid string with '-'
func NewUuidV4String() string {
	id := Must(NewV4())
	return id.String()
}
