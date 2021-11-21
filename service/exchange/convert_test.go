package exchange

import (
	"reflect"
	"testing"
)

func Test_ConvertExchange(t *testing.T) {
	type args struct {
		oldMoney []float64
	}
	tests := []struct {
		name         string
		args         args
		wantNewMoney []float64
		wantErr      bool
		err          error
	}{
		{
			name: "USD to TWD",
			args: args{
				oldMoney: []float64{1, 10, 50},
			},
			wantNewMoney: []float64{27.7995, 277.995, 1389.975},
			wantErr:      false,
			err:          nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNewMoney, err := ConvertExchange(tt.args.oldMoney)
			if err != nil && tt.wantErr {
				if err.Error() != tt.err.Error() {
					t.Errorf("Quote() error = %v, wantErr %v", err.Error(), tt.err.Error())
					return
				}
			}
			if !reflect.DeepEqual(gotNewMoney, tt.wantNewMoney) {
				t.Errorf("ConvertExchange() = %v, want %v", gotNewMoney, tt.wantNewMoney)
			}
		})
	}
}
