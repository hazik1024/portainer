package team

import (
	"github.com/hazik1024/portainer/api"
	"github.com/hazik1024/portainer/api/bolt/internal"

	"github.com/boltdb/bolt"
)

const (
	// BucketName represents the name of the bucket where this service stores data.
	BucketName = "teams"
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

// Team returns a Team by ID
func (service *Service) Team(ID portainer.TeamID) (*portainer.Team, error) {
	var team portainer.Team
	identifier := internal.Itob(int(ID))

	err := internal.GetObject(service.db, BucketName, identifier, &team)
	if err != nil {
		return nil, err
	}

	return &team, nil
}

// TeamByName returns a team by name.
func (service *Service) TeamByName(name string) (*portainer.Team, error) {
	var team *portainer.Team

	err := service.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var t portainer.Team
			err := internal.UnmarshalObject(v, &t)
			if err != nil {
				return err
			}

			if t.Name == name {
				team = &t
				break
			}
		}

		if team == nil {
			return portainer.ErrObjectNotFound
		}

		return nil
	})

	return team, err
}

// Teams return an array containing all the teams.
func (service *Service) Teams() ([]portainer.Team, error) {
	var teams = make([]portainer.Team, 0)

	err := service.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var team portainer.Team
			err := internal.UnmarshalObject(v, &team)
			if err != nil {
				return err
			}
			teams = append(teams, team)
		}

		return nil
	})

	return teams, err
}

// UpdateTeam saves a Team.
func (service *Service) UpdateTeam(ID portainer.TeamID, team *portainer.Team) error {
	identifier := internal.Itob(int(ID))
	return internal.UpdateObject(service.db, BucketName, identifier, team)
}

// CreateTeam creates a new Team.
func (service *Service) CreateTeam(team *portainer.Team) error {
	return service.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		id, _ := bucket.NextSequence()
		team.ID = portainer.TeamID(id)

		data, err := internal.MarshalObject(team)
		if err != nil {
			return err
		}

		return bucket.Put(internal.Itob(int(team.ID)), data)
	})
}

// DeleteTeam deletes a Team.
func (service *Service) DeleteTeam(ID portainer.TeamID) error {
	identifier := internal.Itob(int(ID))
	return internal.DeleteObject(service.db, BucketName, identifier)
}
