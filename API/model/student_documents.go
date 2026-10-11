package model

import "mysql/model/base"

type StudentDocument struct {
	base.ModelBase
	StudentID      int     `gorm:"not null;uniqueIndex:uk_student_document" json:"student_id"`
	DocumentTypeID int     `gorm:"not null;uniqueIndex:uk_student_document;index:fk_student_documents_document_type" json:"document_type_id"`
	RequiredQty    int     `gorm:"type:tinyint unsigned;not null;default:1" json:"required_qty"`
	ReceivedQty    int     `gorm:"type:tinyint unsigned;not null;default:0" json:"received_qty"`
	Remark         *string `gorm:"type:text;default:null" json:"remark"`
}

func (StudentDocument) TableName() string {
	return "student_documents"
}
