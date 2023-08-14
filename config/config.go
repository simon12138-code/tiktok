/*
* @Author: zgy
* @Date:   2023/7/25 15:08
 */
package config

// 一定要注意 settings-dev.yaml每个字段名称和结构体的tag一一对应
type MysqlConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbName"`
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ServerConfig struct {
	Name       string      `mapstructure:"name"`
	Port       int         `mapstructure:"port"`
	Mysqlinfo  MysqlConfig `mapstructure:"mysql"`
	Redisinfo  RedisConfig `mapstructure:"redis"`
	LogAddress string      `mapstructure:"logsAddress"`
	JWTKey     JWTconfig   `mapstructure:"jwt"`
	MinioInfo  MinioConfig `mapstructure:"minio"`
}
type JWTconfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type MinioConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyId     string `mapstructure:"accessKeyID"`
	SecretAccessKey string `mapstructure:"secretAccessKey"`
}
