package domain

// CategoryType はカテゴリの種類を表す
type CategoryType int

const (
	// CategoryTypeIncome は収入を表す
	CategoryTypeIncome CategoryType = 1
	// CategoryTypeOutgoing は支出を表す
	CategoryTypeOutgoing CategoryType = 2
	// CategoryTypeSaving は貯金を表す
	CategoryTypeSaving CategoryType = 3
	// CategoryTypeInvesting は投資を表す
	CategoryTypeInvesting CategoryType = 4
)

// String はCategoryTypeを文字列に変換する
func (ct CategoryType) String() string {
	switch ct {
	case CategoryTypeIncome:
		return "income"
	case CategoryTypeOutgoing:
		return "outgoing"
	case CategoryTypeSaving:
		return "saving"
	case CategoryTypeInvesting:
		return "investing"
	default:
		return "unknown"
	}
}

// CategoryTypeLookup は文字列をCategoryTypeに変換するマップ
var CategoryTypeLookup = map[string]CategoryType{
	"income":    CategoryTypeIncome,
	"outgoing":  CategoryTypeOutgoing,
	"saving":    CategoryTypeSaving,
	"investing": CategoryTypeInvesting,
}

// Category はカテゴリを表すドメインエンティティ
type Category struct {
	ID           int
	CategoryID   int
	Name         string
	CategoryType CategoryType
}
