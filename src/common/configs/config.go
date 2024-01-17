package configs

type Config struct {
	Rsa     *Rsa      `mapstructure:"rsa"`
	PgpKeys []*PgpKey `mapstructure:"pgp"`
	Sftp    *Sftp     `mapstructure:"sftp"`
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
