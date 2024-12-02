package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type DimDepartment struct {
	ent.Schema
}

func (DimDepartment) Fields() []ent.Field {
	return []ent.Field{
		field.Int("dbId"),
		field.String("name"),
		field.String("description"),
	}
}

func (DimDepartment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("dim_process", DimProcess.Type).
			Ref("dimDepartment"),
	}
}

func (DimDepartment) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "dim_department",
		},
	}
}
