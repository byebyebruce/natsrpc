package natsrpc

import "testing"

func Benchmark_CombineSubsets(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = JoinSubject("a", "b", "c")
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JoinSubject(tt.args.s...); got != tt.want {
				t.Errorf("JoinSubject() = %v, want %v", got, tt.want)
			}
		})
	}
}
