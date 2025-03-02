package services

import "testing"

func TestGetPacks(t *testing.T) {

	tests := []struct {
		name        string
		orderAmount int
		packSizes   []int
		want        map[int]int
	}{
		{"no packs", 0, []int{1, 2, 3}, map[int]int{}},
		{"250 packs", 201, []int{250, 500, 1000, 2000, 5000}, map[int]int{250: 1}},
		{"251 packs", 251, []int{250, 500, 1000, 2000, 5000}, map[int]int{500: 1}},
		{"501 packs", 501, []int{250, 500, 1000, 2000, 5000}, map[int]int{500: 1, 250: 1}},
		{"12001 packs", 12001, []int{250, 500, 1000, 2000, 5000}, map[int]int{5000: 2, 2000: 1, 250: 1}},
		{"500000 packs", 500000, []int{23, 31, 53}, map[int]int{23: 2, 31: 7, 53: 9429}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetPacks(tt.orderAmount, tt.packSizes)
			if len(got) != len(tt.want) {
				t.Errorf("GetPacks() = %v, want %v", got, tt.want)
			}
			for k, v := range got {
				if tt.want[k] != v {
					t.Errorf("GetPacks() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
