package config

import (
	"github.com/cosmos/cosmos-sdk/server/config"
	"github.com/spf13/viper"

	ethermint "github.com/jiujiteiro/ethermint/types"
)

const (
	// DefaultGRPCAddress is the default address the gRPC server binds to.
	DefaultGRPCAddress = "0.0.0.0:9900"

	// DefaultEVMAddress is the default address the EVM JSON-RPC server binds to.
	DefaultEVMAddress = "0.0.0.0:8545"

	// DefaultEVMWSAddress is the default address the EVM WebSocket server binds to.
	DefaultEVMWSAddress = "0.0.0.0:8546"
)

// AppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func AppConfig() (string, interface{}) {
	// Optionally allow the chain developer to overwrite the SDK's default
	// server config.
	srvCfg := config.DefaultConfig()

	// The SDK's default minimum gas price is set to "" (empty value) inside
	// app.toml. If left empty by validators, the node will halt on startup.
	// However, the chain developer can set a default app.toml value for their
	// validators here.
	//
	// In summary:
	// - if you leave srvCfg.MinGasPrices = "", all validators MUST tweak their
	//   own app.toml config,
	// - if you set srvCfg.MinGasPrices non-empty, validators CAN tweak their
	//   own app.toml to override, or use this default value.
	//
	// In ethermint, we set the min gas prices to 0.
	srvCfg.MinGasPrices = "0" + ethermint.AttoPhoton

	customAppConfig := Config{
		Config: *srvCfg,
		EVMRPC: *DefaultEVMConfig(),
	}

	customAppTemplate := config.DefaultConfigTemplate + DefaultConfigTemplate

	return customAppTemplate, customAppConfig
}

// DefaultConfig returns server's default configuration.
func DefaultConfig() *Config {
	return &Config{
		Config: *config.DefaultConfig(),
		EVMRPC: *DefaultEVMConfig(),
	}
}

// DefaultEVMConfig returns an EVM config with the JSON-RPC API enabled by default
func DefaultEVMConfig() *EVMRPCConfig {
	return &EVMRPCConfig{
		Enable:     true,
		RPCAddress: DefaultEVMAddress,
		WsAddress:  DefaultEVMWSAddress,
	}
}

// EVMRPCConfig defines configuration for the EVM RPC server.
type EVMRPCConfig struct {
	// Enable defines if the EVM RPC server should be enabled.
	Enable bool `mapstructure:"enable"`
	// Address defines the HTTP server to listen on
	RPCAddress string `mapstructure:"address"`
	// Address defines the WebSocket server to listen on
	WsAddress string `mapstructure:"ws-address"`
}

// Config defines the server's top level configuration. It includes the default app config
// from the SDK as well as the EVM configuration to enable the JSON-RPC APIs.
type Config struct {
	config.Config

	EVMRPC EVMRPCConfig `mapstructure:"evm-rpc"`
}

// GetConfig returns a fully parsed Config object.
func GetConfig(v *viper.Viper) Config {

	cfg := config.GetConfig(v)

	return Config{
		Config: cfg,
		EVMRPC: EVMRPCConfig{
			Enable:     v.GetBool("evm-rpc.enable"),
			RPCAddress: v.GetString("evm-rpc.address"),
			WsAddress:  v.GetString("evm-rpc.ws-address"),
		},
	}
}
