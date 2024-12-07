package utils

import "testing"

func TestGetUrlLastSegment(t *testing.T) {
	type args struct {
		inputURL string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test_1",
			args: args{
				inputURL: "https://www.xbanxia.com/books/307150.html",
			},
			want:    "307150",
			wantErr: false,
		},
		{
			name: "test_1",
			args: args{
				inputURL: "https://www.52shuku.vip/yanqing/pm/h2nh.html",
			},
			want:    "h2nh",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUrlLastSegment(tt.args.inputURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUrlLastSegment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUrlLastSegment() got = %v, want %v", got, tt.want)
			}
		})
	}
}
