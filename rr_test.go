package main

import (
	"testing"
)

func Test_roundrobin(t *testing.T) {
	gResolv.resolv["foo"] = []string{"foo", "bar"}
	gResolv.index["foo"] = -1
	type args struct {
		host string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{name: "r1", args: args{host: "foo"}, want: "foo", want1: true},
		{name: "r2", args: args{host: "foo"}, want: "bar", want1: true},
		{name: "r3", args: args{host: "foo"}, want: "foo", want1: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := roundrobin(tt.args.host)
			if got != tt.want {
				t.Errorf("roundrobin() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("roundrobin() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
