package apps

type appRequest struct {
	Name string `json:"name"`
}

func appRequestBuilder() interface{} {
	return appRequest{}
}
