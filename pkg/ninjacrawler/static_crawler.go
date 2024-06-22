package ninjacrawler

import (
	"crypto/tls"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/proxy"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (app *Crawler) getHttpClient() *http.Client {
	client := &http.Client{
		//Timeout: 30 * time.Second,
		Timeout: (app.engine.Timeout / 1000) * time.Second,
	}
	return client
}

func (app *Crawler) NavigateToStaticURL(client *http.Client, urlString string, proxyServer Proxy) (*goquery.Document, error) {
	if len(app.engine.ProxyServers) > 0 {
		// Create the proxy URL
		proxyURL, err := url.Parse(proxyServer.Server)
		if err != nil {
			log.Fatalf("Failed to parse proxy URL: %v", err)
		}

		// Set the username and password in the proxy URL
		if proxyServer.Username != "" && proxyServer.Password != "" {
			proxyURL.User = url.UserPassword(proxyServer.Username, proxyServer.Password)
		}

		// Create a proxy dialer
		dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
		if err != nil {
			log.Fatalf("Failed to obtain proxy dialer: %v", err)
		}

		// Create an HTTP client and set the transport to use the proxy dialer
		httpTransport := &http.Transport{
			Dial: dialer.Dial,
		}

		// Add TLS configuration for HTTPS proxy if needed
		if proxyURL.Scheme == "https" {
			httpTransport.TLSClientConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}

		client.Transport = httpTransport
	}
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %v", err)
	}
	// Set headers (optional but recommended)
	req.Header.Set("User-Agent", app.Config.GetString("USER_AGENT"))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("failed to fetch page: %v", resp.Status)
		app.Logger.Html(string(body), urlString, msg)
		return nil, fmt.Errorf(msg)
	}

	document, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	return document, nil
}
