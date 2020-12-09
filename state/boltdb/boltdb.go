package boltdb

import (
	"encoding/json"
	"github.com/cbrgm/cloudburst/cloudburst"
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
)

const (
	bucketScrapeTargets = `scrapetargets`
	bucketInstances     = `instances`
)

type BoltDB struct {
	db *bolt.DB
}

func NewDB(path string, initState []cloudburst.ScrapeTarget) (*BoltDB, func() error, error) {
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return nil, nil, errors.Errorf("failed to open boltdb: %s", err)
	}

	err = db.Batch(func(tx *bolt.Tx) error {
		// reset rootBucket
		_ = tx.DeleteBucket([]byte(bucketScrapeTargets))
		_ = tx.DeleteBucket([]byte(bucketInstances))

		// create new buckets with initState
		targets, err := tx.CreateBucketIfNotExists([]byte(bucketScrapeTargets))
		if err != nil {
			return err
		}

		// create new buckets with initState
		_, err = tx.CreateBucketIfNotExists([]byte(bucketInstances))
		if err != nil {
			return err
		}

		for _, st := range initState {
			value, _ := json.Marshal(st)
			err = targets.Put([]byte(st.Name), value)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return &BoltDB{db: db}, db.Close, err
}

func (bdb *BoltDB) ListScrapeTargets() ([]cloudburst.ScrapeTarget, error) {
	var scrapeTargets []cloudburst.ScrapeTarget

	err := bdb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketScrapeTargets))
		c := b.Cursor()

		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			var st cloudburst.ScrapeTarget
			_ = json.Unmarshal(v, &st)
			scrapeTargets = append(scrapeTargets, st)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return scrapeTargets, nil
}

func (bdb *BoltDB) GetInstance(scrapeTarget string) []cloudburst.Instance {
	return nil	
}

func (bdb *BoltDB) GetInstances(scrapeTarget string) []cloudburst.Instance {
	return nil
}

func (bdb *BoltDB) UpdateInstances(scrapeTarget string, instances []cloudburst.Instance) ([]cloudburst.Instance, error) {
	return nil, nil
}
func (bdb *BoltDB) RemoveInstances(scrapeTarget string, instances []cloudburst.Instance) error {
	return nil
}
func (bdb *BoltDB) RemoveInstance(scrapeTarget string, instance cloudburst.Instance) error {
	return nil
}
func (bdb *BoltDB) CreateInstances(scrapeTarget string, instances []cloudburst.Instance) ([]cloudburst.Instance, error) {
	return nil, nil
}
func (bdb *BoltDB) CreateInstance(scrapeTarget string, instance cloudburst.Instance) (cloudburst.Instance, error) {
	return cloudburst.Instance{}, nil
}