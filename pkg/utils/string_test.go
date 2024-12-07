package utils

import (
	"testing"
)

func TestConvertTraditionalToSimplified(t *testing.T) {
	type args struct {
		input string
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
				input: `（堅韌商戶女VS京城貴公子、腹黑心機、醋王、情敵面前狠女主面前嬌弱的綠茶大師）`,
			},
			want:    `（坚韧商户女VS京城贵公子、腹黑心机、醋王、情敌面前狠女主面前娇弱的绿茶大师）`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertTraditionalToSimplified(tt.args.input)
			if got != tt.want {
				t.Errorf("ConvertTraditionalToSimplified() got = %v, want %v", got, tt.want)
			}
		})
	}
}
