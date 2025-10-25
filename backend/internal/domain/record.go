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

// CategoryYearSummary はカテゴリ別年次サマリーを表す
type CategoryYearSummary struct {
	CategoryID   int
	CategoryName string
	CategoryType CategoryType
	Count        int      // 該当カテゴリの取引回数
	Price        [12]int  // 12ヶ月分の金額（4月〜3月の順）
	Total        int      // 合計金額
}
