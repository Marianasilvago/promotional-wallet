package config

type ClientConfig struct {
	timeout            int
	stockClientBaseURL string
}

func newClientConfig() ClientConfig {
	return ClientConfig{
		timeout:            getInt("HTTP_CLIENT_TIMEOUT_IN_SEC", 10),
		stockClientBaseURL: getString("STOCK_CLIENT_BASE_URL"),
	}
}

func (cc ClientConfig) GetTimeout() int {
	return cc.timeout
}

func (cc ClientConfig) GetStockClientBaseURL() string {
	return cc.stockClientBaseURL
}
