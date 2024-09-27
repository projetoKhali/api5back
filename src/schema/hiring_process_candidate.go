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

type HiringProcessCandidate struct {
	ent.Schema
}

func (HiringProcessCandidate) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("email").
			Unique(),
		field.String("phone").
			Unique(),
		field.Float("score"),
		field.Int("factHiringProcessId").
			Immutable().
			Unique(),
		field.Other("applyDate", &pgtype.Date{}).SchemaType(map[string]string{
			dialect.Postgres: "date",
		}),
		field.Enum("status").
			GoType(property.HiringProcessCandidateStatus(0)).
			Immutable(),
		field.Other("updatedAt", &pgtype.Date{}).SchemaType(map[string]string{
			dialect.Postgres: "date",
		}).Optional(),
	}
}

func (HiringProcessCandidate) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("factHiringProcess", FactHiringProcess.Type).
			Field("factHiringProcessId").
			Unique().
			Required().
			Immutable(),
	}
}

func (HiringProcessCandidate) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "hiring_process_candidate",
		},
	}
}
