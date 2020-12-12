package boltdb

import (
	"github.com/cbrgm/cloudburst/cloudburst"
	"reflect"
	"testing"
	"time"
)

func initState() []cloudburst.ScrapeTarget {
	return []cloudburst.ScrapeTarget{
		{
			Name:        "foo",
			Description: "",
			Query:       "null",
			InstanceSpec: cloudburst.InstanceSpec{
				Container: cloudburst.ContainerSpec{
					Name:  "foo-container",
					Image: "cbrgm/example-app",
				},
			},
		},
		{
			Name:        "bar",
			Description: "",
			Query:       "null",
			InstanceSpec: cloudburst.InstanceSpec{
				Container: cloudburst.ContainerSpec{
					Name:  "bar-container",
					Image: "cbrgm/example-app",
				},
			},
		},
	}
}

func testInstances() []cloudburst.Instance {
	return []cloudburst.Instance{
		{
			Name:     "foo-instance",
			Endpoint: "example.com",
			Target:   "",
			Active:   true,
			Status: cloudburst.InstanceStatus{
				Agent:   "fake-agent",
				Status:  cloudburst.Pending,
				Started: time.Now(),
			},
		},
		{
			Name:     "bar-instance",
			Endpoint: "example.com",
			Target:   "",
			Active:   true,
			Status: cloudburst.InstanceStatus{
				Agent:   "fake-agent",
				Status:  cloudburst.Pending,
				Started: time.Now(),
			},
		},
	}
}

func TestNew(t *testing.T) {
	_, dbClose, err := NewDB("../../development/data", initState())
	if err != nil {
		t.Errorf("err %s", err)
	}
	defer dbClose()
}

func TestListScrapeTargets(t *testing.T) {
	db, dbClose, err := NewDB("../../development/data", initState())
	if err != nil {
		t.Errorf("err %s", err)
	}
	defer dbClose()

	type args struct {
		scrapeTarget []cloudburst.ScrapeTarget
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test list scrape instances",
			args: args{
				scrapeTarget: initState(),
			},
			want: true,
		},
	}

	scrapeTargets, err := db.ListScrapeTargets()
	if err != nil {
		t.Errorf("err %s", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reflect.DeepEqual(scrapeTargets, tt.args.scrapeTarget); !got {
				t.Errorf("%s, got %v, want %v", "TestListScrapeTargets()", got, tt.want)
			}
		})
	}
}

func TestSaveGetInstance(t *testing.T) {
	scrapeTargets := initState()
	instances := testInstances()
	db, dbClose, err := NewDB("../../development/data", scrapeTargets)
	if err != nil {
		t.Errorf("err %s", err)
	}
	defer dbClose()

	type args struct {
		scrapeTarget cloudburst.ScrapeTarget
		instance     cloudburst.Instance
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test save and get single instance",
			args: args{
				scrapeTarget: scrapeTargets[0],
				instance:     instances[0],
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			saved, err := db.SaveInstance(tt.args.scrapeTarget.Name, tt.args.instance)
			if err != nil {
				t.Errorf("%s, %s", "TestSaveGetInstance()", err)
			}

			result, err := db.GetInstance(tt.args.scrapeTarget.Name, saved.Name)
			if err != nil {
				t.Errorf("%s, %s", "TestSaveGetInstance()", err)
			}

			if got := reflect.DeepEqual(result, tt.args.instance); got == tt.want {
				t.Errorf("%s, got %v, want %v", "TestSaveGetInstance()", got, tt.want)
			}
		})
	}
}

func TestSaveGetInstances(t *testing.T) {
	scrapeTargets := initState()
	instances := testInstances()
	db, dbClose, err := NewDB("../../development/data", scrapeTargets)
	if err != nil {
		t.Errorf("err %s", err)
	}
	defer dbClose()

	type args struct {
		scrapeTargets []cloudburst.ScrapeTarget
		instances     []cloudburst.Instance
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test save and get single instance",
			args: args{
				scrapeTargets: scrapeTargets,
				instances:     instances,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			for _, target := range tt.args.scrapeTargets {

				_, err := db.SaveInstances(target.Name, tt.args.instances)
				if err != nil {
					t.Errorf("%s, %s", "TestSaveGetInstances()", err)
				}

				result, err := db.GetInstances(target.Name)
				if err != nil {
					t.Errorf("%s, %s", "TestSaveGetInstances()", err)
				}

				if got := reflect.DeepEqual(result, tt.args.instances); got == tt.want {
					t.Errorf("%s, got %v, want %v", "TestSaveGetInstances()", got, tt.want)
				}
			}

		})
	}
}