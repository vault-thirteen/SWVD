package pw

import (
	"bytes"
	"encoding/base64"
	"os"
	"time"

	"github.com/mxschmitt/playwright-go"
	"github.com/vault-thirteen/auxie/errors"
)

//go:generate go install github.com/mxschmitt/playwright-go/cmd/playwright@v0.6401.0
//go:generate playwright install --with-deps

const (
	TimeoutMs            = 5 * 60 * 1000
	DoH_Server_Url       = "https://dns.google/dns-query"
	Dns_BootstrapAddress = "8.8.8.8"
)

func DownloadFile(url string, dstFilePath string) (err error) {
	err = playwright.Install()
	if err != nil {
		return err
	}

	var pw *playwright.Playwright
	pw, err = playwright.Run()
	if err != nil {
		return err
	}
	defer func() {
		derr := pw.Stop()
		err = errors.Combine(err, derr)
	}()

	userPrefs := map[string]any{
		// 2 = DoH first with fallback to native DNS
		// 3 = DoH only (strict mode, no fallback)
		"network.trr.mode": 3,
		// Custom DoH provider endpoint
		"network.trr.uri": DoH_Server_Url,

		// Optional: Bootstrap IP so Firefox doesn't need plain DNS to resolve the DoH host
		"network.trr.bootstrapAddress": Dns_BootstrapAddress,

		// Disable tracking protection.
		"browser.contentblocking.category":                  "custom",
		"network.cookie.cookieBehavior":                     0, // 0 = Accept all cookies
		"privacy.trackingprotection.enabled":                false,
		"privacy.trackingprotection.pbmode.enabled":         false,
		"privacy.trackingprotection.socialtracking.enabled": false,
		"privacy.trackingprotection.cryptomining.enabled":   false,
		"privacy.trackingprotection.fingerprinting.enabled": false,
		"privacy.annotate_channels.strict_list.enabled":     false,

		// ?
		"network.proxy.allow_hijacking_localhost": true,
	}

	var browser playwright.Browser
	browser, err = pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{
		Headless:         new(false),
		FirefoxUserPrefs: userPrefs,
	})
	if err != nil {
		return err
	}
	defer func() {
		derr := browser.Close()
		err = errors.Combine(err, derr)
	}()

	var context playwright.BrowserContext
	context, err = browser.NewContext()
	if err != nil {
		return err
	}
	defer func() {
		derr := context.Close()
		err = errors.Combine(err, derr)
	}()

	context.SetDefaultTimeout(TimeoutMs)
	context.SetDefaultNavigationTimeout(TimeoutMs)

	var page playwright.Page
	page, err = context.NewPage()
	if err != nil {
		return err
	}

	page.SetDefaultTimeout(TimeoutMs)
	page.SetDefaultNavigationTimeout(TimeoutMs)

	var videoBuffer bytes.Buffer
	err = page.ExposeFunction("pushChunk", func(args ...any) any {
		b64Str, ok := args[0].(string)
		if !ok {
			return nil
		}

		chunk, _ := base64.StdEncoding.DecodeString(b64Str)
		videoBuffer.Write(chunk)
		return nil
	})
	if err != nil {
		return err
	}

	_, err = page.Goto(url)
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 1)

	_, err = page.Evaluate(`async () => {
	const res = await fetch(window.location.href);
	const reader = res.body.getReader();

	while (true) {
		const { done, value } = await reader.read();
		if (done) break;

		// Convert Uint8Array to base64
		let binary = '';
		for (let i = 0; i < value.byteLength; i++) {
			binary += String.fromCharCode(value[i]);
		}
		await window.pushChunk(btoa(binary));
	}
}`)
	if err != nil {
		return err
	}

	videoBytes := videoBuffer.Bytes()

	err = os.WriteFile(dstFilePath, videoBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
