package common

// AppConfig is the configuration that is loaded in one way or another, and then provided to the
// different services and controllers
type AppConfig struct {
	AuthPublicKey  string
	AuthPrivateKey string
}
