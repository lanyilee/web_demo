package errpkg

import "errors"

var ErrNoRows = errors.New("sql: no rows in result set")
