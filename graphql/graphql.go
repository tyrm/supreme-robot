package graphql

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/tyrm/supreme-robot/models"
	"net/http"
)

type postData struct {
	Query     string                 `json:"query"`
	Operation string                 `json:"operation"`
	Variables map[string]interface{} `json:"variables"`
}

// input types
var soaInputType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "SOA",
		Fields: graphql.InputObjectConfigFieldMap{
			"ttl": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"mbox": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"refresh": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"retry": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"expire": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
	},
)

// types
var domainType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Domain",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"domain": &graphql.Field{
			Type: graphql.String,
		},
		"records": &graphql.Field{
			Type: graphql.NewList(recordType),
		},
		"createdAt": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if d, ok := p.Source.(models.Domain); ok {
					return d.CreatedAt.Unix(), nil
				}
				if d, ok := p.Source.(*models.Domain); ok {
					return d.CreatedAt.Unix(), nil
				}
				return nil, nil
			},
		},
		"updatedAt": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if d, ok := p.Source.(models.Domain); ok {
					return d.UpdatedAt.Unix(), nil
				}
				if d, ok := p.Source.(*models.Domain); ok {
					return d.UpdatedAt.Unix(), nil
				}
				return nil, nil
			},
		},
	},
})

var jwtTokensType = graphql.NewObject(graphql.ObjectConfig{
	Name: "JwtToken",
	Fields: graphql.Fields{
		"accessToken": &graphql.Field{
			Type: graphql.String,
		},
		"refreshToken": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var successType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Success",
	Fields: graphql.Fields{
		"success": &graphql.Field{
			Type: graphql.Boolean,
		},
	},
})

var recordType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Record",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"type": &graphql.Field{
			Type: graphql.String,
		},
		"value": &graphql.Field{
			Type: graphql.String,
		},
		"ttl": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if r, ok := p.Source.(models.Record); ok {
					if r.TTL.Valid {
						return r.TTL.Int32, nil
					} else {
						return nil, nil
					}
				}
				if r, ok := p.Source.(*models.Record); ok {
					if r.TTL.Valid {
						return r.TTL.Int32, nil
					} else {
						return nil, nil
					}
				}
				return nil, nil
			},
		},
		"priority": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if r, ok := p.Source.(models.Record); ok {
					if r.Priority.Valid {
						return r.Priority.Int32, nil
					} else {
						return nil, nil
					}
				}
				if r, ok := p.Source.(*models.Record); ok {
					if r.Priority.Valid {
						return r.Priority.Int32, nil
					} else {
						return nil, nil
					}
				}
				return nil, nil
			},
		},
		"port": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if r, ok := p.Source.(models.Record); ok {
					if r.Port.Valid {
						return r.Port.Int32, nil
					} else {
						return nil, nil
					}
				}
				if r, ok := p.Source.(*models.Record); ok {
					if r.Port.Valid {
						return r.Port.Int32, nil
					} else {
						return nil, nil
					}
				}
				return nil, nil
			},
		},
		"weight": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if r, ok := p.Source.(models.Record); ok {
					if r.Weight.Valid {
						return r.Weight.Int32, nil
					} else {
						return nil, nil
					}
				}
				if r, ok := p.Source.(*models.Record); ok {
					if r.Weight.Valid {
						return r.Weight.Int32, nil
					} else {
						return nil, nil
					}
				}
				return nil, nil
			},
		},
		"refresh": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if r, ok := p.Source.(models.Record); ok {
					if r.Refresh.Valid {
						return r.Refresh.Int32, nil
					} else {
						return nil, nil
					}
				}
				if r, ok := p.Source.(*models.Record); ok {
					if r.Refresh.Valid {
						return r.Refresh.Int32, nil
					} else {
						return nil, nil
					}
				}
				return nil, nil
			},
		},
		"retry": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if r, ok := p.Source.(models.Record); ok {
					if r.Retry.Valid {
						return r.Retry.Int32, nil
					} else {
						return nil, nil
					}
				}
				if r, ok := p.Source.(*models.Record); ok {
					if r.Retry.Valid {
						return r.Retry.Int32, nil
					} else {
						return nil, nil
					}
				}
				return nil, nil
			},
		},
		"expire": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if r, ok := p.Source.(models.Record); ok {
					if r.Expire.Valid {
						return r.Expire.Int32, nil
					} else {
						return nil, nil
					}
				}
				if r, ok := p.Source.(*models.Record); ok {
					if r.Expire.Valid {
						return r.Expire.Int32, nil
					} else {
						return nil, nil
					}
				}
				return nil, nil
			},
		},
		"mbox": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if r, ok := p.Source.(models.Record); ok {
					if r.MBox.Valid {
						return r.MBox.String, nil
					} else {
						return nil, nil
					}
				}
				if r, ok := p.Source.(*models.Record); ok {
					if r.MBox.Valid {
						return r.MBox.String, nil
					} else {
						return nil, nil
					}
				}
				return nil, nil
			},
		},
		"tag": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if r, ok := p.Source.(models.Record); ok {
					if r.Tag.Valid {
						return r.Tag.String, nil
					} else {
						return nil, nil
					}
				}
				if r, ok := p.Source.(*models.Record); ok {
					if r.Tag.Valid {
						return r.Tag.String, nil
					} else {
						return nil, nil
					}
				}
				return nil, nil
			},
		},
		"createdAt": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if r, ok := p.Source.(models.Record); ok {
					return r.CreatedAt.Unix(), nil
				}
				if r, ok := p.Source.(*models.Record); ok {
					return r.CreatedAt.Unix(), nil
				}
				return nil, nil
			},
		},
		"updatedAt": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if r, ok := p.Source.(models.Record); ok {
					return r.UpdatedAt.Unix(), nil
				}
				if r, ok := p.Source.(*models.Record); ok {
					return r.UpdatedAt.Unix(), nil
				}
				return nil, nil
			},
		},
	},
})

