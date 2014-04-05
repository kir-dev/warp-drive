package main

type environment string

func (e environment) dev() bool {
	return e == "development" || e == "dev"
}

func (e environment) test() bool {
	return e == "test"
}

func (e environment) prod() bool {
	return e == "prod" || e == "production"
}
