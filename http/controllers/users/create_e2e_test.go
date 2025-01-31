package users

import (
	"github.com/tebeka/selenium"
	"testing"
	"time"
)

const (
	baseURL  = "http://localhost:8080/static/main.html"
	Nickname = "E2Etest"
	Password = "testpassword123"
)

func TestSendEmailE2E(t *testing.T) {
	caps := selenium.Capabilities{"browserName": "chrome"}
	driver, err := selenium.NewRemote(caps, "http://localhost:4444")
	if err != nil {
		t.Fatalf("Error connecting to WebDriver instance: %v", err)
	}
	defer driver.Quit()

	if err := driver.Get(baseURL); err != nil {
		t.Fatalf("Failed to load page: %v", err)
	}

	time.Sleep(3 * time.Second)

	currentURL, _ := driver.CurrentURL()
	t.Logf("Current URL: %s", currentURL)

	NicknameElement, err := driver.FindElement(selenium.ByID, "nickname")
	if err != nil {
		t.Fatalf("Failed to locate nickname recipient field: %v", err)
	}
	err = NicknameElement.SendKeys(Nickname)
	if err != nil {
		t.Fatalf("Failed to enter nickname recipient: %v", err)
	}

	PasswordElement, err := driver.FindElement(selenium.ByID, "password")
	if err != nil {
		t.Fatalf("Failed to locate password subject field: %v", err)
	}
	err = PasswordElement.SendKeys(Password)
	if err != nil {
		t.Fatalf("Failed to enter password subject: %v", err)
	}

	sendButton, err := driver.FindElement(selenium.ByCSSSelector, "button[onclick='addUser()']")
	if err != nil {
		t.Fatalf("Failed to locate add button: %v", err)
	}
	err = sendButton.Click()
	if err != nil {
		t.Fatalf("Failed to click add button: %v", err)
	}

	time.Sleep(10 * time.Second)
	alertText, err := driver.AlertText()
	if err != nil {
		t.Fatalf("No alert found: %v", err)
	}

	expectedMessage := "User added successfully!"
	if alertText != expectedMessage {
		t.Fatalf("Unexpected alert message: got %s, want %s", alertText, expectedMessage)
	}

	err = driver.AcceptAlert()
	if err != nil {
		t.Fatalf("Failed to accept alert: %v", err)
	}

	t.Logf("E2E Test Passed: Email sent successfully")
}
