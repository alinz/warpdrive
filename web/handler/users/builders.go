package users

type createUserRequest struct {
	Name     *string `json:"name,required"`
	Email    *string `json:"email,required"`
	Password *string `json:"password,required"`
}

func createUserRequestBuilder() interface{} {
	return &createUserRequest{}
}
