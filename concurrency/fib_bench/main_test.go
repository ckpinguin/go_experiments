package main

import (
	"testing"
)

func BenchmarkFibIt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fibIt(46)
	}
}

func BenchmarkFibP(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fibP(46)
	}
}

func BenchmarkFib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fib(46)
	}
}
func Test_fib(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fib(tt.args.n); got != tt.want {
				t.Errorf("fib() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fibIt(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fibIt(tt.args.n); got != tt.want {
				t.Errorf("fibIt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fibP(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fibP(tt.args.n); got != tt.want {
				t.Errorf("fibP() = %v, want %v", got, tt.want)
			}
		})
	}
}
