package mylittlesalesman_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/goccy/go-json"
	scraperapi "github.com/zenrows/zenrows-go-sdk/service/api"
)

// PageResponse for MLS listing page
type PageResponse struct {
	Title string   `json:"title" extractor:"h1"`
	URLs  []string `json:"urls" extractor:".content-card-inner > .row > .prhead > h3 > a[href*='/listing/'] @href"`
}

// ListingResponse for MLS detail page
type ListingResponse struct {
	Title         string   `json:"title" extractor:"h1.pb3"`
	Price         string   `json:"price" extractor:"h2#ctl00_ctl00_mc_mc_retailHeader > span > .text-darkred"`
	Description   string   `json:"description" extractor:"#prddesc > div"`
	DetailsLabels []string `json:"details_labels" extractor:"#prddtl > table > tbody > tr > th"`
	DetailsValues []string `json:"details_values" extractor:"#prddtl > table > tbody > tr > td"`
}

// TestZenRowsMLSListingPage tests ZenRows with MLS listing page.
// Run with: ZENROWS_API_KEY=your_key go test -v -run TestZenRowsMLSListingPage
func TestZenRowsMLSListingPage(t *testing.T) {
	apiKey := os.Getenv("ZENROWS_API_KEY")
	if apiKey == "" {
		t.Skip("ZENROWS_API_KEY not set, skipping integration test")
	}

	client := scraperapi.NewClient(scraperapi.WithAPIKey(apiKey))
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	targetURL := "https://www.mylittlesalesman.com/trucks-for-sale-i2c0f0m0?ptid=1&s=11"

	pageResp := PageResponse{}

	t.Logf("Fetching MLS listing page: %s", targetURL)

	resp, err := client.Get(ctx, targetURL, &scraperapi.RequestParameters{
		JSRender:          true,
		UsePremiumProxies: true,
		ProxyCountry:      "US",
		WaitMilliseconds:  30000,
		CSSExtractor:      MarshalExtractor(pageResp),
		CustomParams: map[string]string{
			"antibot": "true", // Try anti-bot bypass
		},
	})
	if err != nil {
		t.Fatalf("Failed to fetch page: %v", err)
	}

	if resp.Error() != nil {
		t.Fatalf("ZenRows returned error: %v", resp.Error())
	}

	t.Logf("Response status: %d", resp.StatusCode())
	t.Logf("Response body length: %d bytes", len(resp.Body()))
	t.Logf("Raw response: %s", string(resp.Body()))

	if err := json.Unmarshal(resp.Body(), &pageResp); err != nil {
		t.Logf("Failed to unmarshal response: %v", err)
		t.Logf("Continuing with empty response...")
	}

	t.Logf("Page title: %s", pageResp.Title)
	t.Logf("Found %d listing URLs", len(pageResp.URLs))

	if len(pageResp.URLs) > 0 {
		t.Logf("First 3 URLs:")
		for i, url := range pageResp.URLs {
			if i >= 3 {
				break
			}
			t.Logf("  - %s", url)
		}
	}

	// Check if blocked
	if pageResp.Title == "" || pageResp.Title == "Checking your browser" {
		t.Logf("WARNING: Possibly blocked")
		t.Logf("Raw response: %s", string(resp.Body()[:min(500, len(resp.Body()))]))
	}

	if pageResp.Title != "Trucks For Sale" {
		t.Logf("WARNING: Unexpected title - might be blocked or different page")
	}
}

// TestZenRowsMLSRawHTML tests ZenRows with MLS to get raw HTML (no CSS extraction).
// Run with: ZENROWS_API_KEY=your_key go test -v -run TestZenRowsMLSRawHTML
func TestZenRowsMLSRawHTML(t *testing.T) {
	apiKey := os.Getenv("ZENROWS_API_KEY")
	if apiKey == "" {
		t.Skip("ZENROWS_API_KEY not set, skipping integration test")
	}

	client := scraperapi.NewClient(scraperapi.WithAPIKey(apiKey))
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	targetURL := "https://www.mylittlesalesman.com/trucks-for-sale-i2c0f0m0?ptid=1&s=11"

	t.Logf("Fetching MLS raw HTML: %s", targetURL)

	// Try with session_id for IP persistence
	resp, err := client.Get(ctx, targetURL, &scraperapi.RequestParameters{
		JSRender:          true,
		UsePremiumProxies: true,
		ProxyCountry:      "US",
		WaitMilliseconds:  30000,
		CustomParams: map[string]string{
			"session_id": "12345", // Maintain same IP across requests
		},
	})
	if err != nil {
		t.Fatalf("Failed to fetch page: %v", err)
	}

	if resp.Error() != nil {
		t.Fatalf("ZenRows returned error: %v", resp.Error())
	}

	t.Logf("Response status: %d", resp.StatusCode())
	t.Logf("Response body length: %d bytes", len(resp.Body()))
	t.Logf("First 1000 bytes: %s", string(resp.Body()[:min(1000, len(resp.Body()))]))
}

