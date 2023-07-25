package qr

import (
	"reflect"
	"testing"
)

func Test_mergeBlocks(t *testing.T) {
	type args struct {
		src [][]byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "1",
			args: args{src: [][]byte{{1, 2, 3}, {4, 5, 6, 7}, {8, 9, 10, 11}}},
			want: []byte{1, 4, 8, 2, 5, 9, 3, 6, 10, 7, 11},
		},
		{
			name: "2",
			args: args{src: [][]byte{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}},
			want: []byte{1, 4, 7, 2, 5, 8, 3, 6, 9},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeBlocks(tt.args.src); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mergeBlocks() = %v, want %v", got, tt.want)
			}
		})
	}
}
