package hash

import (
	"crypto/sha256"
	"hash"
	"testing"
)

func TestCrateHash(t *testing.T) {
	type args struct {
		secretKey string
		data      []byte
		hashNew   func() hash.Hash
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				secretKey: "key",
				data:      []byte{},
				hashNew:   sha256.New,
			},
			want: "5d5d139563c95b5967b9bd9a8c9b233a9dedb45072794cd232dc1b74832607d0",
		},
		{
			args: args{
				secretKey: "Key",
				data:      []byte{},
				hashNew:   sha256.New,
			},
			want: "09f6d7714f332549c6f00817a1cfa5d524bc5ed437ffe74dbf8775a3e3b3ff12",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CrateHash(tt.args.secretKey, tt.args.data, tt.args.hashNew); got != tt.want {
				t.Errorf("CrateHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
