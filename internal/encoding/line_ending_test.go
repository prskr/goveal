package encoding_test

import (
	"io"
	"strings"
	"testing"

	"github.com/baez90/goveal/internal/encoding"
)

func TestDetect(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    encoding.LineEnding
		wantErr bool
	}{
		{
			name: "Empty file expect unknown",
			args: args{
				reader: strings.NewReader(""),
			},
			want:    encoding.LineEndingUnknown,
			wantErr: true,
		},
		{
			name: "File with only Unix line ending",
			args: args{
				reader: strings.NewReader("\n"),
			},
			want:    encoding.LineEndingUnix,
			wantErr: false,
		},
		{
			name: "File with only Windows line ending",
			args: args{
				reader: strings.NewReader("\r\n"),
			},
			want:    encoding.LineEndingWindows,
			wantErr: false,
		},
		{
			name: "File with multiple lines - Unix file ending",
			args: args{
				reader: strings.NewReader("Hello, World\nThis comes from Unix!\n"),
			},
			want:    encoding.LineEndingUnix,
			wantErr: false,
		},
		{
			name: "File with multiple lines - Windows file ending",
			args: args{
				reader: strings.NewReader("Hello, World\r\nThis comes from Windows!\r\n"),
			},
			want:    encoding.LineEndingWindows,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encoding.Detect(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("Detect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Detect() got = %v, want %v", got, tt.want)
			}
		})
	}
}
