package request

type ShiftRequestCreate struct {
	Name     string `gorm:"column:name;size:255;not null" json:"name"`
	Session1 string `gorm:"column:session1;not null;default:false" json:"session1"`
	Session2 string `gorm:"column:session2;not null;default:false" json:"session2"`
	Session3 string `gorm:"column:session3;not null;default:false" json:"session3"`
	Session4 string `gorm:"column:session4;not null;default:false" json:"session4"`
	Session5 string `gorm:"column:session5;not null;default:false" json:"session5"`
}

type ShiftRequestUpdate struct {
	Name     string `gorm:"column:name;size:255;not null" json:"name"`
	Session1 string `gorm:"column:session1;not null;default:false" json:"session1"`
	Session2 string `gorm:"column:session2;not null;default:false" json:"session2"`
	Session3 string `gorm:"column:session3;not null;default:false" json:"session3"`
	Session4 string `gorm:"column:session4;not null;default:false" json:"session4"`
	Session5 string `gorm:"column:session5;not null;default:false" json:"session5"`
}
