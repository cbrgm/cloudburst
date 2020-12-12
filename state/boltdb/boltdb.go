package boltdb

import (
	"encoding/json"
	"fmt"
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

		// create new buckets with initState
		root, err := tx.CreateBucketIfNotExists([]byte(bucketScrapeTargets))
		if err != nil {
			return err
		}

		for _, scrapeTarget := range initState {

			// create new buckets with initState
			_, err = root.CreateBucketIfNotExists([]byte(bucketInstancesName(scrapeTarget.Name)))
			if err != nil {
				return err
			}

			// add each scrape target to the root bucket
			value, _ := json.Marshal(scrapeTarget)
			err = root.Put([]byte(scrapeTarget.Name), value)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return &BoltDB{db: db}, db.Close, err
}

func instancesBucketForScrapeTarget(tx *bolt.Tx, scrapeTarget string) *bolt.Bucket {
	b := tx.Bucket([]byte(bucketScrapeTargets))
	return b.Bucket([]byte(bucketInstancesName(scrapeTarget)))
}

func bucketInstancesName(scrapeTarget string) string {
	return fmt.Sprintf("%s-%s", scrapeTarget, bucketInstances)
}

func (bdb *BoltDB) ListScrapeTargets() ([]cloudburst.ScrapeTarget, error) {
	scrapeTargets := []cloudburst.ScrapeTarget{}

	err := bdb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketScrapeTargets))
		c := b.Cursor()

		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			if v != nil {
				var st cloudburst.ScrapeTarget
				_ = json.Unmarshal(v, &st)
				scrapeTargets = append(scrapeTargets, st)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return scrapeTargets, nil
}

func (bdb *BoltDB) GetInstance(scrapeTarget string, name string) (cloudburst.Instance, error) {
	var instance cloudburst.Instance

	err := bdb.db.View(func(tx *bolt.Tx) error {
		b := instancesBucketForScrapeTarget(tx, scrapeTarget)
		bytes := b.Get([]byte(name))

		if err := json.Unmarshal(bytes, &instance); err != nil {
			return err
		}
		return nil
	})

	return instance, err
}

func (bdb *BoltDB) GetInstances(scrapeTarget string) ([]cloudburst.Instance, error) {
	instances := []cloudburst.Instance{}

	err := bdb.db.View(func(tx *bolt.Tx) error {
		b := instancesBucketForScrapeTarget(tx, scrapeTarget)
		c := b.Cursor()

		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			var instance cloudburst.Instance
			_ = json.Unmarshal(v, &instance)
			instances = append(instances, instance)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return instances, nil
}

func (bdb *BoltDB) SaveInstance(scrapeTarget string, instance cloudburst.Instance) (cloudburst.Instance, error) {
	err := bdb.db.Update(func(tx *bolt.Tx) error {
		b := instancesBucketForScrapeTarget(tx, scrapeTarget)

		key := instance.Name
		value, _ := json.Marshal(instance)

		return b.Put([]byte(key), value)
	})

	return instance, err
}

func (bdb *BoltDB) SaveInstances(scrapeTarget string, instances []cloudburst.Instance) ([]cloudburst.Instance, error) {
	var res []cloudburst.Instance
	for _, instance := range instances {
		updated, err := bdb.SaveInstance(scrapeTarget, instance)
		if err != nil {
			return nil, err
		}
		res = append(res, updated)
	}
	return res, nil
}

func (bdb *BoltDB) RemoveInstances(scrapeTarget string, instances []cloudburst.Instance) error {
	for _, instance := range instances {
		err := bdb.RemoveInstance(scrapeTarget, instance)
		if err != nil {
			return err
		}
	}
	return nil
}
func (bdb *BoltDB) RemoveInstance(scrapeTarget string, instance cloudburst.Instance) error {
	err := bdb.db.Update(func(tx *bolt.Tx) error {
		b := instancesBucketForScrapeTarget(tx, scrapeTarget)
		key := []byte(instance.Name)
		return b.Delete(key)
	})
	return err
}
