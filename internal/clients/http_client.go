package clients

type HttpClient interface {
	Get(url string) (int, string, error)
	GetWithHeaders(url string, headers map[string]string) (int, string, error)
}
