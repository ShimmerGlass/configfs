package main

type Provider interface {
	List() ([]string, error)
	Value(k, env string) (string, error)
}
