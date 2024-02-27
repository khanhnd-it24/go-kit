package configs

type Config struct {
	Rsa      *Rsa      `mapstructure:"rsa"`
	PgpKeys  []*PgpKey `mapstructure:"pgp"`
	Sftp     *Sftp     `mapstructure:"sftp"`
	Postgres *Postgres `mapstructure:"postgres"`

	Mongo *Mongo `mapstructure:"mongo"`
}

type Postgres struct {
	Host        string `mapstructure:"host"`
	Port        string `mapstructure:"port"`
	User        string `mapstructure:"user"`
	DbName      string `mapstructure:"db_name"`
	SslMode     string `mapstructure:"ssl_mode"`
	Password    string `mapstructure:"password"`
	AutoMigrate bool   `mapstructure:"auto_migrate"`
	MaxLifeTime int    `mapstructure:"max_life_time"`
}

type Mongo struct {
	Uri string `mapstructure:"uri"`
	DB  string `mapstructure:"db"`
}

type Rsa struct {
	PublicKeys  map[string]string `mapstructure:"public_keys"`
	PrivateKeys map[string]string `mapstructure:"private_keys"`
}

type PgpKey struct {
	Name       string `mapstructure:"name"`
	Path       string `mapstructure:"path"`
	Passphrase string `mapstructure:"passphrase"`
}

type Sftp struct {
	User string `mapstructure:"user"`
	Pass string `mapstructure:"pass"`
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}
