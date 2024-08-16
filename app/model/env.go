package model

type AppConfig struct {
	FileName        string      `json:"-" yaml:"-" gorm:"primaryKey;uniqueIndex;not null"`
	MaxCPUCount     int         `json:"maxCPUCount" yaml:"maxCPUCount" gorm:"not null"`
	MaxMemoryGB     int         `json:"maxMemoryGB" yaml:"maxMemoryGB" gorm:"not null"`
	Description     string      `json:"description" yaml:"description" gorm:"not null"`
	JavaPreCmd      DataStrings `json:"javaPreCmd" yaml:"javaPreCmd" gorm:"not null"`
	JavaPostCmd     DataStrings `json:"javaPostCmd" yaml:"javaPostCmd" gorm:"not null"`
	NodeSelector    DatabaseMap `json:"nodeSelector" yaml:"nodeSelector" gorm:"not null"`
	OverrideNode    bool        `json:"overrideNode" yaml:"overrideNode" gorm:"not null"`
	GatewayName     string      `json:"gatewayName" yaml:"gatewayName" gorm:"not null"`
	GatewayNodePort int         `json:"gatewayNodePort" yaml:"gatewayNodePort" gorm:"not null"`
	HostPort        int         `json:"hostPort" yaml:"hostPort" gorm:"not null"`
	Replicas        int         `json:"replicas" yaml:"replicas" gorm:"not null"`
}
