// Mango is a modular web-application framework for Go, inspired by Rack and PEP333.
package mango

import (
	"fmt"
)

func Version() []int {
	return []int{1, 0, 0}
}

func VersionString() string {
	v := Version()
	return fmt.Sprintf("%d.%02d.%02d", v[0], v[1], v[2])
}
