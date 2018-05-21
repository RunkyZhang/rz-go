package configuration

type Config struct {
	Web web `json:"web"`
}

type web struct {
	Listen string `json:"listen"`
}

type redis struct {
	Sentinels []string `json:"sentinels"`
	Password  string   `json:"password"`
}
