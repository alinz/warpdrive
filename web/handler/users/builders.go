package users

type createUserRequest struct {
	Name     *string `json:"name,required"`
	Email    *string `json:"email,required"`
	Password *string `json:"password,required"`
}

func createUserRequestBuilder() interface{} {
	return &createUserRequest{}
}

type updateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func updateUserRequestBuilder() interface{} {
	return &updateUserRequest{}
}
