package Entity

import (
	"github.com/Enrikerf/pfm/commandExecutor/app/Domain/Entity"
	"testing"
)

func Test_controlAlgorithm_Calculate(t *testing.T) {
	type fields struct {
		goal         float64
		P            float64
		I            float64
		D            float64
		integralTerm float64
		sampleTime   float64
		currentValue float64
		currentError float64
		pastError    float64
	}
	type args struct {
		currentValue float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ca := Entity.NewControlAlgorithm()
			if got := ca.Calculate(tt.args.currentValue); got != tt.want {
				t.Errorf("Calculate() = %v, want %v", got, tt.want)
			}
		})
	}
}
