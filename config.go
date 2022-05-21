package seras

type Config struct {
	Mods []string
}

func (c *Config) Connections() []Connection {
	return nil
}
