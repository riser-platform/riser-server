package core

// RuntimeConfig provides config for the srver.
type RuntimeConfig struct {
	BootstrapApikey          string `split_words:"true"`
	BindAddress              string `split_words:"true" default:":8000"`
	DeveloperMode            bool   `split_words:"true"`
	GitUrl                   string `split_words:"true" required:"true"`
	GitDir                   string `split_words:"true" default:"/tmp/riser/git/"`
	GitBranch                string `split_words:"true" default:"master"`
	GitUsername              string `split_words:"true"`
	GitPassword              string `split_words:"true"`
	PostgresUrl              string `split_words:"true" default:"postgres://postgres.riser-system.svc.cluster.local/riserdb?sslmode=disable&connect_timeout=3"`
	PostgresUsername         string `split_words:"true" required:"true"`
	PostgresPassword         string `split_words:"true" required:"true"`
	PostgresMigrateOnStartup bool   `split_words:"true" default:"true"`
}
