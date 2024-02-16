package centrifugo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	baseURL string
	apiKey  string
	client  http.Client
}

type transport struct {
	apiKey           string
	defaultTransport http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("x-api-key", t.apiKey)
	return t.defaultTransport.RoundTrip(req)
}

func New(baseURL, apiKey string) Client {
	return Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		client: http.Client{Transport: &transport{
			apiKey:           apiKey,
			defaultTransport: http.DefaultTransport,
		}},
	}
}

func (c *Client) Broadcast(channels []string, data map[string]interface{}) error {
	url := fmt.Sprintf("%s/api/broadcast", c.baseURL)
	rawPayload := map[string]interface{}{
		"channels": channels,
		"data":     data,
	}
	payload, err := json.Marshal(rawPayload)
	if err != nil {
		return fmt.Errorf("unable to marshal json in centrifugo broadcast request: %w", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("unable to create new broadcast request to centrifugo: %w", err)
	}
	_, err = c.client.Do(req)
	if err != nil {
		return fmt.Errorf("unable to get response from centrifugo: %w", err)
	}
	return nil
}

func (c *Client) BroadcastAllExcept(except string, data map[string]interface{}) error {
	channels, err := c.Channels("*")
	if err != nil {
		return fmt.Errorf("unable to get channels in broadcast except: %w", err)
	}
	filteredChannels := make([]string, 0)
	for _, v := range channels {
		if v == except {
			continue
		}
		filteredChannels = append(filteredChannels, v)
	}
	if len(filteredChannels) == 0 {
		return nil
	}
	return c.Broadcast(filteredChannels, data)
}

type ChannelInfo struct {
	NumClients int64 `json:"num_clients"`
}

type ChannelsOutResult struct {
	Channels map[string]ChannelInfo `json:"channels"`
}

type ChannelsOut struct {
	Result ChannelsOutResult `json:"result"`
}

func (c *Client) Channels(pattern string) ([]string, error) {
	url := fmt.Sprintf("%s/api/channels", c.baseURL)
	rawPayload := map[string]string{"patter": pattern}
	payload, err := json.Marshal(rawPayload)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal json in centrifugo channels request: %w", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("unable to create new channels request to centrifugo: %w", err)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to get response from centrifugo: %w", err)
	}
	defer resp.Body.Close()
	var data ChannelsOut
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("unable to unmarshal channels response: %w", err)
	}
	channels := make([]string, 0, len(data.Result.Channels))
	for k := range data.Result.Channels {
		channels = append(channels, k)
	}
	return channels, nil
}
