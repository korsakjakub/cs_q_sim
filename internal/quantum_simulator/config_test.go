package quantum_simulator

import (
	"os"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	type args struct {
		filename      []string
		lines          []string
	}
	tests := []struct {
		name string
		args args
		want Config
	}{
		{
			name: "single file",
			args: args{
				filename: []string{"/tmp/c1.yaml"},
				lines: []string{
					"physics:",
					"  spin: 1.0",
					"  tiltangle: 0.0",
				},
			},
			want: Config{
				Physics: PhysicsConfig {
					Spin:             1,
					TiltAngle:        0.0,
				},
			},
		},
	}
	for _, tt := range tests {
		f, err := os.Create(tt.args.filename[0])
		if err != nil {
			t.Error(err)
		}
		defer f.Close()

		for _, line := range tt.args.lines {
			_, err := f.WriteString(line + "\n")
			if err != nil {
				t.Error(err)
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := LoadConfig(tt.args.filename); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadConfig() = %v, want %v", got, tt.want)
			}
		})
		err = os.Remove(tt.args.filename[0])
		if err != nil {
			t.Error(err)
		}
	}
}
