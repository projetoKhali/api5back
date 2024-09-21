package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type DimVacancy struct {
	ent.Schema
}

func (DimVacancy) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.String("title"),
		field.Int("numPositions"),
		field.Int("reqId"),
		field.Int("status").Default(1),
		field.String("location"),
		field.Int("dimUsrId"),
		field.Time("openingDate"),
		field.Int("closingDate"),
	}
}

func (DimVacancy) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("fact_hiring_process", FactHiringProcess.Type).
			Ref("dimVacancy"),
	}
}

func (DimVacancy) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "dim_vacancy",
		},
	}
}
