package model

import (
	"errors"
	"fmt"

	"tukangku/features/skill"
	"tukangku/features/skill/repository"
	"tukangku/features/users"

	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	Nama     string
	UserName string `gorm:"unique"`
	Password string
	Email    string `json:"email" gorm:"unique"`
	NoHp     string
	Alamat   string
	Foto     string
	Role     string
	Skill    []repository.SkillModel `gorm:"many2many:user_skills;"`
	Jobs     []JobModel              `gorm:"foreignKey:WorkerID"`
	Requests []JobModel              `gorm:"foreignKey:ClientID"`
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

	var response []skill.Skills
	for _, v := range user.Skill {
		response = append(response, skill.Skills{
			ID:        v.ID,
			NamaSkill: v.NamaSkill,
		})
	}

	result := users.Users{
		ID:       idUser,
		Nama:     exitingUser.Nama,
		Email:    exitingUser.Email,
		NoHp:     exitingUser.NoHp,
		Alamat:   exitingUser.Alamat,
		Foto:     exitingUser.Foto,
		UserName: exitingUser.UserName,
		// Skills:   response,
		Skill: response,
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

func (gu *userQuery) GetUserBySKill(idSkill uint, page, pageSize int) ([]users.Users, int, error) {
	var result []UserModel
	var totalCount int64

	offset := (page - 1) * pageSize

	if err := gu.db.Model(&UserModel{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	if err := gu.db.Preload("Skill").
		Where("skill_models.id = ?", idSkill).
		Joins("JOIN user_skills ON user_models.id = user_skills.user_model_id").
		Joins("JOIN skill_models ON user_skills.skill_model_id = skill_models.id").
		Offset(offset).
		Limit(pageSize).
		Find(&result).Error; err != nil {
		return []users.Users{}, 0, err
	}
	// if err := gu.db.
	// 	Where("user_models.id IN (SELECT distinct(user_skills.user_model_id) FROM user_skills WHERE user_skills.skill_model_id = ?)", idSkill).
	// 	Offset(offset).
	// 	Limit(pageSize).
	// 	Find(&result).Error; err != nil {
	// 	return []users.Users{}, 0, err
	// }
	fmt.Println(result)
	var response []users.Users
	for _, v := range result {
		tmp := new(users.Users)
		tmp.ID = v.ID
		tmp.Nama = v.Nama
		tmp.UserName = v.UserName
		tmp.Alamat = v.Alamat
		tmp.Foto = v.Foto
		for _, v := range v.Skill {
			tmp.Skill = append(tmp.Skill, skill.Skills{
				ID:        v.ID,
				NamaSkill: v.NamaSkill,
			})
		}
		response = append(response, *tmp)
	}
	return response, int(totalCount), nil
}

func (tk *userQuery) TakeWorker(idUser uint) (users.Users, error) {
	var result UserModel

	if err := tk.db.Preload("Skill").Where("id = ?", idUser).Find(&result).Error; err != nil {
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
