package utils

import "testing"

func TestStringInSlice(t *testing.T) {
	type args struct {
		a    string
		list []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "slice contains the sting",
			args: args{
				a: "test",
				list: []string{
					"www",
					"test",
					"xxx",
				},
			},
			want: true,
		},
		{
			name: "slice doesn't contain the sting",
			args: args{
				a: "test",
				list: []string{
					"www",
					"xxx",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringInSlice(tt.args.a, tt.args.list); got != tt.want {
				t.Errorf("StringInSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
