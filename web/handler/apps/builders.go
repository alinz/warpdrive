package apps

type createAppRequest struct {
	Name *string `json:"name,required"`
}

func createAppRequestBuilder() interface{} {
	return &createAppRequest{}
}

type createCycleRequest struct {
	Name *string `json:"name,required"`
}

func createCycleRequestBuilder() interface{} {
	return &createCycleRequest{}
}

type updateCycleRequest struct {
	Name *string `json:"name,required"`
}

func updateCycleRequestBuilder() interface{} {
	return &updateCycleRequest{}
}
