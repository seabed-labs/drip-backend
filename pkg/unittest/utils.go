package unittest

import "github.com/dcaf-labs/drip/pkg/service/configs"

func GetTestPrivateKey() string {
	return "[95,189,40,215,74,154,138,123,245,115,184,90,2,187,104,25,241,164,79,247,14,69,207,235,40,245,13,157,149,20,13,227,252,155,201,43,89,96,76,119,162,241,148,53,80,193,126,159,80,213,140,166,144,139,205,143,160,238,11,34,192,249,59,31]"
}

func GetTestAppConfig() *configs.AppConfig {
	return &configs.AppConfig{
		Network:                   configs.DevnetNetwork,
		Environment:               configs.StagingEnv,
		Wallet:                    GetTestPrivateKey(),
		DripProgramID:             "F1NyoZsUhJzcpGyoEqpDNbUMKVvCnSXcCki1nN3ycAeo",
		GoogleClientID:            "",
		Port:                      8080,
		DiscordWebhookID:          "",
		DiscordWebhookAccessToken: "",
		ShouldByPassAdminAuth:     false,
	}
}
