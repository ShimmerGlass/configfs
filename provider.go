package main

type Provider interface {
	List() ([]string, error)
	Value(k, project, env string) (string, error)
}
