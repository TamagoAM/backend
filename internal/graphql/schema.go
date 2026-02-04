package graphql

import (
	"time"

	"github.com/graphql-go/graphql"
	"github.com/jmoiron/sqlx"

	"tamagoam/internal/models"
)

type CreateUserInput struct {
	Name               string
	LastName           string
	UserName           string
	Email              string
	ProfilPicture      *string
	GamingTime         int
	LastConnectionDate *time.Time
}

func NewSchema(db *sqlx.DB) (graphql.Schema, error) {
	store := NewSQLStore(db)

	userType := graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if u, ok := p.Source.(models.User); ok {
					return u.UserID, nil
				}
				return nil, nil
			}},
			"name": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if u, ok := p.Source.(models.User); ok {
					return u.Name, nil
				}
				return nil, nil
			}},
			"lastName": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if u, ok := p.Source.(models.User); ok {
					return u.LastName, nil
				}
				return nil, nil
			}},
			"userName": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if u, ok := p.Source.(models.User); ok {
					return u.UserName, nil
				}
				return nil, nil
			}},
			"email": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if u, ok := p.Source.(models.User); ok {
					return u.Email, nil
				}
				return nil, nil
			}},
			"profilPicture": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if u, ok := p.Source.(models.User); ok {
					return u.ProfilPicture, nil
				}
				return nil, nil
			}},
			"gamingTime": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if u, ok := p.Source.(models.User); ok {
					return u.GamingTime, nil
				}
				return nil, nil
			}},
			"creationDate": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if u, ok := p.Source.(models.User); ok {
					return formatTimeValue(&u.CreationDate), nil
				}
				return nil, nil
			}},
			"lastConnectionDate": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if u, ok := p.Source.(models.User); ok {
					return formatTimeValue(u.LastConnectionDate), nil
				}
				return nil, nil
			}},
		},
	})

	raceType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Race",
		Fields: graphql.Fields{
			"id": &graphql.Field{Type: graphql.Int, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if r, ok := p.Source.(models.Race); ok {
					return r.RaceID, nil
				}
				return nil, nil
			}},
			"name": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if r, ok := p.Source.(models.Race); ok {
					return r.Name, nil
				}
				return nil, nil
			}},
			"desc": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if r, ok := p.Source.(models.Race); ok {
					return r.Desc, nil
				}
				return nil, nil
			}},
			"bonus": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if r, ok := p.Source.(models.Race); ok {
					return r.Bonus, nil
				}
				return nil, nil
			}},
			"malus": &graphql.Field{Type: graphql.String, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if r, ok := p.Source.(models.Race); ok {
					return r.Malus, nil
				}
				return nil, nil
			}},
		},
	})

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"users": &graphql.Field{
				Type: graphql.NewList(userType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return store.ListUsers(p.Context)
				},
			},
			"user": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					return store.GetUser(p.Context, id)
				},
			},
			"races": &graphql.Field{
				Type: graphql.NewList(raceType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return store.ListRaces(p.Context)
				},
			},
		},
	})

	createUserInput := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "CreateUserInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"name":               &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"lastName":           &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"userName":           &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"email":              &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"profilPicture":      &graphql.InputObjectFieldConfig{Type: graphql.String},
			"gamingTime":         &graphql.InputObjectFieldConfig{Type: graphql.Int},
			"lastConnectionDate": &graphql.InputObjectFieldConfig{Type: graphql.String},
		},
	})

	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createUser": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(createUserInput)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					inputMap := p.Args["input"].(map[string]interface{})
					input := CreateUserInput{
						Name:     inputMap["name"].(string),
						LastName: inputMap["lastName"].(string),
						UserName: inputMap["userName"].(string),
						Email:    inputMap["email"].(string),
					}
					if v, ok := inputMap["profilPicture"]; ok {
						if s, ok := v.(string); ok {
							input.ProfilPicture = &s
						}
					}
					if v, ok := inputMap["gamingTime"]; ok {
						if i, ok := v.(int); ok {
							input.GamingTime = i
						}
					}
					if v, ok := inputMap["lastConnectionDate"]; ok {
						if s, ok := v.(string); ok && s != "" {
							if t, err := time.Parse(time.RFC3339, s); err == nil {
								input.LastConnectionDate = &t
							}
						}
					}
					return store.CreateUser(p.Context, input)
				},
			},
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
}

func formatTimeValue(t *time.Time) interface{} {
	if t == nil {
		return nil
	}
	return t.UTC().Format(time.RFC3339)
}
