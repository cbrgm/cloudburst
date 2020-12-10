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
	scrapeTargets := []cloudburst.ScrapeTarget{}

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

func (bdb *BoltDB) GetInstance(name string) (cloudburst.Instance, error) {
	var instance cloudburst.Instance

	err := bdb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketInstances))
		bytes := b.Get([]byte(name))

		if err := json.Unmarshal(bytes, &instance); err != nil {
			return err
		}
		return nil
	})

	return instance, err
}

func (bdb *BoltDB) GetInstancesForTarget(scrapeTarget string) ([]cloudburst.Instance, error) {
	instances := []cloudburst.Instance{}

	err := bdb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketInstances))
		c := b.Cursor()

		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			var instance cloudburst.Instance
			_ = json.Unmarshal(v, &instance)
			if scrapeTarget == instance.Target {
				instances = append(instances, instance)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return instances, nil
}

func (bdb *BoltDB) SaveInstance(instance cloudburst.Instance) (cloudburst.Instance, error) {
	err := bdb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketInstances))

		key := instance.Name
		value, _ := json.Marshal(instance)

		return b.Put([]byte(key), value)
	})

	return instance, err
}

func (bdb *BoltDB) SaveInstances(instances []cloudburst.Instance) ([]cloudburst.Instance, error) {
	var res []cloudburst.Instance
	for _, instance := range instances {
		updated, err := bdb.SaveInstance(instance)
		if err != nil {
			return nil, err
		}
		res = append(res, updated)
	}
	return res, nil
}

func (bdb *BoltDB) RemoveInstances(instances []cloudburst.Instance) error {
	for _, instance := range instances {
		err := bdb.RemoveInstance(instance)
		if err != nil {
			return err
		}
	}
	return nil
}
func (bdb *BoltDB) RemoveInstance(instance cloudburst.Instance) error {
	err := bdb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketInstances))
		key := []byte(instance.Name)
		return b.Delete(key)
	})
	return err
}
