package model

type AppBaseConfig struct {
	MaxCPUCount  int         `json:"maxCPUCount" yaml:"maxCPUCount" gorm:"not null"`
	MaxMemoryGB  int         `json:"maxMemoryGB" yaml:"maxMemoryGB" gorm:"not null"`
	Description  string      `json:"description" yaml:"description" gorm:"not null"`
	PreCmd       DataStrings `json:"preCmd" yaml:"preCmd" gorm:"not null"`
	Args         DataStrings `json:"args" yaml:"args" gorm:"not null"`
	PostCmd      DataStrings `json:"postCmd" yaml:"postCmd" gorm:"not null"`
	NodeSelector DatabaseMap `json:"nodeSelector" yaml:"nodeSelector" gorm:"not null"`
	Replicas     int         `json:"replicas" yaml:"replicas" gorm:"not null"`
}

type AppConfig struct {
	Name          string `json:"name" yaml:"name" gorm:"primaryKey;uniqueIndex;not null"`
	AppBaseConfig `json:",inline" yaml:",inline"  gorm:"embedded"`
}

type AppDeploy struct {
	AppName string `json:"appName" yaml:"appName" gorm:"not null;uniqueIndex:uidx_app_name_env_name"`
	EnvName string `json:"envName" yaml:"envName" gorm:"not null;uniqueIndex:uidx_app_name_env_name"`
	Image   string `json:"image" yaml:"image" gorm:"not null"`
	Tag     string `json:"tag" yaml:"tag" gorm:"not null"`
}

type EnvConfig struct {
	AppBaseConfig   `json:",inline" yaml:",inline" gorm:"embedded"`
	FileName        string      `json:"-" yaml:"-" gorm:"primaryKey;uniqueIndex;not null"`
	Namespace       string      `json:"namespace" yaml:"namespace" gorm:"not null"`
	OverrideNode    bool        `json:"overrideNode" yaml:"overrideNode" gorm:"not null"`
	GatewayName     string      `json:"gatewayName" yaml:"gatewayName" gorm:"not null"`
	GatewayNodePort int         `json:"gatewayNodePort" yaml:"gatewayNodePort" gorm:"not null"`
	HostPort        int         `json:"hostPort" yaml:"hostPort" gorm:"not null"`
	EnvValues       DatabaseMap `json:"envValues" yaml:"envValues" gorm:"not null"`
	UserGuide       DatabaseMap `json:"-" yaml:"-" gorm:"not null"`
}
