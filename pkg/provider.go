package configfs

type ConfigEntry struct {
	Name    string
	Project string
}

type Provider interface {
	List() ([]ConfigEntry, error)
	Value(k, project, env string) (string, error)
}
