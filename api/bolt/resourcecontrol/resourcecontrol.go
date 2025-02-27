package resourcecontrol

import (
	"github.com/hazik1024/portainer/api"
	"github.com/hazik1024/portainer/api/bolt/internal"

	"github.com/boltdb/bolt"
)

const (
	// BucketName represents the name of the bucket where this service stores data.
	BucketName = "resource_control"
)

// Service represents a service for managing endpoint data.
type Service struct {
	db *bolt.DB
}

// NewService creates a new instance of a service.
func NewService(db *bolt.DB) (*Service, error) {
	err := internal.CreateBucket(db, BucketName)
	if err != nil {
		return nil, err
	}

	return &Service{
		db: db,
	}, nil
}

// ResourceControl returns a ResourceControl object by ID
func (service *Service) ResourceControl(ID portainer.ResourceControlID) (*portainer.ResourceControl, error) {
	var resourceControl portainer.ResourceControl
	identifier := internal.Itob(int(ID))

	err := internal.GetObject(service.db, BucketName, identifier, &resourceControl)
	if err != nil {
		return nil, err
	}

	return &resourceControl, nil
}

// ResourceControlByResourceID returns a ResourceControl object by checking if the resourceID is equal
// to the main ResourceID or in SubResourceIDs
func (service *Service) ResourceControlByResourceID(resourceID string) (*portainer.ResourceControl, error) {
	var resourceControl *portainer.ResourceControl

	err := service.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		cursor := bucket.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var rc portainer.ResourceControl
			err := internal.UnmarshalObject(v, &rc)
			if err != nil {
				return err
			}

			if rc.ResourceID == resourceID {
				resourceControl = &rc
				break
			}

			for _, subResourceID := range rc.SubResourceIDs {
				if subResourceID == resourceID {
					resourceControl = &rc
					break
				}
			}
		}

		if resourceControl == nil {
			return portainer.ErrObjectNotFound
		}

		return nil
	})

	return resourceControl, err
}

// ResourceControls returns all the ResourceControl objects
func (service *Service) ResourceControls() ([]portainer.ResourceControl, error) {
	var rcs = make([]portainer.ResourceControl, 0)

	err := service.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var resourceControl portainer.ResourceControl
			err := internal.UnmarshalObject(v, &resourceControl)
			if err != nil {
				return err
			}
			rcs = append(rcs, resourceControl)
		}

		return nil
	})

	return rcs, err
}

// CreateResourceControl creates a new ResourceControl object
func (service *Service) CreateResourceControl(resourceControl *portainer.ResourceControl) error {
	return service.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		id, _ := bucket.NextSequence()
		resourceControl.ID = portainer.ResourceControlID(id)

		data, err := internal.MarshalObject(resourceControl)
		if err != nil {
			return err
		}

		return bucket.Put(internal.Itob(int(resourceControl.ID)), data)
	})
}

// UpdateResourceControl saves a ResourceControl object.
func (service *Service) UpdateResourceControl(ID portainer.ResourceControlID, resourceControl *portainer.ResourceControl) error {
	identifier := internal.Itob(int(ID))
	return internal.UpdateObject(service.db, BucketName, identifier, resourceControl)
}

// DeleteResourceControl deletes a ResourceControl object by ID
func (service *Service) DeleteResourceControl(ID portainer.ResourceControlID) error {
	identifier := internal.Itob(int(ID))
	return internal.DeleteObject(service.db, BucketName, identifier)
}
