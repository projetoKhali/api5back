package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Department struct {
	ent.Schema
}

func (Department) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("description"),
	}
}

func (Department) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("access_group", AccessGroup.Type),
	}
}

func (Department) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "department",
		},
	}
}
