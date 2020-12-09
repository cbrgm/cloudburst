package boltdb

import (
	"github.com/cbrgm/cloudburst/cloudburst"
	"reflect"
	"testing"
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
			name: "test list scrape scrapeTargets",
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
