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

type DimCandidate struct {
	ent.Schema
}

func (DimCandidate) Fields() []ent.Field {
	return []ent.Field{
		field.Int("dbId"),
		field.String("name"),
		field.String("email"),
		field.String("phone"),
		field.Float("score"),
		field.Int("dimVacancyDbId").
			Immutable(),
		field.Other("applyDate", &pgtype.Date{}).SchemaType(map[string]string{
			dialect.Postgres: "date",
		}),
		field.Enum("status").
			GoType(property.DimCandidateStatus(0)).
			SchemaType(map[string]string{
				dialect.Postgres: "character varying",
			}).
			Immutable(),
		field.Other("updatedAt", &pgtype.Date{}).SchemaType(map[string]string{
			dialect.Postgres: "date",
		}).Optional(),
	}
}

func (DimCandidate) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("dimVacancy", DimVacancy.Type).
			Field("dimVacancyDbId").
			Unique().
			Required().
			Immutable(),
	}
}

func (DimCandidate) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "dim_candidate",
		},
	}
}
