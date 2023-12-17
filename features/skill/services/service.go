package services

import (
	"tukangku/features/skill"
)

type SKillService struct {
	m skill.Repository
}

func New(model skill.Repository) skill.Service {
	return &SKillService{
		m: model,
	}
}

func (ss *SKillService) AddSkill(newSkill skill.Skills) (skill.Skills, error) {

	result, err := ss.m.AddSkill(newSkill)
	if err != nil {
		return skill.Skills{}, err
	}

	return result, nil
}

func (ss *SKillService) ShowSkill() ([]skill.Skills, error) {
	skills, err := ss.m.ShowSkill()
	if err != nil {
		return nil, err
	}

	return skills, nil
}
