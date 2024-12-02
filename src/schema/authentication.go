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
		field.String("name"),
		field.String("email").
			Unique(),
		field.String("password"),
		field.Int("groupId").
			Immutable(),
	}
}

func (Authentication) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("access_group", AccessGroup.Type).
			Unique().
			Immutable().
			Required().
			Field("groupId"),
	}
}

func (Authentication) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "authentication",
		},
	}
}
