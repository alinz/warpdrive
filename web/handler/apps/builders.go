package apps

type createAppRequest struct {
	Name *string `json:"name,required"`
}

func createAppRequestBuilder() interface{} {
	return &createAppRequest{}
}
