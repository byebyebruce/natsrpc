package natsrpc

import "testing"

func Benchmark_CombineSubsets(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = joinSubject("a", "b", "c", "d", "1", "2", "3", "4")
	}
}

func TestJoinSubject(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{[]string{"a", "b", "c"}}, "a.b.c"},
		{"", args{[]string{"", "b", "c"}}, "b.c"},
		{"", args{[]string{"", "", "b", "c"}}, "b.c"},
		{"", args{[]string{"a", "", "c"}}, "a.c"},
		{"", args{[]string{"a", "b", ""}}, "a.b"},
		{"", args{[]string{"", "a", "", "b"}}, "a.b"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := joinSubject(tt.args.s...); got != tt.want {
				t.Errorf("joinSubject() = %v, want %v", got, tt.want)
			}
		})
	}
}
