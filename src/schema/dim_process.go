package schema

import (
	"api5back/src/property"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/jackc/pgx/v5/pgtype"
)

type DimProcess struct {
	ent.Schema
}

func (DimProcess) Fields() []ent.Field {
	return []ent.Field{
		field.Int("dbId"),
		field.String("title"),
		field.Other("initialDate", &pgtype.Date{}).SchemaType(map[string]string{
			dialect.Postgres: "date",
		}),
		field.Other("finishDate", &pgtype.Date{}).SchemaType(map[string]string{
			dialect.Postgres: "date",
		}).Optional(),
		field.Int("status").
			GoType(property.DimProcessStatus(1)).
			SchemaType(map[string]string{
				dialect.Postgres: "character varying",
			}).
			Default(int(property.DimProcessStatusOpen)),
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
