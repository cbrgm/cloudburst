package cloudburst

import (
	"testing"
	"time"
)

func testInstances() []*Instance {
	return []*Instance{
		{
			Name:     "foo-instance",
			Endpoint: "example.com",
			Active:   true,
			Status: InstanceStatus{
				Agent:   "fake-agent",
				Status:  Pending,
				Started: time.Now(),
			},
		},
		{
			Name:     "foo-instance",
			Endpoint: "example.com",
			Active:   true,
			Status: InstanceStatus{
				Agent:   "fake-agent",
				Status:  Progress,
				Started: time.Now(),
			},
		},
		{
			Name:     "bar-instance",
			Endpoint: "example.com",
			Active:   true,
			Status: InstanceStatus{
				Agent:   "fake-agent",
				Status:  Running,
				Started: time.Now(),
			},
		},
		{
			Name:     "foobar-instance",
			Endpoint: "example.com",
			Active:   false,
			Status: InstanceStatus{
				Agent:   "fake-agent",
				Status:  Running,
				Started: time.Now(),
			},
		},
		{
			Name:     "foobar-instance",
			Endpoint: "example.com",
			Active:   true,
			Status: InstanceStatus{
				Agent:   "fake-agent",
				Status:  Pending,
				Started: time.Now(),
			},
		},
		{
			Name:     "foobar-instance",
			Endpoint: "example.com",
			Active:   true,
			Status: InstanceStatus{
				Agent:   "fake-agent",
				Status:  Failure,
				Started: time.Now(),
			},
		},
		{
			Name:     "foobar-instance",
			Endpoint: "example.com",
			Active:   false,
			Status: InstanceStatus{
				Agent:   "fake-agent",
				Status:  Terminated,
				Started: time.Now(),
			},
		},
	}
}

func TestCountInstancesByActiveStatus(t *testing.T) {

	type args struct {
		active    bool
		instances []*Instance
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "get all active false",
			args: args{
				active:    false,
				instances: testInstances(),
			},
			want: 2,
		},
		{
			name: "get all active true",
			args: args{
				active:    true,
				instances: testInstances(),
			},
			want: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CountInstancesByActiveStatus(tt.args.instances, tt.args.active)
			if got != tt.want {
				t.Errorf("want %d, got %d", tt.want, got)
			}
		})
	}
}

func TestCountInstancesByStatus(t *testing.T) {
	type args struct {
		status    Status
		instances []*Instance
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "get all pending",
			args: args{
				status:    Pending,
				instances: testInstances(),
			},
			want: 2,
		},
		{
			name: "get all process",
			args: args{
				status:    Progress,
				instances: testInstances(),
			},
			want: 1,
		},
		{
			name: "get all running",
			args: args{
				status:    Running,
				instances: testInstances(),
			},
			want: 2,
		},
		{
			name: "get all terminated",
			args: args{
				status:    Terminated,
				instances: testInstances(),
			},
			want: 1,
		},
		{
			name: "get all failure",
			args: args{
				status:    Failure,
				instances: testInstances(),
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CountInstancesByStatus(tt.args.instances, tt.args.status)
			if got != tt.want {
				t.Errorf("want %d, got %d", tt.want, got)
			}
		})
	}
}

func TestThresholdInRange(t *testing.T) {
	type args struct {
		queryResult int
		threshold   Threshold
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "-2 in range upper 0, lower -2",
			args: args{
				threshold:   newThreshold(0, -2),
				queryResult: -2,
			},
			want: true,
		},
		{
			name: "0 in range upper 0, lower -2",
			args: args{
				threshold:   newThreshold(0, -2),
				queryResult: -0,
			},
			want: true,
		},
		{
			name: "0 in range upper 0, lower -2",
			args: args{
				threshold:   newThreshold(0, -2),
				queryResult: 1,
			},
			want: false,
		},
		{
			name: "0 in range upper 0, lower -2",
			args: args{
				threshold:   newThreshold(0, -2),
				queryResult: -3,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.threshold.inRange(tt.args.queryResult)
			if got != tt.want {
				t.Errorf("want %t, got %t", tt.want, got)
			}
		})
	}
}

func TestGetInstancesByActiveStatus(t *testing.T) {

}

func TestGetInstancesByStatus(t *testing.T) {

}

func TestIsMatchingStatus(t *testing.T) {

}
