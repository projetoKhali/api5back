package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type FactHiringProcess struct {
	ent.Schema
}

func (FactHiringProcess) Fields() []ent.Field {
	return []ent.Field{
		field.Int("dimProcessId").Immutable(),
		field.Int("dimVacancyId").Immutable(),
		field.Int("dimUserId").Immutable(),
		field.Int("dimDateId").Immutable(),
		field.Int("metTotalCandidatesApplied"),
		field.Int("metTotalCandidatesInterviewed"),
		field.Int("metTotalCandidatesHired"),
		field.Int("metSumDurationHiringProces"),
		field.Int("metSumSalaryInitial"),
		field.Int("metTotalFeedbackPositive"),
		field.Int("metTotalNeutral"),
		field.Int("metTotalNegative"),
	}
}

func (FactHiringProcess) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("dimProcess", DimProcess.Type).
			Unique().
			Immutable().
			Required().
			Field("dimProcessId"),
		edge.To("dimVacancy", DimVacancy.Type).
			Unique().
			Immutable().
			Required().
			Field("dimVacancyId"),
		edge.To("dimUser", DimUser.Type).
			Unique().
			Immutable().
			Required().
			Field("dimUserId"),
		edge.To("dimDatetime", DimDatetime.Type).
			Unique().
			Immutable().
			Required().
			Field("dimDateId"),
	}
}

func (FactHiringProcess) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "fact_hiring_process",
		},
	}
}
