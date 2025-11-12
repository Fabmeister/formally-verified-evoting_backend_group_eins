package configuration

type DatabaseConfig struct {
	Address  string
	Database string
	Port     string
	User     string
	Pwd      string
}

func (this *DatabaseConfig) SetStandard() {
	this.Address = "localhost"
	this.Database = "evoting"
	this.Port = "3306"
	this.User = "user"
	this.Pwd = "userpwd"
}
