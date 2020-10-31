package views

//RegisterBody model
type RegisterBody struct {
	Username string `form:"username" binding:"required,username"`
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

// TODO Add password validation

//LoginBody struct
type LoginBody struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}
