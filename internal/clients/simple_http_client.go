package clients

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type SimpleHttpClient struct {
}

func (c SimpleHttpClient) Get(url string) (int, string, error) {
	return c.GetWithHeaders(url, nil)
}

func (c SimpleHttpClient) GetWithHeaders(url string, headers map[string]string) (int, string, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("error creating request: %v", err)
	}

	if headers != nil {
		for name, value := range headers {
			request.Header.Add(name, value)
		}
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return 0, "", fmt.Errorf("error sending http request: %w", err)
	}
	defer response.Body.Close()
	//fmt.Printf("status code: %d\n", response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, "", fmt.Errorf("error reading http response: %w", err)
	}
	// fmt.Printf("Body: %s\n", body)
	return response.StatusCode, string(body), nil
}
