package graphql

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
	"net/http"
)

type postData struct {
	Query     string                 `json:"query"`
	Operation string                 `json:"operation"`
	Variables map[string]interface{} `json:"variables"`
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

			"addRecordA": &graphql.Field{
				Type:        recordType,
				Description: "Add A record",
				Args: graphql.FieldConfigArgument{
					"domainId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"ip": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"ttl": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: s.addRecordAMutator,
			},

			"addRecordAAAA": &graphql.Field{
				Type:        recordType,
				Description: "Add AAAA record",
				Args: graphql.FieldConfigArgument{
					"domainId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"ip": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"ttl": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: s.addRecordAAAAMutator,
			},

			"addRecordCNAME": &graphql.Field{
				Type:        recordType,
				Description: "Add CNAME record",
				Args: graphql.FieldConfigArgument{
					"domainId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"host": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"ttl": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: s.addRecordCNAMEMutator,
			},

			"addRecordMX": &graphql.Field{
				Type:        recordType,
				Description: "Add MX record",
				Args: graphql.FieldConfigArgument{
					"domainId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"host": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"priority": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"ttl": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: s.addRecordMXMutator,
			},

			"addRecordNS": &graphql.Field{
				Type:        recordType,
				Description: "Add NS record",
				Args: graphql.FieldConfigArgument{
					"domainId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"host": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"ttl": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: s.addRecordNSMutator,
			},

			"addRecordSRV": &graphql.Field{
				Type:        recordType,
				Description: "Add SRV record",
				Args: graphql.FieldConfigArgument{
					"domainId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"host": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"port": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"priority": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"weight": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"ttl": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: s.addRecordSRVMutator,
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

			"deleteDomain": &graphql.Field{
				Type:        successType,
				Description: "Delete domain",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: s.deleteDomainMutator,
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
		ctx = context.WithValue(ctx, metadataKey, metadata)
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
