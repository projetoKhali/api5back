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

type DimDatetime struct {
	ent.Schema
}

func (DimDatetime) Fields() []ent.Field {
	return []ent.Field{
		field.Other("date", &pgtype.Date{}).SchemaType(map[string]string{
			dialect.Postgres: "date",
		}),
		field.Int("year"),
		field.Int("month"),
		field.Int("weekday"),
		field.Int("day"),
		field.Int("hour"),
		field.Int("minute"),
		field.Int("second"),
	}
}

func (DimDatetime) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("fact_hiring_process", FactHiringProcess.Type).
			Ref("dimDatetime"),
	}
}

func (DimDatetime) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "dim_datetime",
		},
	}
}
