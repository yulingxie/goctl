package sshx

import (
	"fmt"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				opts: []Option{
					Host("172.16.1.155"),
					Port(65022),
					User("xiaoyl"),
					Auth("xiaoyl"),
					AuthType(AuthTypePass),
					TimeOut(10 * time.Second),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else {
				defer got.Close()
				s, _ := got.NewSession()
				defer s.Close()
				out, _ := s.Output("pwd")
				fmt.Printf("pwd out: %v", string(out[:]))
			}
		})
	}
}
