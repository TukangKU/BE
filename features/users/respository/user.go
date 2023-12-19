package model

import (
	"errors"
	"tukangku/features/skill"
	"tukangku/features/skill/repository"
	"tukangku/features/users"

	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	Nama     string `json:"nama"`
	UserName string `json:"username" gorm:"unique"`
	Password string `json:"password"`
	Email    string `json:"email" gorm:"unique"`
	NoHp     string `json:"nohp"`
	Alamat   string `json:"alamat"`
	Foto     string `json:"foto"`
	Role     string `json:"role"`
	Skills   string
	Skill    []repository.SkillModel `gorm:"many2many:user_skills;"`
}

type userQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) users.Repository {
	return &userQuery{
		db: db,
	}
}

func (ur *userQuery) Register(newUser users.Users) (users.Users, error) {
	var input = new(UserModel)
	input.UserName = newUser.UserName
	input.Email = newUser.Email
	input.Password = newUser.Password
	input.Role = newUser.Role
	input.Foto = "https://res.cloudinary.com/daxpcsncf/image/upload/v1702888962/mfgsrgdlsguqjoskujib.png"

	if err := ur.db.Create(&input).Error; err != nil {
		return users.Users{}, err
	}

	newUser.ID = input.ID

	return newUser, nil

}

func (ul *userQuery) Login(email string) (users.Users, error) {
	var userData = new(UserModel)

	if err := ul.db.Where("email = ?", email).First(userData).Error; err != nil {
		return users.Users{}, nil
	}

	var result = new(users.Users)
	result.ID = userData.ID
	result.UserName = userData.UserName
	result.Password = userData.Password
	result.Email = userData.Email
	result.Role = userData.Role

	return *result, nil
}

func (us *userQuery) UpdateUser(idUser uint, updateWorker users.Users) (users.Users, error) {
	var exitingUser = new(UserModel)
	exitingUser.UserName = updateWorker.UserName
	exitingUser.Nama = updateWorker.Nama
	exitingUser.Email = updateWorker.Email
	exitingUser.NoHp = updateWorker.NoHp
	exitingUser.Alamat = updateWorker.Alamat
	exitingUser.Foto = updateWorker.Foto

	if err := us.db.Where("id = ?", idUser).Updates(exitingUser).Error; err != nil {
		return users.Users{}, err
	}

	if updateWorker.ID != 0 {
		exitingUser.ID = updateWorker.ID
	}

	if updateWorker.Nama != "" {
		exitingUser.Nama = updateWorker.Nama
	}

	if updateWorker.Email != "" {
		exitingUser.Email = updateWorker.Email
	}

	if updateWorker.NoHp != "" {
		exitingUser.NoHp = updateWorker.NoHp
	}

	if updateWorker.Alamat != "" {
		exitingUser.Alamat = updateWorker.Alamat
	}

	if len(updateWorker.Skill) != 0 {
		var userSkill = []repository.SkillModel{}
		for _, v := range updateWorker.Skill {
			userSkill = append(userSkill, repository.SkillModel{
				ID:        v.ID,
				NamaSkill: v.NamaSkill,
			})

		}

		if err := us.db.Model(exitingUser).Association("Skill").Replace(&userSkill); err != nil {

			return users.Users{}, err
		}

	}

	var user UserModel

	if err := us.db.Preload("Skill").Where("id = ?", idUser).First(&user).Error; err != nil {
		return users.Users{}, err
	}

	var response []string
	for _, v := range user.Skill {
		response = append(response, v.NamaSkill)
	}

	result := users.Users{
		ID:       idUser,
		Nama:     exitingUser.Nama,
		Email:    exitingUser.Email,
		NoHp:     exitingUser.NoHp,
		Alamat:   exitingUser.Alamat,
		Foto:     exitingUser.Foto,
		UserName: exitingUser.UserName,
		Skills:   response,
	}
	return result, nil

}

func (gu *userQuery) GetUserByID(idUser uint) (users.Users, error) {
	var result UserModel

	if err := gu.db.Preload("Skill").Where("id = ?", idUser).Find(&result).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return users.Users{}, errors.New("user not found")
		}
		return users.Users{}, err
	}

	response := users.Users{
		ID:       result.ID,
		Nama:     result.Nama,
		Email:    result.Email,
		NoHp:     result.NoHp,
		Alamat:   result.Alamat,
		Foto:     result.Foto,
		UserName: result.UserName,
		Role:     result.Role,
	}
	for _, v := range result.Skill {
		response.Skill = append(response.Skill, skill.Skills{
			ID:        v.ID,
			NamaSkill: v.NamaSkill,
		})
	}

	return response, nil
}
