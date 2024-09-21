package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type DimProcess struct {
	ent.Schema
}

func (DimProcess) Fields() []ent.Field {
	return []ent.Field{
		field.String("title"),
		field.Time("initialDate"),
		field.Time("finishDate"),
		field.Int("status").Default(1),
		field.Int("dimUsrId"),
		field.String("description").Optional(),
	}
}

func (DimProcess) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("fact_hiring_process", FactHiringProcess.Type).
			Ref("dimProcess"),
	}
}

func (DimProcess) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "dim_process",
		},
	}
}
