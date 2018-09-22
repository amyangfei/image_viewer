package viewer

import (
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"log"
)

func StartChrome(port int) (*selenium.Service, selenium.WebDriver, error) {
	var err error
	opts := []selenium.ServiceOption{}
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}

	// Disable image rendering
	imagCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}

	chromeCaps := chrome.Capabilities{
		Prefs: imagCaps,
		Path:  "",
		Args: []string{
			"--headless",
			"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36",
		},
	}
	caps.AddChrome(chromeCaps)

	// Start chrome driver
	service, err := selenium.NewChromeDriverService("/usr/local/bin/chromedriver", port, opts...)
	if err != nil {
		log.Printf("Error starting the ChromeDriver server: %v", err)
		return nil, nil, err
	}

	// Start chrome
	webDriver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		return service, nil, err
	}

	// Don't forget to stop chromedriver and shutdown webDriver
	// service.Stop()   // stop chromedriver
	// webDriver.Quit() // shutdown chrome client

	return service, webDriver, nil
}
