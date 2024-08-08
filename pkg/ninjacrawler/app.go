package ninjacrawler

import (
	"fmt"
	"github.com/playwright-community/playwright-go"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var startTime time.Time

const (
	baseCollection = "sites"
)

type Crawler struct {
	*mongo.Client
	Config                *configService
	Name                  string
	Url                   string
	BaseUrl               string
	pw                    *playwright.Playwright
	UrlSelectors          []UrlSelector
	ProductDetailSelector ProductDetailSelector
	engine                *Engine
	Logger                *defaultLogger
	httpClient            *http.Client
	isLocalEnv            bool
	preference            *AppPreference
	userAgent             string
	CurrentProxy          Proxy
}

func NewCrawler(name, url string, engines ...Engine) *Crawler {

	defaultPreference := getDefaultPreference()
	defaultEngine := getDefaultEngine()
	if len(engines) > 0 {
		eng := engines[0]
		overrideEngineDefaults(&defaultEngine, &eng)
	}
	// Handle other engine overrides as needed
	config := newConfig()

	crawler := &Crawler{
		Name:         name,
		Url:          url,
		engine:       &defaultEngine,
		Config:       config,
		CurrentProxy: Proxy{},
	}

	logger := newDefaultLogger(crawler, name)
	crawler.Logger = logger
	crawler.Client = crawler.mustGetClient()
	crawler.BaseUrl = crawler.getBaseUrl(url)
	crawler.isLocalEnv = config.GetString("APP_ENV") == "local"
	crawler.userAgent = config.GetString("USER_AGENT")
	crawler.preference = &defaultPreference
	return crawler
}

func (app *Crawler) Start() {
	defer func() {
		if r := recover(); r != nil {
			app.Logger.Error("Recovered in Start: %v", r)
		}
	}()
	startTime = time.Now()
	app.Logger.Info("Crawler Started! 🚀")

	deleteDB := app.Config.GetBool("DELETE_DB")
	if deleteDB {
		err := app.dropDatabase()
		if err != nil {
			return
		}
	}
	app.newSite()
	app.toggleClient()
}

func (app *Crawler) toggleClient() {
	if app.engine.IsDynamic {
		pw, err := app.GetPlaywright()
		if err != nil {
			app.Logger.Fatal("failed to initialize playwright: %v\n", err)
			return // exit if playwright initialization fails
		}
		app.pw = pw
	} else {
		app.httpClient = app.GetHttpClient()
	}
}

func (app *Crawler) Stop() {
	defer func() {
		if r := recover(); r != nil {
			app.Logger.Error("Recovered in Stop: %v", r)
		}
	}()
	if app.pw != nil {
		app.pw.Stop()
	}
	if app.Client != nil {
		app.closeClient()
	}
	// upload logs
	app.UploadLogs()
	duration := time.Since(startTime)
	app.Logger.Info("Crawler stopped in ⚡ %v", duration)
}

func (app *Crawler) UploadLogs() {
	storagePath := fmt.Sprintf("storage/logs/%s", app.Name)
	err := filepath.Walk(storagePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			app.Logger.Error("Error accessing path %s: %v", path, err)
			return err
		}

		if !info.IsDir() {
			relativePath := strings.TrimPrefix(path, storagePath+"/")
			uploadToBucket(app, path, fmt.Sprintf("logs/%s", relativePath))
		}

		return nil
	})

	if err != nil {
		app.Logger.Error("Error walking through storage directory: %v", err)
	}
}

func (app *Crawler) GetBaseCollection() string {
	return baseCollection
}

func (app *Crawler) SetPreference(preference AppPreference) *Crawler {

	defaultPreference := getDefaultPreference()

	overridePreferenceDefaults(&defaultPreference, &preference)
	app.preference = &defaultPreference
	return app
}

func (app *Crawler) Handle(handler Handler) {
	defer app.Stop() // Ensure Stop is called after handlers
	app.Start()

	if handler.UrlHandler != nil {
		handler.UrlHandler(app)
	}
	if handler.ProductHandler != nil {
		handler.ProductHandler(app)
	}
}
func (app *Crawler) AutoHandle(configs []ProcessorConfig) {
	defer app.Stop() // Ensure Stop is called after handlers
	app.Start()

	app.CrawlUrls(configs)
}
func getDefaultPreference() AppPreference {
	return AppPreference{
		ExcludeUniqueUrlEntities: []string{},
	}
}
func overridePreferenceDefaults(defaultPreference *AppPreference, preference *AppPreference) {
	if len(preference.ExcludeUniqueUrlEntities) > 0 {
		defaultPreference.ExcludeUniqueUrlEntities = preference.ExcludeUniqueUrlEntities
	}
}

