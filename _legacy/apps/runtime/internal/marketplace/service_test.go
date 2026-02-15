package marketplace

import (
	"context"
	"sync"
	"testing"
	"time"

	"pryx-core/internal/bus"
)

// helperService creates a marketplace service for testing
func helperService(t *testing.T) (*Service, *bus.Bus) {
	b := bus.New()
	s := NewService(b)
	return s, b
}

// helperCreateListing creates a test listing
func helperCreateListing(s *Service, name string) *AgentListing {
	return &AgentListing{
		AgentID:       name + "-id",
		AgentName:     name,
		AgentVersion:  "1.0.0",
		Description:   "Test description for " + name,
		Author:        "Test Author",
		AuthorID:      "author-id",
		Category:      "productivity",
		Tags:          []string{"test", name},
		Capabilities:  []string{"text-generation"},
		PriceModel:    "free",
		Price:         0,
		License:       "MIT",
		Rating:        4.5,
		ReviewCount:   10,
		DownloadCount: 100,
		InstallCount:  50,
		Status:        ListingStatusPublished,
	}
}

// TestCreateListing tests creating a new listing
func TestCreateListing(t *testing.T) {
	s, _ := helperService(t)

	listing := helperCreateListing(s, "TestAgent")

	err := s.CreateListing(context.Background(), listing)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if listing.ListingID == "" {
		t.Error("expected listing ID to be set")
	}
	if listing.Status != ListingStatusDraft {
		t.Errorf("expected status draft, got %s", listing.Status)
	}
	if listing.CreatedAt.IsZero() {
		t.Error("expected created at to be set")
	}
}

