package mango

func Redirect(status Status, location string) (Status, Headers, Body) {
	return status, Headers{"Location": []string{location}}, Body("")
}
