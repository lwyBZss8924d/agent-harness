package secretpolicy

type RevealPolicy string
type UsagePolicy string
type Action string
type Category string

const (
	RevealNever     RevealPolicy = "never"
	RevealAdminOnly RevealPolicy = "admin_only"
	RevealAllowed   RevealPolicy = "allowed"
)

const (
	UsageOpaqueUse UsagePolicy = "opaque_use"
	UsageInjectEnv UsagePolicy = "inject_env"
	UsageReveal    UsagePolicy = "reveal"
)

const (
	ActionLLMVerify     Action = "llm.verify"
	ActionLLMRequest    Action = "llm.request"
	ActionLLMModels     Action = "llm.models"
	ActionRegistryLogin Action = "registry.login"
	ActionWebhookSend   Action = "webhook.send"
	ActionK8sAuth       Action = "k8s.auth"
	ActionSignPayload   Action = "sign.payload"
)

const (
	CategoryLLMAPI     Category = "llm_api"
	CategoryCredential Category = "credential"
	CategoryToken      Category = "token"
	CategoryWebhook    Category = "webhook"
	CategorySigningKey Category = "signing_key"
)
