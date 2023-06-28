package hash

import (
	"crypto/sha256"
	"hash"
	"testing"
)

func TestCheckHash(t *testing.T) {
	type args struct {
		data      []byte
		secretKey string
		checksum  string
		hashNew   func() hash.Hash
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{
				secretKey: "key",
				data:      []byte{},
				checksum:  "5d5d139563c95b5967b9bd9a8c9b233a9dedb45072794cd232dc1b74832607d0",
				hashNew:   sha256.New,
			},
			wantErr: false,
		},
		{
			args: args{
				secretKey: "Key",
				data:      []byte{},
				checksum:  "5d5d139563c95b5967b9bd9a8c9b233a9dedb45072794cd232dc1b74832607d0",
				hashNew:   sha256.New,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckHash(tt.args.data, tt.args.secretKey, tt.args.checksum, tt.args.hashNew); (err != nil) != tt.wantErr {
				t.Errorf("CheckHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
