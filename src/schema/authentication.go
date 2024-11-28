package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Authentication struct {
	ent.Schema
}

func (Authentication) Fields() []ent.Field {
	return []ent.Field{
		field.Int("groupId").Immutable(),
		field.Int("id"),
		field.String("name"),
		field.String("email"),
		field.String("password"),
	}
}

func (Authentication) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("group_acess", GroupAcess.Type).
			Unique().
			Immutable().
			Required().
			Field("groupId"),
	}
}

func (Authentication) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "autentication",
		},
	}
}