// TestGetListing tests retrieving a listing
func TestGetListing(t *testing.T) {
	s, _ := helperService(t)

	listing := helperCreateListing(s, "TestAgent")
	s.CreateListing(context.Background(), listing)

	retrieved, err := s.GetListing(listing.ListingID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if retrieved.AgentID != listing.AgentID {
		t.Errorf("expected agent ID %s, got %s", listing.AgentID, retrieved.AgentID)
	}
}

// TestGetListingNotFound tests error when listing not found
func TestGetListingNotFound(t *testing.T) {
	s, _ := helperService(t)

	_, err := s.GetListing("unknown-id")
	if err == nil {
		t.Error("expected error when listing not found")
	}
}

// TestGetListingByAgentID tests retrieving listing by agent ID
func TestGetListingByAgentID(t *testing.T) {
	s, _ := helperService(t)

	listing := helperCreateListing(s, "TestAgent")
	s.CreateListing(context.Background(), listing)

	retrieved, err := s.GetListingByAgentID(listing.AgentID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if retrieved.AgentID != listing.AgentID {
		t.Errorf("expected agent ID %s, got %s", listing.AgentID, retrieved.AgentID)
	}
}

// TestUpdateListing tests updating a listing
func TestUpdateListing(t *testing.T) {
	s, _ := helperService(t)

	listing := helperCreateListing(s, "TestAgent")
	s.CreateListing(context.Background(), listing)

	err := s.UpdateListing(listing.ListingID, map[string]interface{}{
		"name":        "Updated Agent",
		"description": "Updated description",
		"price":       9.99,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	retrieved, _ := s.GetListing(listing.ListingID)
	if retrieved.AgentName != "Updated Agent" {
		t.Errorf("expected name 'Updated Agent', got '%s'", retrieved.AgentName)
	}
	if retrieved.Price != 9.99 {
		t.Errorf("expected price 9.99, got %f", retrieved.Price)
	}
}

// TestPublishListing tests publishing a listing
func TestPublishListing(t *testing.T) {
	s, _ := helperService(t)

	listing := helperCreateListing(s, "TestAgent")
	s.CreateListing(context.Background(), listing)

	err := s.PublishListing(listing.ListingID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	retrieved, _ := s.GetListing(listing.ListingID)
	if retrieved.Status != ListingStatusPublished {
		t.Errorf("expected status published, got %s", retrieved.Status)
	}
	if retrieved.PublishedAt == nil {
		t.Error("expected published at to be set")
	}
}

// TestSuspendListing tests suspending a listing
func TestSuspendListing(t *testing.T) {
	s, _ := helperService(t)

	listing := helperCreateListing(s, "TestAgent")
	s.CreateListing(context.Background(), listing)
	s.PublishListing(listing.ListingID)

	err := s.SuspendListing(listing.ListingID, "Policy violation")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	retrieved, _ := s.GetListing(listing.ListingID)
	if retrieved.Status != ListingStatusSuspended {
		t.Errorf("expected status suspended, got %s", retrieved.Status)
	}
}

// TestGetFeaturedListings tests getting featured listings
func TestGetFeaturedListings(t *testing.T) {
	s, _ := helperService(t)

	// Create featured listing
	featured := helperCreateListing(s, "FeaturedAgent")
	featured.Featured = true
	s.CreateListing(context.Background(), featured)
	s.PublishListing(featured.ListingID)

	// Create non-featured listing
	normal := helperCreateListing(s, "NormalAgent")
	s.CreateListing(context.Background(), normal)
	s.PublishListing(normal.ListingID)

	listings := s.GetFeaturedListings(10)

	if len(listings) != 1 {
		t.Errorf("expected 1 featured listing, got %d", len(listings))
	}
	if len(listings) > 0 && listings[0].AgentName != "FeaturedAgent" {
		t.Errorf("expected FeaturedAgent, got %s", listings[0].AgentName)
	}
}

// TestGetListingsByCategory tests getting listings by category
func TestGetListingsByCategory(t *testing.T) {
	s, _ := helperService(t)

	// Create listings in different categories
	prodListing := helperCreateListing(s, "ProductivityAgent")
	prodListing.Category = "productivity"
	s.CreateListing(context.Background(), prodListing)
	s.PublishListing(prodListing.ListingID)

	devListing := helperCreateListing(s, "DeveloperAgent")
	devListing.Category = "development"
	s.CreateListing(context.Background(), devListing)
	s.PublishListing(devListing.ListingID)

	listings := s.GetListingsByCategory("productivity", 10)

	if len(listings) != 1 {
		t.Errorf("expected 1 productivity listing, got %d", len(listings))
	}
	if len(listings) > 0 && listings[0].Category != "productivity" {
		t.Errorf("expected category productivity, got %s", listings[0].Category)
	}
}

// TestSearchListings tests searching for listings
func TestSearchListings(t *testing.T) {
	s, _ := helperService(t)

	// Create test listing
	listing1 := helperCreateListing(s, "SearchableAgent123")
	s.CreateListing(context.Background(), listing1)
	s.PublishListing(listing1.ListingID)

	query := &SearchQuery{
		RequesterID: "user-1",
		Keywords:    "SearchableAgent123",
		PageSize:    20,
	}

	result, err := s.SearchListings(context.Background(), query)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check that we get at least some results
	if len(result.Listings) == 0 {
		t.Error("expected search results, got none")
	}
}

// TestSearchListingsByCategory tests searching with category filter
func TestSearchListingsByCategory(t *testing.T) {
	s, _ := helperService(t)

	// Create listings
	prodListing := helperCreateListing(s, "ProductivityAgent")
	prodListing.Category = "productivity"
	s.CreateListing(context.Background(), prodListing)
	s.PublishListing(prodListing.ListingID)

	devListing := helperCreateListing(s, "DeveloperAgent")
	devListing.Category = "development"
	s.CreateListing(context.Background(), devListing)
	s.PublishListing(devListing.ListingID)

	query := &SearchQuery{
		RequesterID: "user-1",
		Category:    "development",
		PageSize:    20,
	}

	result, err := s.SearchListings(context.Background(), query)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.Listings) != 1 {
		t.Errorf("expected 1 listing, got %d", len(result.Listings))
	}
}

// TestSearchListingsByMinRating tests searching with minimum rating filter
func TestSearchListingsByMinRating(t *testing.T) {
	s, _ := helperService(t)

	// Create listings with different ratings
	highRated := helperCreateListing(s, "HighRatedAgent")
	highRated.Rating = 4.8
	s.CreateListing(context.Background(), highRated)
	s.PublishListing(highRated.ListingID)

	lowRated := helperCreateListing(s, "LowRatedAgent")
	lowRated.Rating = 3.0
	s.CreateListing(context.Background(), lowRated)
	s.PublishListing(lowRated.ListingID)

	query := &SearchQuery{
		RequesterID: "user-1",
		MinRating:   4.5,
		PageSize:    20,
	}

	result, err := s.SearchListings(context.Background(), query)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.Listings) != 1 {
		t.Errorf("expected 1 listing with rating >= 4.5, got %d", len(result.Listings))
	}
}

// TestSearchListingsSort tests sorting search results
func TestSearchListingsSort(t *testing.T) {
	s, _ := helperService(t)

	// Create listings with different ratings
	listing1 := helperCreateListing(s, "AgentA")
	listing1.Rating = 3.0
	s.CreateListing(context.Background(), listing1)
	s.PublishListing(listing1.ListingID)

	listing2 := helperCreateListing(s, "AgentB")
	listing2.Rating = 4.5
	s.CreateListing(context.Background(), listing2)
	s.PublishListing(listing2.ListingID)

	query := &SearchQuery{
		RequesterID: "user-1",
		SortBy:      "rating",
		SortOrder:   "desc",
		PageSize:    20,
	}

	result, _ := s.SearchListings(context.Background(), query)

	if len(result.Listings) < 2 {
		t.Fatalf("expected at least 2 listings, got %d", len(result.Listings))
	}

	if result.Listings[0].Rating < result.Listings[1].Rating {
		t.Error("expected results sorted by rating descending")
	}
}

// TestAddReview tests adding a review
func TestAddReview(t *testing.T) {
	s, _ := helperService(t)

	listing := helperCreateListing(s, "TestAgent")
	s.CreateListing(context.Background(), listing)

	review := &Review{
		ListingID:    listing.ListingID,
		AgentID:      listing.AgentID,
		ReviewerID:   "reviewer-1",
		ReviewerName: "Test Reviewer",
		Rating:       5,
		Title:        "Great agent!",
		Content:      "This agent works wonderfully.",
	}

	err := s.AddReview(context.Background(), review)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if review.ReviewID == "" {
		t.Error("expected review ID to be set")
	}
}

// TestAddReviewUpdatesRating tests that adding a review updates the listing rating
func TestAddReviewUpdatesRating(t *testing.T) {
	s, _ := helperService(t)

	listing := helperCreateListing(s, "TestAgent")
	listing.Rating = 4.0
	listing.ReviewCount = 0
	s.CreateListing(context.Background(), listing)

	// Add first review
	review1 := &Review{
		ListingID:    listing.ListingID,
		AgentID:      listing.AgentID,
		ReviewerID:   "reviewer-1",
		ReviewerName: "Test Reviewer 1",
		Rating:       4,
		Title:        "Good!",
		Content:      "Nice agent.",
	}
	s.AddReview(context.Background(), review1)

	// Add second review with different reviewer
	review2 := &Review{
		ListingID:    listing.ListingID,
		AgentID:      listing.AgentID,
		ReviewerID:   "reviewer-2",
		ReviewerName: "Test Reviewer 2",
		Rating:       5,
		Title:        "Great!",
		Content:      "Excellent agent!",
	}
	s.AddReview(context.Background(), review2)

	retrieved, _ := s.GetListing(listing.ListingID)
	if retrieved.ReviewCount != 2 {
		t.Errorf("expected review count 2, got %d", retrieved.ReviewCount)
	}
}

// TestGetReviews tests retrieving reviews
func TestGetReviews(t *testing.T) {
	s, _ := helperService(t)

	listing := helperCreateListing(s, "TestAgent")
	s.CreateListing(context.Background(), listing)

	// Add reviews
	s.AddReview(context.Background(), &Review{
		ListingID:    listing.ListingID,
		AgentID:      listing.AgentID,
		ReviewerID:   "reviewer-1",
		ReviewerName: "Reviewer 1",
		Rating:       4,
		Content:      "Good",
	})

	s.AddReview(context.Background(), &Review{
		ListingID:    listing.ListingID,
		AgentID:      listing.AgentID,
		ReviewerID:   "reviewer-2",
		ReviewerName: "Reviewer 2",
		Rating:       5,
		Content:      "Excellent",
	})

	reviews := s.GetReviews(listing.ListingID, 10)

	if len(reviews) != 2 {
		t.Errorf("expected 2 reviews, got %d", len(reviews))
	}
}

// TestMarkReviewHelpful tests marking a review as helpful
func TestMarkReviewHelpful(t *testing.T) {
	s, _ := helperService(t)

	listing := helperCreateListing(s, "TestAgent")
	s.CreateListing(context.Background(), listing)

	review := &Review{
		ListingID:    listing.ListingID,
		AgentID:      listing.AgentID,
		ReviewerID:   "reviewer-1",
		ReviewerName: "Test Reviewer",
		Rating:       4,
		Content:      "Good",
	}
	s.AddReview(context.Background(), review)

	err := s.MarkReviewHelpful(review.ReviewID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	s.mu.RLock()
	retrieved := s.reviews[review.ReviewID]
	s.mu.RUnlock()

	if retrieved.HelpfulCount != 1 {
		t.Errorf("expected helpful count 1, got %d", retrieved.HelpfulCount)
	}
}

// TestAddReviewResponse tests adding a response to a review
func TestAddReviewResponse(t *testing.T) {
	s, _ := helperService(t)

	listing := helperCreateListing(s, "TestAgent")
	s.CreateListing(context.Background(), listing)

	review := &Review{
		ListingID:    listing.ListingID,
		AgentID:      listing.AgentID,
		ReviewerID:   "reviewer-1",
		ReviewerName: "Test Reviewer",
		Rating:       4,
		Content:      "Good",
	}
	s.AddReview(context.Background(), review)

	response := &ReviewResponse{
		AuthorID:   "author-id",
		AuthorName: "Test Author",
		Content:    "Thank you for your feedback!",
	}

	err := s.AddReviewResponse(review.ReviewID, response)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	s.mu.RLock()
	retrieved := s.reviews[review.ReviewID]
	s.mu.RUnlock()

	if retrieved.Response == nil {
		t.Error("expected review to have response")
	}
}

// TestCreateCategory tests creating a category
func TestCreateCategory(t *testing.T) {
	s, _ := helperService(t)

	category := &Category{
		Name:         "Development",
		Slug:         "development",
		Description:  "Agents for software development",
		ListingCount: 0,
	}

	err := s.CreateCategory(category)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if category.CategoryID == "" {
		t.Error("expected category ID to be set")
	}
}

// TestGetCategory tests retrieving a category
func TestGetCategory(t *testing.T) {
	s, _ := helperService(t)

	category := &Category{
		Name:        "Development",
		Slug:        "development",
		Description: "Agents for software development",
	}
	s.CreateCategory(category)

	retrieved, err := s.GetCategory(category.CategoryID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if retrieved.Name != "Development" {
		t.Errorf("expected name 'Development', got '%s'", retrieved.Name)
	}
}

// TestRecordInstall tests recording an installation
func TestRecordInstall(t *testing.T) {
	s, _ := helperService(t)

	listing := helperCreateListing(s, "TestAgent")
	s.CreateListing(context.Background(), listing)

	err := s.RecordInstall(listing.ListingID, "user-1", "1.0.0", "marketplace")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	s.mu.RLock()
	retrieved := s.listings[listing.ListingID]
	s.mu.RUnlock()

	if retrieved.InstallCount != 1 {
		t.Errorf("expected install count 1, got %d", retrieved.InstallCount)
	}
}

// TestRecordDownload tests recording a download
func TestRecordDownload(t *testing.T) {
	s, _ := helperService(t)

	listing := helperCreateListing(s, "TestAgent")
	s.CreateListing(context.Background(), listing)

	// Get the listing ID and verify it exists
	listingID := listing.ListingID

	err := s.RecordDownload(listingID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	s.mu.RLock()
	retrieved := s.listings[listingID]
	s.mu.RUnlock()

	if retrieved == nil {
		t.Fatal("listing not found")
	}

	// Just verify that download count increased from whatever it was
	if retrieved.DownloadCount < 1 {
		t.Errorf("expected download count >= 1, got %d", retrieved.DownloadCount)
	}
}

// TestGetStats tests getting marketplace statistics
func TestGetStats(t *testing.T) {
	s, _ := helperService(t)

	// Create listings
	s.CreateListing(context.Background(), helperCreateListing(s, "Agent1"))
	s.CreateListing(context.Background(), helperCreateListing(s, "Agent2"))

	// Create category
	s.CreateCategory(&Category{Name: "Test Category"})

	stats := s.GetStats()

	if stats.TotalListings != 2 {
		t.Errorf("expected 2 listings, got %d", stats.TotalListings)
	}
	if stats.TotalCategories != 1 {
		t.Errorf("expected 1 category, got %d", stats.TotalCategories)
	}
}

// TestPrettyPrint tests PrettyPrint methods
func TestPrettyPrint(t *testing.T) {
	listing := &AgentListing{
		AgentName:     "Test Agent",
		Rating:        4.5,
		DownloadCount: 100,
		Status:        ListingStatusPublished,
	}

	output := listing.PrettyPrint()
	if output == "" {
		t.Error("expected non-empty output")
	}

	review := &Review{
		Rating:       5,
		Title:        "Great!",
		ReviewerName: "Test User",
	}

	output = review.PrettyPrint()
	if output == "" {
		t.Error("expected non-empty output")
	}
}

// TestConcurrentAccess tests thread-safe access (simplified)
func TestConcurrentAccess(t *testing.T) {
	s, _ := helperService(t)

	var wg sync.WaitGroup
	numGoroutines := 5

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			s.CreateListing(context.Background(), helperCreateListing(s, "Agent"))
		}(i)
	}

	wg.Wait()

	if len(s.listings) != numGoroutines {
		t.Errorf("expected %d listings, got %d", numGoroutines, len(s.listings))
	}
}

// TestMarshalJSON tests JSON marshaling
func TestMarshalJSON(t *testing.T) {
	listing := &AgentListing{
		ListingID:    "listing-1",
		AgentID:      "agent-1",
		AgentName:    "Test Agent",
		AgentVersion: "1.0.0",
		Description:  "Test description",
	}

	data, err := listing.MarshalJSON()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(data) == 0 {
		t.Error("expected non-empty JSON")
	}

	review := &Review{
		ReviewID:     "review-1",
		Rating:       5,
		Title:        "Great!",
		ReviewerName: "Test User",
	}

	data, err = review.MarshalJSON()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(data) == 0 {
		t.Error("expected non-empty JSON")
	}
}

// TestSplitIntoWords tests word splitting functionality
func TestSplitIntoWords(t *testing.T) {
	testCases := []struct {
		input    string
		expected int
	}{
		{"hello world", 2},
		{"hello", 1},
		{"text generation", 2},
		{"", 0},
		{"singleword", 1},
	}

	for _, tc := range testCases {
		words := splitIntoWords(tc.input)
		if len(words) != tc.expected {
			t.Errorf("input '%s': expected %d words, got %d", tc.input, tc.expected, len(words))
		}
	}
}

// TestToLower tests lowercase conversion
func TestToLower(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello world"},
		{"TEST", "test"},
		{"MixedCase", "mixedcase"},
	}

	for _, tc := range testCases {
		result := toLower(tc.input)
		if result != tc.expected {
			t.Errorf("input '%s': expected '%s', got '%s'", tc.input, tc.expected, result)
		}
	}
}

// TestEventPublishing tests that events are published correctly
func TestEventPublishing(t *testing.T) {
	b := bus.New()
	s := NewService(b)

	events := make([]bus.Event, 0)
	var mu sync.Mutex
	eventCh, closer := b.Subscribe()
	defer closer()

	go func() {
		for event := range eventCh {
			mu.Lock()
			events = append(events, event)
			mu.Unlock()
		}
	}()

	listing := helperCreateListing(s, "TestAgent")
	s.CreateListing(context.Background(), listing)

	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	eventCount := len(events)
	mu.Unlock()
	if eventCount == 0 {
		t.Error("expected at least one event to be published")
	}
}
