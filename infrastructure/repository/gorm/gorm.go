package gormRepo

import "gorm.io/gorm"

func Init(db *gorm.DB) {
	db.AutoMigrate(&gormBeneficiary{}, &gormTransaction{})
}
