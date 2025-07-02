package config

import (
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/viper"
)

type Config struct {
	Name                         string `mapstructure:"NAME"`
	Environment                  string `mapstructure:"ENVIRONMENT"`
	Port                         int    `mapstructure:"PORT"`
	Version                      string `mapstructure:"VERSION"`
	AppUrl                       string `mapstructure:"APP_URL"`
	AccessTokenSigningKey        string `mapstructure:"ACCESS_TOKEN_SIGNING_KEY"`
	AccessTokenTokenExpiration   int    `mapstructure:"ACCESS_TOKEN_EXPIRATION"`
	RefreshTokenSigningKey       string `mapstructure:"REFRESH_TOKEN_SIGNING_KEY"`
	RefreshTokenExpiration       int    `mapstructure:"REFRESH_TOKEN_EXPIRATION"`
	RabbitMQUserName             string `mapstructure:"RABBITMQ_USERNAME"`
	RabbitMQPassword             string `mapstructure:"RABBITMQ_PASSWORD"`
	RabbitMQHost                 string `mapstructure:"RABBITMQ_HOST"`
	RabbitMQPort                 string `mapstructure:"RABBITMQ_PORT"`
	DatabasePort                 int    `mapstructure:"DATABASE_PORT"`
	DatabaseName                 string `mapstructure:"DATABASE_NAME"`
	DatabaseHost                 string `mapstructure:"DATABASE_HOST"`
	DatabaseUsername             string `mapstructure:"DATABASE_USERNAME"`
	DatabasePassword             string `mapstructure:"DATABASE_PASSWORD"`
	DatabaseSslMode              string `mapstructure:"DATABASE_SSL_MODE"`
	SMTPHost                     string `mapstructure:"SMTP_HOST"`
	SMTPPort                     int    `mapstructure:"SMTP_PORT"`
	StaticDirPath                string `mapstructure:"STATIC_DIR_PATH"`
	SMTPUsername                 string `mapstructure:"SMTP_USERNAME"`
	SMTPPassword                 string `mapstructure:"SMTP_PASSWORD"`
	SMTPFrom                     string `mapstructure:"SMTP_FROM"`
	LogEnable                    bool   `mapstructure:"LOG_ENABLE"`
	MailEnable                   bool   `mapstructure:"MAIL_ENABLE"`
	PayPingCallbackTargetWebsite string `mapstructure:"PAYPING_CALLBACK_TARGET_WEBSITE"`
	RequiredSignUpInviteCode     bool   `mapstructure:"REQUIRED_SIGN_UP_INVITE_CODE"`
	KavenegarApiKey              string `mapstructure:"KAVENEGAR_API_KEY"`
	KavenegarSender              string `mapstructure:"KAVENEGAR_SENDER"`
	MaxLoginDeviceCount          int    `mapstructure:"MAX_LOGIN_DEVICE_COUNT"`
	AutoDeleteDevice             bool   `mapstructure:"AUTO_DELETE_DEVICE"`
	SendSMS                      bool   `mapstructure:"SEND_SMS"`
	SendEmail                    bool   `mapstructure:"SEND_EMAIL"`
	PayPingToken                 string `mapstructure:"PAYPING_TOKEN"`
	PayPingCallback              string `mapstructure:"PAYPING_CALLBACK"`
	AwsEndpoint                  string `mapstructure:"AWS_S3_HOST"`
	AwsBucketName                string `mapstructure:"AWS_S3_PRIVATE_BUCKET_NAME"`
	AwsAccessKey                 string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AwsSecretKey                 string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	AwsRegion                    string `mapstructure:"AWS_REGION"`
	// Change IPServer
	IPServer string `mapstructure:"IP_ADDRESS_SERVER"`

	RoutePrefix string
}

var AppConfig *Config = &Config{}

func Load() {
	viper.AddConfigPath(".")
	viper.SetConfigType("env")
	fmt.Println("ENVIRONMENT: ", os.Getenv("ENVIRONMENT"))
	if os.Getenv("ENVIRONMENT") != "" {
		env := fmt.Sprintf(".env.%s", os.Getenv("ENVIRONMENT"))
		viper.SetConfigName(env)
	} else {
		viper.SetConfigName(".env")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// read config from system environment
	elem := reflect.TypeOf(AppConfig).Elem()
	for i := 0; i < elem.NumField(); i++ {
		key := elem.Field(i).Tag.Get("mapstructure")
		value := os.Getenv(key)
		if value != "" {
			viper.Set(key, value)
		}
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	err := viper.Unmarshal(&AppConfig)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Config unmarshal error: ", err)
	}
}
