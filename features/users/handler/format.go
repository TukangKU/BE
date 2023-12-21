package user

type RegisterRequest struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type RegisterResponse struct {
	ID       uint   `json:"id"`
	Role     string `json:"role"`
	UserName string `json:"username"`
	Email    string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID       uint   `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Token    string `json:"token"`
}

type UserUpdate struct {
	ID       uint   `json:"id" form:"id"`
	UserName string `json:"username" form:"username"`
	Nama     string `json:"nama" form:"nama"`
	Email    string `json:"email" form:"email"`
	NoHp     string `json:"nohp" form:"nohp"`
	Alamat   string `json:"alamat" form:"alamat"`
	Foto     string `json:"foto" form:"foto"`
	Role     string `json:"role" form:"role"`
	Skills   []uint `json:"skill" form:"skill"`
}

type UserResponse struct {
	ID       uint      `json:"id" form:"id"`
	UserName string    `json:"username" form:"username"`
	Nama     string    `json:"nama" form:"nama"`
	Email    string    `json:"email" form:"email"`
	NoHp     string    `json:"nohp" form:"nohp"`
	Alamat   string    `json:"alamat" form:"alamat"`
	Foto     string    `json:"foto" form:"foto"`
	Role     string    `json:"role" form:"role"`
	Skill    []string  `json:"skill" form:"skill"`
	Skilll UserSkill
}

type UserResponseUpdate struct {
	ID       uint      `json:"id" form:"id"`
	UserName string    `json:"username" form:"username"`
	Nama     string    `json:"nama" form:"nama"`
	Email    string    `json:"email" form:"email"`
	NoHp     string    `json:"nohp" form:"nohp"`
	Alamat   string    `json:"alamat" form:"alamat"`
	Foto     string    `json:"foto" form:"foto"`
	Role     string    `json:"role" form:"role"`
	Skill    []UserSkill  `json:"skill" form:"skill"`
}

type GetUserResponse struct {
	ID       uint      `json:"id" form:"id"`
	UserName string    `json:"username" form:"username"`
	Nama     string    `json:"nama" form:"nama"`
	Email    string    `json:"email" form:"email"`
	Alamat   string    `json:"alamat" form:"alamat"`
	Foto     string    `json:"foto" form:"foto"`
	Skill    []UserSkill  `json:"skill" form:"skill"`
}

type UserSkill struct {
	SkillID   uint   `json:"skill_id"`
	NamaSKill string `json:"skill"`
}


