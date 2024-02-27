package rolling

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBucket_Avg(t *testing.T) {
	tests := []struct {
		name   string
		bucket *Bucket
		want   float64
	}{
		{
			name:   "Empty bucket",
			bucket: &Bucket{},
			want:   0,
		},
		{
			name: "Non-empty bucket",
			bucket: &Bucket{
				count: 3,
				sum:   9,
			},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := tt.bucket
			got := b.Avg()
			assert.Equal(t, tt.want, got, "Avg() returned unexpected value")
		})
	}
}

func TestBucket_Count(t *testing.T) {
	tests := []struct {
		name   string
		bucket *Bucket
		want   uint64
	}{
		{
			name:   "Empty bucket",
			bucket: &Bucket{},
			want:   0,
		},
		{
			name: "Non-empty bucket",
			bucket: &Bucket{
				count: 3,
				sum:   9,
			},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := tt.bucket
			got := b.Count()
			assert.Equal(t, tt.want, got, "Count() returned unexpected value")
		})
	}
}

func TestBucket_Sum(t *testing.T) {
	tests := []struct {
		name   string
		bucket *Bucket
		want   float64
	}{
		{
			name:   "Empty bucket",
			bucket: &Bucket{},
			want:   0,
		},
		{
			name: "Non-empty bucket",
			bucket: &Bucket{
				count: 3,
				sum:   9,
			},
			want: 9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := tt.bucket
			got := b.Sum()
			assert.Equal(t, tt.want, got, "Sum() returned unexpected value")
		})
	}
}
