package config

type SecretConfigInfo struct {
	SecretFileEnvName string
	SecretFileName    string
	SecretValuePtr    *string
}

type IConfig interface {
	GetSecretInfos() []SecretConfigInfo
}
