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

type DimUser struct {
	ent.Schema
}

func (DimUser) Fields() []ent.Field {
	return []ent.Field{
		field.Int("dbId"),
		field.String("name"),
		field.String("occupation"),
		field.Other("updated_at", &pgtype.Date{}).SchemaType(map[string]string{
			dialect.Postgres: "date",
		}).Optional(),
	}
}

func (DimUser) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("fact_hiring_process", FactHiringProcess.Type).
			Ref("dimUser"),
	}
}

func (DimUser) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "dim_user",
		},
	}
}
