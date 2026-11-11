package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// TenantMixin 租户唯一标识
type TenantMixin struct {
	mixin.Schema
}

func (TenantMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("tenant_id").
			Positive().
			Immutable().
			StructTag(`json:"tenantId"`).
			Comment("租户id"),
	}
}

// CreateAuditMixin 创建审计
type CreateAuditMixin struct {
	mixin.Schema
}

func (CreateAuditMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			SchemaType(map[string]string{dialect.MySQL: "datetime(3)"}).
			Annotations(entsql.DefaultExprs(map[string]string{dialect.MySQL: "CURRENT_TIMESTAMP(3)"})).
			StructTag(`json:"createdAt"`).
			Comment("创建时间"),
		field.String("created_by").
			MaxLen(64).
			Optional().
			Nillable().
			Immutable().
			StructTag(`json:"createdBy"`).
			Comment("创建人/主体"),
	}
}

// MutableAuditMixin 修改审计
type MutableAuditMixin struct {
	mixin.Schema
}

func (MutableAuditMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			SchemaType(map[string]string{dialect.MySQL: "datetime(3)"}).
			Annotations(entsql.DefaultExprs(map[string]string{dialect.MySQL: "CURRENT_TIMESTAMP(3)"})).
			StructTag(`json:"updatedAt"`).
			Comment("最后更新时间"),
		field.String("updated_by").
			MaxLen(64).
			Optional().
			Nillable().
			StructTag(`json:"updatedBy"`).
			Comment("最后更新人/主体"),
		field.Time("deleted_at").
			Optional().
			Nillable().
			SchemaType(map[string]string{dialect.MySQL: "datetime(3)"}).
			StructTag(`json:"deletedAt"`).
			Comment("逻辑删除时间,NULL表示未删除"),
		field.String("deleted_by").
			MaxLen(64).
			Optional().
			Nillable().
			StructTag(`json:"deletedBy"`).
			Comment("删除人/主体"),
	}
}
