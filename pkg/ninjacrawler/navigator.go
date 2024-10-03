package ninjacrawler

import (
	"strings"
	"sync/atomic"
	"time"
)

func (app *Crawler) Navigate(url string) (*NavigationContext, error) {
	app.openBrowsers(app.CurrentProxy)
	defer app.closeBrowsers()
	app.openPages()
	defer app.closePages()

	atomic.AddInt32(&app.ReqCount, 1)

	// Freeze all goroutines after n requests
	if app.ReqCount > 0 && atomic.LoadInt32(&app.ReqCount)%int32(app.engine.SleepAfter) == 0 {
		app.Logger.Info("Sleeping %d seconds after %d operations", app.engine.SleepDuration, app.engine.SleepAfter)
		time.Sleep(time.Duration(app.engine.SleepDuration) * time.Second)
	}

	navigationContext, err := app.navigateTo(url, "DeepLink", false)

	if err != nil {
		if strings.Contains(err.Error(), "StatusCode:404") {
			return nil, err
		} else if strings.Contains(err.Error(), "isRetryable") {
			// Rotate proxy if it's a retryable error
			if len(app.engine.ProxyServers) > 0 && app.engine.ProxyStrategy == ProxyStrategyRotation {
				nextProxyIndex := (app.CurrentProxyIndex + 1) % len(app.engine.ProxyServers)
				app.Logger.Summary("Error with proxy %s: %v. Retrying with a different proxy: %s", app.CurrentProxy.Server, err.Error(), app.engine.ProxyServers[nextProxyIndex].Server)

				app.CurrentProxyIndex = nextProxyIndex
				app.CurrentProxy = app.engine.ProxyServers[nextProxyIndex]

				app.Logger.Info("Sleeping %d seconds before retrying", app.engine.SleepDuration)
				time.Sleep(time.Duration(app.engine.SleepDuration) * time.Second)

				// Retry with the next proxy and return the result
				return app.Navigate(url)
			}
			if app.engine.RetrySleepDuration > 0 {
				app.HandleThrottling(1, 0)
			}

			return nil, err
		} else {
			return nil, err
		}
	}
	return navigationContext, nil

}
