package domain

import "time"

// Record は金融取引記録を表すドメインエンティティ
type Record struct {
	ID           int
	CategoryID   int
	CategoryName string
	Datetime     time.Time
	From         string
	Type         string
	Price        int
	Memo         string
}