var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"username": &graphql.Field{
			Type: graphql.String,
		},
		"groups": &graphql.Field{
			Type: graphql.NewList(graphql.String),
		},
		"createdAt": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if u, ok := p.Source.(models.User); ok {
					return u.CreatedAt.Unix(), nil
				}
				if u, ok := p.Source.(*models.User); ok {
					return u.CreatedAt.Unix(), nil
				}
				return nil, nil
			},
		},
		"updatedAt": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if u, ok := p.Source.(models.User); ok {
					return u.UpdatedAt.Unix(), nil
				}
				if u, ok := p.Source.(*models.User); ok {
					return u.UpdatedAt.Unix(), nil
				}
				return nil, nil
			},
		},
	},
})

// root mutation
func (s *Server) rootMutation() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"addDomain": &graphql.Field{
				Type:        domainType,
				Description: "Add new domain",
				Args: graphql.FieldConfigArgument{
					"domain": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"soa": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(soaInputType),
					},
				},
				Resolve: s.addDomainMutator,
			},

			"addUser": &graphql.Field{
				Type:        userType,
				Description: "Add new user",
				Args: graphql.FieldConfigArgument{
					"username": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"groups": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.String),
					},
				},
				Resolve: s.addUserMutation,
			},

			"login": &graphql.Field{
				Type:        jwtTokensType,
				Description: "Login to system",
				Args: graphql.FieldConfigArgument{
					"username": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: s.loginMutator,
			},

			"logout": &graphql.Field{
				Type:        successType,
				Description: "Logout of the system",
				Resolve:     s.logoutMutator,
			},

			"refreshAccessToken": &graphql.Field{
				Type:        jwtTokensType,
				Description: "Refresh jwt token",
				Args: graphql.FieldConfigArgument{
					"refreshToken": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: s.refreshAccessTokenMutator,
			},
		},
	})
}

// root queries
func (s *Server) rootQuery() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"domain": &graphql.Field{
				Type:        domainType,
				Description: "Get info about a domain",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: s.domainQuery,
			},

			"me": &graphql.Field{
				Type:        userType,
				Description: "Get logged in user",
				Resolve:     s.meQuery,
			},

			"myDomains": &graphql.Field{
				Type:        graphql.NewList(domainType),
				Description: "Get my domains",
				Resolve:     s.myDomainsQuery,
			},

			"user": &graphql.Field{
				Type:        userType,
				Description: "Get single user",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: s.userQuery,
			},
		},
	})
}

func (s *Server) schema() graphql.Schema {
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    s.rootQuery(),
		Mutation: s.rootMutation(),
	})
	return schema
}

func (s *Server) graphqlHandler(w http.ResponseWriter, r *http.Request) {
	var p postData
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		w.WriteHeader(400)
		return
	}

	logger.Tracef("query: %s", p.Query)
	logger.Tracef("operation: %s", p.Operation)
	logger.Tracef("variables: %v", p.Variables)

	ctx := r.Context()

	// check auth
	metadata, err := s.extractTokenMetadata(r)
	if err == nil {
		// if success add metadata to context
		ctx = context.WithValue(ctx, MetadataKey, metadata)
	}

	result := graphql.Do(graphql.Params{
		Context:        ctx,
		Schema:         s.schema(),
		RequestString:  p.Query,
		VariableValues: p.Variables,
		OperationName:  p.Operation,
	})
	if err := json.NewEncoder(w).Encode(result); err != nil {
		fmt.Printf("could not write result to response: %s", err)
	}
}