// TestZenRowsMLSSupportURL tests the exact URL that ZenRows support claimed worked.
// Run with: ZENROWS_API_KEY=your_key go test -v -run TestZenRowsMLSSupportURL
func TestZenRowsMLSSupportURL(t *testing.T) {
	apiKey := os.Getenv("ZENROWS_API_KEY")
	if apiKey == "" {
		t.Skip("ZENROWS_API_KEY not set, skipping integration test")
	}

	client := scraperapi.NewClient(scraperapi.WithAPIKey(apiKey))
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// This is the exact URL ZenRows support said works
	targetURL := "https://www.mylittlesalesman.com/international-lt625-sleeper-semi-trucks-for-sale-i2c55f339m564445"

	t.Logf("Fetching URL that ZenRows support tested: %s", targetURL)

	resp, err := client.Get(ctx, targetURL, &scraperapi.RequestParameters{
		JSRender:          true,
		UsePremiumProxies: true,
	})
	if err != nil {
		t.Fatalf("Failed to fetch page: %v", err)
	}

	if resp.Error() != nil {
		t.Fatalf("ZenRows returned error: %v", resp.Error())
	}

	body := string(resp.Body())
	t.Logf("Response status: %d", resp.StatusCode())
	t.Logf("Response body length: %d bytes", len(body))
	t.Logf("First 1500 bytes: %s", body[:min(1500, len(body))])

	// Check for blocking indicators
	if contains(body, "Checking your browser") {
		t.Logf("❌ BLOCKED: Got Cloudflare challenge page")
	} else if contains(body, "recaptcha") {
		t.Logf("❌ BLOCKED: Got reCAPTCHA challenge")
	} else if contains(body, "International LT625") || contains(body, "Sleeper") {
		t.Logf("✅ SUCCESS: Got actual content!")
	} else {
		t.Logf("⚠️ UNKNOWN: Check response manually")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestZenRowsMLSDetailPageRaw tests raw HTML for MLS detail page.
// Run with: ZENROWS_API_KEY=your_key go test -v -run TestZenRowsMLSDetailPageRaw
func TestZenRowsMLSDetailPageRaw(t *testing.T) {
	apiKey := os.Getenv("ZENROWS_API_KEY")
	if apiKey == "" {
		t.Skip("ZENROWS_API_KEY not set, skipping integration test")
	}

	client := scraperapi.NewClient(scraperapi.WithAPIKey(apiKey))
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	targetURL := "https://www.mylittlesalesman.com/2022-freightliner-cascadia-126-sleeper-semi-truck-72-condo-sleeper-455hp-12-speed-automatic-14042329"

	t.Logf("Fetching MLS detail page raw HTML: %s", targetURL)

	resp, err := client.Get(ctx, targetURL, &scraperapi.RequestParameters{
		JSRender:          true,
		UsePremiumProxies: true,
		ProxyCountry:      "US",
		WaitMilliseconds:  30000,
	})
	if err != nil {
		t.Fatalf("Failed to fetch page: %v", err)
	}

	body := string(resp.Body())
	t.Logf("Response status: %d", resp.StatusCode())
	t.Logf("Response body length: %d bytes", len(body))
	t.Logf("First 1500 bytes: %s", body[:min(1500, len(body))])

	if contains(body, "Checking your browser") {
		t.Logf("❌ BLOCKED: Got Cloudflare challenge page")
	} else if contains(body, "Freightliner") || contains(body, "Cascadia") {
		t.Logf("✅ SUCCESS: Got actual content!")
	}
}

// TestZenRowsMLSDetailPage tests ZenRows with MLS detail page.
// Run with: ZENROWS_API_KEY=your_key go test -v -run TestZenRowsMLSDetailPage
func TestZenRowsMLSDetailPage(t *testing.T) {
	apiKey := os.Getenv("ZENROWS_API_KEY")
	if apiKey == "" {
		t.Skip("ZENROWS_API_KEY not set, skipping integration test")
	}

	client := scraperapi.NewClient(scraperapi.WithAPIKey(apiKey))
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	targetURL := "https://www.mylittlesalesman.com/2022-freightliner-cascadia-126-sleeper-semi-truck-72-condo-sleeper-455hp-12-speed-automatic-14042329"

	listingResp := ListingResponse{}

	t.Logf("Fetching MLS detail page: %s", targetURL)

	resp, err := client.Get(ctx, targetURL, &scraperapi.RequestParameters{
		JSRender:          true,
		UsePremiumProxies: true,
		ProxyCountry:      "US",
		WaitMilliseconds:  30000,
		CSSExtractor:      MarshalExtractor(listingResp),
		CustomParams: map[string]string{
			"antibot": "true",
		},
	})
	if err != nil {
		t.Fatalf("Failed to fetch page: %v", err)
	}

	if resp.Error() != nil {
		t.Fatalf("ZenRows returned error: %v", resp.Error())
	}

	t.Logf("Response status: %d", resp.StatusCode())
	t.Logf("Response body length: %d bytes", len(resp.Body()))
	t.Logf("Raw response: %s", string(resp.Body()[:min(500, len(resp.Body()))]))

	if err := json.Unmarshal(resp.Body(), &listingResp); err != nil {
		t.Logf("Failed to unmarshal response: %v", err)
	}

	t.Logf("Title: %s", listingResp.Title)
	t.Logf("Price: %s", listingResp.Price)
	t.Logf("Description length: %d chars", len(listingResp.Description))
	t.Logf("Details labels: %v", listingResp.DetailsLabels)
	t.Logf("Details values: %v", listingResp.DetailsValues)

	// Check if blocked
	if listingResp.Title == "" {
		t.Logf("WARNING: No title - possibly blocked or page not found")
		t.Logf("Raw response: %s", string(resp.Body()[:min(500, len(resp.Body()))]))
	}
}
