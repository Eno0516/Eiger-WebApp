// Package entity は、車の予約ミニアプリの Entity とそのルールを定義する。
package entity

// Car はサークルが管理する共有車。
type Car struct {
	ID          string
	CircleID    string
	Name        string
	PlateNumber string
	// Available が false の車は、予約できない(故障中・整備中など)。
	Available bool
}
