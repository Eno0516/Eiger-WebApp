// Package port は、usecase が必要とする、時計と ID 採番の抽象。
package port

import "time"

type Clock interface {
	Now() time.Time
}

type IDGenerator interface {
	NewID() string
}
