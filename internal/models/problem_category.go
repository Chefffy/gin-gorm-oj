package models

import "gorm.io/gorm"

type ProblemCategory struct {
	gorm.Model
	ProblemId	uint               `gorm:"column:problem_id;type:int;"json:"problem_id"`
	CategoryId	uint              `goem:"column:category_id;type:int;"json:"category_id"`
	CategoryBasic	*CategoryBasic `gorm:"foreignKey:id;references:category_id"` //关联分类的基础信息表
}

func (table *ProblemCategory)TableName() string{
	return "problem_category"
}