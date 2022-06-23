package model

type DynamicGroupsFields struct {
	DynamicFieldID int32 `gorm:"column:dynamic_field_id;type:INT4;" json:"dynamic_field_id"`
	DynamicGroupID int64 `gorm:"column:dynamic_group_id;type:INT4;" json:"dynamic_group_id"`
}

func (t *DynamicGroupsFields) TableName() string {
	return "dynamic_groups_fields"
}