func getDefaultEngine() Engine {
	return Engine{
		BrowserType:             "chromium",
		Provider:                "http",
		ConcurrentLimit:         1,
		IsDynamic:               false,
		WaitForDynamicRendering: false,
		DevCrawlLimit:           100,
		BlockResources:          false,
		JavaScriptEnabled:       true,
		BlockedURLs: []string{
			"www.googletagmanager.com",
			"google.com",
			"googleapis.com",
			"gstatic.com",
		},
		BoostCrawling:          false,
		ProxyServers:           []Proxy{},
		CookieConsent:          nil,
		Timeout:                30 * 1000, // 30 sec
		SleepAfter:             1000,
		MaxRetryAttempts:       3,
		ForceInstallPlaywright: false,
		Args:                   []string{},
		ProviderOption: ProviderQueryOption{
			JsRender:             false,
			UsePremiumProxyRetry: false,
		},
		SleepDuration: 10,
	}
}

func overrideEngineDefaults(defaultEngine *Engine, eng *Engine) {
	if eng.BrowserType != "" {
		defaultEngine.BrowserType = eng.BrowserType
	}
	if eng.Provider != "" {
		defaultEngine.Provider = eng.Provider
	}
	if eng.ConcurrentLimit > 0 {
		defaultEngine.ConcurrentLimit = eng.ConcurrentLimit
	}
	if eng.IsDynamic {
		defaultEngine.IsDynamic = eng.IsDynamic
	}
	if eng.WaitForDynamicRendering {
		defaultEngine.WaitForDynamicRendering = eng.WaitForDynamicRendering
	}
	if eng.DevCrawlLimit > 0 {
		defaultEngine.DevCrawlLimit = eng.DevCrawlLimit
	}
	if eng.BlockResources {
		defaultEngine.BlockResources = eng.BlockResources
	}
	if eng.JavaScriptEnabled {
		defaultEngine.JavaScriptEnabled = eng.JavaScriptEnabled
	}
	if eng.BoostCrawling {
		defaultEngine.BoostCrawling = eng.BoostCrawling
		defaultEngine.ProxyServers = eng.getProxyList()
	}
	if len(eng.ProxyServers) > 0 {
		config := newConfig()
		zenrowsApiKey := config.EnvString("ZENROWS_API_KEY")
		for _, proxy := range eng.ProxyServers {
			if proxy.Server == ZENROWS {
				proxy.Server = fmt.Sprintf("http://%s:@proxy.zenrows.com:8001", zenrowsApiKey)
			}
			defaultEngine.ProxyServers = append(defaultEngine.ProxyServers, proxy)
		}
	}
	if eng.CookieConsent != nil {
		defaultEngine.CookieConsent = eng.CookieConsent
	}
	if eng.Timeout > 0 {
		defaultEngine.Timeout = eng.Timeout * 1000
	}
	if eng.SleepAfter > 0 {
		defaultEngine.SleepAfter = eng.SleepAfter
	}
	if eng.MaxRetryAttempts > 0 {
		defaultEngine.MaxRetryAttempts = eng.MaxRetryAttempts
	}
	if eng.ForceInstallPlaywright {
		defaultEngine.ForceInstallPlaywright = eng.ForceInstallPlaywright
	}
	if len(eng.Args) > 0 {
		defaultEngine.Args = eng.Args
	}

	if eng.ProviderOption.JsRender {
		defaultEngine.ProviderOption.JsRender = eng.ProviderOption.JsRender
	}

	if eng.ProviderOption.UsePremiumProxyRetry {
		defaultEngine.ProviderOption.UsePremiumProxyRetry = eng.ProviderOption.UsePremiumProxyRetry
	}
	defaultEngine.BlockedURLs = append(defaultEngine.BlockedURLs, eng.BlockedURLs...)
	if eng.SleepDuration > 0 {
		defaultEngine.SleepDuration = eng.SleepDuration
	}
}
