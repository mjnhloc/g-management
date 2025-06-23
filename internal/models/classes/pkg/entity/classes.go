package entity

import (
	"g-management/pkg/shared/gorm/model"
)

type Classes struct {
	ID          int     `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement" mapstructure:"id"`
	Name        string  `gorm:"column:name;type:varchar(50);not null" mapstructure:"name"`
	TrainerID   int     `gorm:"column:trainer_id;type:bigint;not null" mapstructure:"trainer_id"`
	Schedule    string  `gorm:"column:schedule;type:datetime;not null" mapstructure:"schedule"`
	Duration    int     `gorm:"column:duration;type:int(5) unsigned;not null" mapstructure:"duration"`
	MaxCapacity int     `gorm:"column:max_capacity;type:int(3) unsigned;not null" mapstructure:"max_capacity"`
	Description *string `gorm:"column:description;type:text" mapstructure:"description"`
	model.BaseModel
}

const ClassesIndexNameElasticSearch = "classes"

type ClassDocument struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Schedule    string  `json:"schedule"`
	Description *string `json:"description,omitempty"`
}

func (c *ClassDocument) IndexName() string {
	return ClassesIndexNameElasticSearch
}

func (c *ClassDocument) DocumentID() int {
	return c.ID
}
