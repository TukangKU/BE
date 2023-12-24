package repository

import (
	"time"
	"tukangku/features/skill"

	"gorm.io/gorm"
)

type SkillModel struct {
	ID        uint   `gorm:"primarykey"`
	NamaSkill string `json:"skill"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Jobs      []JobModel     `gorm:"foreignKey:Category"`
}
type JobModel struct {
	gorm.Model
	WorkerID  uint   `gorm:"not null"`
	ClientID  uint   `gorm:"not null"`
	Category  uint   `gorm:"not null"`
	StartDate string `gorm:"not null"`
	EndDate   string `gorm:"not null"`
	Price     int
	Deskripsi string
	Status    string
	Address   string
	NoteNego  string
}

type SkillQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) skill.Repository {
	return &SkillQuery{
		db: db,
	}
}

// AddSkill implements skill.Repository.
func (as *SkillQuery) AddSkill(newSkill skill.Skills) (skill.Skills, error) {
	var inputData = new(SkillModel)
	inputData.NamaSkill = newSkill.NamaSkill

	if err := as.db.Create(&inputData).Error; err != nil {
		return skill.Skills{}, err
	}

	newSkill.ID = inputData.ID

	return newSkill, nil

}

// ShowSkill implements skill.Repository.
func (sq *SkillQuery) ShowSkill() ([]skill.Skills, error) {
	var skills []SkillModel
	if err := sq.db.Find(&skills).Error; err != nil {
		return nil, err
	}

	// Transform the database model to the desired Skills type.
	var result []skill.Skills
	for _, s := range skills {
		result = append(result, skill.Skills{
			ID:        s.ID,
			NamaSkill: s.NamaSkill,
			// Add other fields if needed
		})
	}

	return result, nil

}
