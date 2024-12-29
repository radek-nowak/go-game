package vector

import (
	"math"
	"reflect"
	"testing"
)

func floatsAlmostEqual(f1, f2, epsilon float64) bool {
	return math.Abs(f1-f2) <= epsilon
}

func TestNewVector(t *testing.T) {
	type args struct {
		x float64
		y float64
	}
	tests := []struct {
		name string
		args args
		want *Vector
	}{
		{
			name: "Zero vector",
			args: args{x: 0, y: 0},
			want: &Vector{x: 0, y: 0},
		},
		{
			name: "Positive values",
			args: args{x: 1.5, y: 2.5},
			want: &Vector{x: 1.5, y: 2.5},
		},
		{
			name: "Negative values",
			args: args{x: -3, y: -4},
			want: &Vector{x: -3, y: -4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewVector(tt.args.x, tt.args.y); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewVector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector_Add(t *testing.T) {
	type fields struct {
		x float64
		y float64
	}
	type args struct {
		dx float64
		dy float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Vector
	}{
		{
			name:   "Add positive values",
			fields: fields{x: 1, y: 1},
			args:   args{dx: 2, dy: 3},
			want:   &Vector{x: 3, y: 4},
		},
		{
			name:   "Add zero values",
			fields: fields{x: 1, y: 1},
			args:   args{dx: 0, dy: 0},
			want:   &Vector{x: 1, y: 1},
		},
		{
			name:   "Add negative values",
			fields: fields{x: 5, y: 5},
			args:   args{dx: -2, dy: -3},
			want:   &Vector{x: 3, y: 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Vector{
				x: tt.fields.x,
				y: tt.fields.y,
			}
			v.Add(tt.args.dx, tt.args.dy)
			if !floatsAlmostEqual(v.x, tt.want.x, 1e-9) || !floatsAlmostEqual(v.y, tt.want.y, 1e-9) {
				t.Errorf("Vector.Add() = %v, want %v", v, tt.want)
			}
		})
	}
}

func TestVector_Rotate(t *testing.T) {
	type fields struct {
		x float64
		y float64
	}
	type args struct {
		angle float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Vector
	}{
		{
			name:   "Rotate 180 degrees",
			fields: fields{x: 1.0, y: 1.0},
			args:   args{angle: 180},
			want:   &Vector{x: -1.0, y: -1.0},
		},
		{
			name:   "Rotate 90 degrees",
			fields: fields{x: 1.0, y: 0.0},
			args:   args{angle: 90},
			want:   &Vector{x: 0.0, y: 1.0},
		},
		{
			name:   "Rotate 360 degrees",
			fields: fields{x: 1.0, y: 1.0},
			args:   args{angle: 360},
			want:   &Vector{x: 1.0, y: 1.0},
		},
		{
			name:   "Rotate with negative angle",
			fields: fields{x: 1.0, y: 0.0},
			args:   args{angle: -90},
			want:   &Vector{x: 0.0, y: -1.0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Vector{
				x: tt.fields.x,
				y: tt.fields.y,
			}
			got := v.Rotate(tt.args.angle)
			if !floatsAlmostEqual(got.x, tt.want.x, 1e-9) || !floatsAlmostEqual(got.y, tt.want.y, 1e-9) {
				t.Errorf("Vector.Rotate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector_Normalize(t *testing.T) {
	type fields struct {
		x float64
		y float64
	}
	tests := []struct {
		name   string
		fields fields
		want   *Vector
	}{
		{
			name:   "Normalize non-zero vector",
			fields: fields{x: 3, y: 4},
			want:   &Vector{x: 0.6, y: 0.8},
		},
		{
			name:   "Normalize unit vector",
			fields: fields{x: 1, y: 0},
			want:   &Vector{x: 1, y: 0},
		},
		{
			name:   "Normalize zero vector",
			fields: fields{x: 0, y: 0},
			want:   &Vector{x: 0, y: 0},
		},
		{
			name:   "Normalize tiny vector",
			fields: fields{x: float64(1e-10), y: float64(1e-10)},
			want:   &Vector{x: 0, y: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Vector{
				x: tt.fields.x,
				y: tt.fields.y,
			}
			got := v.Normalize()
			if !floatsAlmostEqual(got.x, tt.want.x, 1e-9) || !floatsAlmostEqual(got.y, tt.want.y, 1e-9) {
				t.Errorf("Vector.Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}

