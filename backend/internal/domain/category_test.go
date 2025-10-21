package domain

import "testing"

func TestCategoryType_String(t *testing.T) {
	tests := []struct {
		name string
		ct   CategoryType
		want string
	}{
		{
			name: "income",
			ct:   CategoryTypeIncome,
			want: "income",
		},
		{
			name: "outgoing",
			ct:   CategoryTypeOutgoing,
			want: "outgoing",
		},
		{
			name: "investing",
			ct:   CategoryTypeInvesting,
			want: "investing",
		},
		{
			name: "unknown",
			ct:   CategoryType(999),
			want: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ct.String(); got != tt.want {
				t.Errorf("CategoryType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
