package rand

import (
	"testing"
	"time"
)

func TestSeed(t *testing.T) {
	type args struct {
		seed int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"0",
			args{
				0,
			},
		},
		{
			"100",
			args{
				100,
			},
		},
		{
			"10000",
			args{
				10000,
			},
		},
		{
			"Now",
			args{
				time.Now().Unix(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Seed(tt.args.seed)
		})
	}
}
