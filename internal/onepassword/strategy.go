package onepassword

const (
	DefaultAuthMode = "service_account"
	SDKModule       = "github.com/1password/onepassword-sdk-go"
	SDKVersion      = "v0.4.0"
)

type Strategy struct {
	AuthMode               string
	UsesConnectServer      bool
	PrimaryTokenEnv        string
	RecommendedIntegration string
}

func DefaultStrategy() Strategy {
	return Strategy{
		AuthMode:               DefaultAuthMode,
		UsesConnectServer:      false,
		PrimaryTokenEnv:        "OP_SERVICE_ACCOUNT_TOKEN",
		RecommendedIntegration: "official_go_sdk_with_service_account",
	}
}
