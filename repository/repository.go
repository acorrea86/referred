package repository

import (
	"blumer-ms-refers/contracts"
	"blumer-ms-refers/graph"
	"blumer-ms-refers/model"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Repository struct {
	Session contracts.NeoSession
}

func (r *Repository) Find(id string) (*model.Profile, error) {
	result, err := r.Session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (p:Profile {user_id: $id}) RETURN p", map[string]interface{}{
			"id": id,
		})
		if err != nil {
			return nil, err
		}

		record, err := result.Single()
		if err != nil {
			//...this error indicates there is no  record for this ID, so we consider that there is no error
			return nil, nil
		}

		return record, nil
	})
	if err != nil {
		return nil, err
	}

	profile := result.(model.Profile)
	return &profile, err
}

func (r *Repository) Save(data *model.Profile) error {
	_, err := r.Session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(
			"CREATE (p:Profile {user_id: $id, username: $username, is_active: $isActive, reward: $reward}) RETURN p",
			map[string]interface{}{
				"id":       data.UserID,
				"username": data.Username,
				"isActive": data.IsActive,
				"reward":   data.Reward,
			},
		)

		return nil, err
	})

	return err
}

func (r *Repository) Update(data *model.Profile) error {
	_, err := r.Session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(
			"MATCH (p:Profile {user_id: $id}) SET p.username = $username, "+
				"p.is_active = $isActive, p.reward = $reward RETURN p",
			map[string]interface{}{
				"id":       data.UserID,
				"username": data.Username,
				"isActive": data.IsActive,
				"reward":   data.Reward,
			},
		)
		return nil, err
	})

	return err
}

func (r *Repository) Delete(id string) error {
	_, err := r.Session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(
			"MATCH (p:Profile {user_id: $id}) SET p.is_active = $isActive RETURN p",
			map[string]interface{}{
				"id":       id,
				"isActive": false,
			},
		)
		return nil, err
	})

	return err
}

func (r *Repository) Refer(referred string, referrer string) *graph.ErrorExtensionParams {
	_, err := r.Session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err :=
			tx.Run(
				"MATCH (A:Profile {user_id:$referred}),(B:Profile {user_id:$referrer}) CREATE (A)-[r:[REFERRED]->(B)",
				map[string]interface{}{
					"referred": referred,
					"referrer": referrer,
				})

		return nil, err
	})
	if err != nil {
		return &graph.ErrorExtensionParams{AppError: graph.DataSourceError, Reason: err.Error()}
	}

	return nil
}

func (r *Repository) GetReferrer(referred string) (*model.Profile, error) {
	result, err := r.Session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(
			"MATCH (A:Profile {user_id: $referred})<-[:REFERRED]-(B:Profile) RETURN B.user_id, B.reward",
			map[string]interface{}{
				"referred": referred,
			})
		if err != nil {
			return nil, err
		}

		record, err := result.Single()
		if err != nil {
			//...this error indicates there is no  record for this ID, so we consider that there is no error
			return nil, nil
		}

		return record, nil
	})
	if err != nil {
		return nil, err
	}

	profile := result.(model.Profile)
	return &profile, nil
}

func NewRepository(session contracts.NeoSession) *Repository {
	return &Repository{
		Session: session,
	}
}
