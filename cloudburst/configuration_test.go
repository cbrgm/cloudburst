package cloudburst

import (
	"reflect"
	"testing"
)

func TestParseConfiguration(t *testing.T) {
	type args struct {
		config Configuration
	}
	tests := []struct {
		name    string
		args    args
		want    []*ScrapeTarget
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseConfiguration(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseConfiguration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseConfiguration() got = %v, want %v", got, tt.want)
			}
		})
	}
}