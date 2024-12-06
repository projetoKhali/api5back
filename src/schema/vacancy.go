package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Vacancy holds the schema definition for the Vacancy entity.
type Vacancy struct {
	ent.Schema
}

// Fields of the Vacancy.
func (Vacancy) Fields() []ent.Field {
	return []ent.Field{
		field.Int("proccessId"),
		field.String("title"),
		field.Int("positions"),
		field.Int("status").Default(1),
		field.String("location"),
		field.Int("userId"),
		field.Time("openingDate"),
		field.Int("closingDate"),
	}
}

// Edges of the Vacancy.
func (Vacancy) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type),
	}
}

func (Vacancy) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "vacancy",
		},
	}
}
