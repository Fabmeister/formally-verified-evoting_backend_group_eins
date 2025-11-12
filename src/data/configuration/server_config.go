package configuration

type ServerConfig struct {
	Port int32
}

func (this *ServerConfig) SetStandard() {
	this.Port = 50051
}
