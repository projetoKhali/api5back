package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/jackc/pgx/v5/pgtype"
)

type DimVacancy struct {
	ent.Schema
}

func (DimVacancy) Fields() []ent.Field {
	return []ent.Field{
		field.Int("dbId"),
		field.String("title"),
		field.Int("numPositions"),
		field.Int("reqId").Optional(),
		field.Int("status").Default(1),
		field.String("location"),
		field.Int("dimUsrId"),
		field.Other("openingDate", &pgtype.Date{}).SchemaType(map[string]string{
			dialect.Postgres: "date",
		}),
		field.Other("closingDate", &pgtype.Date{}).SchemaType(map[string]string{
			dialect.Postgres: "date",
		}).Optional(),
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
