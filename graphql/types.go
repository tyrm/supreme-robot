package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/tyrm/supreme-robot/models"
)

type success struct {
	Success bool `json:"success"`
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

var statusType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Status",
	Fields: graphql.Fields{
		"version": &graphql.Field{
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
		},
		"priority": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if r, ok := p.Source.(models.Record); ok {
					if r.Priority.Valid {
						return r.Priority.Int32, nil
					}
					return nil, nil
				}
				if r, ok := p.Source.(*models.Record); ok {
					if r.Priority.Valid {
						return r.Priority.Int32, nil
					}
					return nil, nil
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
					}
					return nil, nil
				}
				if r, ok := p.Source.(*models.Record); ok {
					if r.Port.Valid {
						return r.Port.Int32, nil
					}
					return nil, nil
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
					}
					return nil, nil
				}
				if r, ok := p.Source.(*models.Record); ok {
					if r.Weight.Valid {
						return r.Weight.Int32, nil
					}
					return nil, nil
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
					}
					return nil, nil
				}
				if r, ok := p.Source.(*models.Record); ok {
					if r.Refresh.Valid {
						return r.Refresh.Int32, nil
					}
					return nil, nil
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
					}
					return nil, nil
				}
				if r, ok := p.Source.(*models.Record); ok {
					if r.Retry.Valid {
						return r.Retry.Int32, nil
					}
					return nil, nil
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
					}
					return nil, nil
				}
				if r, ok := p.Source.(*models.Record); ok {
					if r.Expire.Valid {
						return r.Expire.Int32, nil
					}
					return nil, nil
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
					}
					return nil, nil
				}
				if r, ok := p.Source.(*models.Record); ok {
					if r.MBox.Valid {
						return r.MBox.String, nil
					}
					return nil, nil
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
					}
					return nil, nil
				}
				if r, ok := p.Source.(*models.Record); ok {
					if r.Tag.Valid {
						return r.Tag.String, nil
					}
					return nil, nil
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
