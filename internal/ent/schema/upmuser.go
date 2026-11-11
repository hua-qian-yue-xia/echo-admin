package schema

import "entgo.io/ent"

// UpmUser holds the schema definition for the UpmUser entity.
type UpmUser struct {
	ent.Schema
}

// Fields of the UpmUser.
func (UpmUser) Fields() []ent.Field {
	return nil
}

// Edges of the UpmUser.
func (UpmUser) Edges() []ent.Edge {
	return nil
}
