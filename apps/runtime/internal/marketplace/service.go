package marketplace

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"pryx-core/internal/bus"
)

// AgentListing represents an agent listing in the marketplace
type AgentListing struct {
	ListingID        string                 `json:"listing_id"`
	AgentID          string                 `json:"agent_id"`
	AgentName        string                 `json:"agent_name"`
	AgentVersion     string                 `json:"agent_version"`
	Description      string                 `json:"description"`
	LongDescription  string                 `json:"long_description,omitempty"`
	Author           string                 `json:"author"`
	AuthorID         string                 `json:"author_id"`
	Category         string                 `json:"category"`
	Tags             []string               `json:"tags"`
	Capabilities     []string               `json:"capabilities"`
	PriceModel       string                 `json:"price_model"`
	Price            float64                `json:"price"`
	Currency         string                 `json:"currency"`
	License          string                 `json:"license"`
	RepositoryURL    string                 `json:"repository_url,omitempty"`
	DocumentationURL string                 `json:"documentation_url,omitempty"`
	IconURL          string                 `json:"icon_url,omitempty"`
	Screenshots      []string               `json:"screenshots,omitempty"`
	Rating           float64                `json:"rating"`
	ReviewCount      int                    `json:"review_count"`
	DownloadCount    int                    `json:"download_count"`
	InstallCount     int                    `json:"install_count"`
	Status           ListingStatus          `json:"status"`
	Featured         bool                   `json:"featured"`
	Verified         bool                   `json:"verified"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	PublishedAt      *time.Time             `json:"published_at,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// ListingStatus represents the status of a listing
type ListingStatus string

const (
	ListingStatusDraft     ListingStatus = "draft"
	ListingStatusPending   ListingStatus = "pending"
	ListingStatusPublished ListingStatus = "published"
	ListingStatusSuspended ListingStatus = "suspended"
	ListingStatusRejected  ListingStatus = "rejected"
)

// Review represents a review for an agent listing
type Review struct {
	ReviewID     string          `json:"review_id"`
	ListingID    string          `json:"listing_id"`
	AgentID      string          `json:"agent_id"`
	ReviewerID   string          `json:"reviewer_id"`
	ReviewerName string          `json:"reviewer_name"`
	Rating       int             `json:"rating"`
	Title        string          `json:"title"`
	Content      string          `json:"content"`
	Pros         []string        `json:"pros,omitempty"`
	Cons         []string        `json:"cons,omitempty"`
	Verified     bool            `json:"verified"`
	HelpfulCount int             `json:"helpful_count"`
	Response     *ReviewResponse `json:"response,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// ReviewResponse represents an author's response to a review
type ReviewResponse struct {
	ResponseID string    `json:"response_id"`
	ReviewID   string    `json:"review_id"`
	AuthorID   string    `json:"author_id"`
	AuthorName string    `json:"author_name"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

// SearchQuery represents a marketplace search query
type SearchQuery struct {
	QueryID     string                 `json:"query_id"`
	RequesterID string                 `json:"requester_id"`
	Keywords    string                 `json:"keywords"`
	Category    string                 `json:"category,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	MinRating   float64                `json:"min_rating,omitempty"`
	MaxPrice    float64                `json:"max_price,omitempty"`
	PriceModel  string                 `json:"price_model,omitempty"`
	SortBy      string                 `json:"sort_by"`
	SortOrder   string                 `json:"sort_order"`
	Page        int                    `json:"page"`
	PageSize    int                    `json:"page_size"`
	Filters     map[string]interface{} `json:"filters,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
}

// SearchResult represents search results
type SearchResult struct {
	QueryID    string          `json:"query_id"`
	Listings   []*AgentListing `json:"listings"`
	TotalCount int             `json:"total_count"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
	Facets     SearchFacets    `json:"facets"`
	ReturnedAt time.Time       `json:"returned_at"`
}

// SearchFacets represents aggregated search facets
type SearchFacets struct {
	Categories   map[string]int `json:"categories"`
	Tags         map[string]int `json:"tags"`
	PriceModels  map[string]int `json:"price_models"`
	RatingRanges map[string]int `json:"rating_ranges"`
}

// Category represents a marketplace category
type Category struct {
	CategoryID   string                 `json:"category_id"`
	Name         string                 `json:"name"`
	Slug         string                 `json:"slug"`
	Description  string                 `json:"description"`
	IconURL      string                 `json:"icon_url,omitempty"`
	ParentID     string                 `json:"parent_id,omitempty"`
	ListingCount int                    `json:"listing_count"`
	Featured     []string               `json:"featured"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// ServiceStats represents marketplace service statistics
type ServiceStats struct {
	TotalListings     int `json:"total_listings"`
	PublishedListings int `json:"published_listings"`
	TotalReviews      int `json:"total_reviews"`
	TotalCategories   int `json:"total_categories"`
	TotalDownloads    int `json:"total_downloads"`
	TotalInstalls     int `json:"total_installs"`
}

// Service manages the agent marketplace
type Service struct {
	mu          sync.RWMutex
	bus         *bus.Bus
	listings    map[string]*AgentListing
	reviews     map[string]*Review
	categories  map[string]*Category
	searchIndex map[string][]string // keyword -> listing IDs
	installLog  []InstallEvent
}

// InstallEvent represents an installation event
type InstallEvent struct {
	EventID   string    `json:"event_id"`
	ListingID string    `json:"listing_id"`
	AgentID   string    `json:"agent_id"`
	UserID    string    `json:"user_id"`
	Version   string    `json:"version"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
}

// NewService creates a new marketplace service
func NewService(b *bus.Bus) *Service {
	return &Service{
		bus:         b,
		listings:    make(map[string]*AgentListing),
		reviews:     make(map[string]*Review),
		categories:  make(map[string]*Category),
		searchIndex: make(map[string][]string),
		installLog:  make([]InstallEvent, 0),
	}
}

// CreateListing creates a new agent listing
func (s *Service) CreateListing(ctx context.Context, listing *AgentListing) error {
	listing.ListingID = uuid.New().String()
	listing.CreatedAt = time.Now().UTC()
	listing.UpdatedAt = time.Now().UTC()
	listing.Status = ListingStatusDraft
	listing.ReviewCount = 0
	listing.DownloadCount = 0
	listing.InstallCount = 0

	s.mu.Lock()
	s.listings[listing.ListingID] = listing
	s.mu.Unlock()

	// Index the listing for search
	s.indexListing(listing)

	// Publish event
	s.bus.Publish(bus.NewEvent("marketplace.listing.created", "", map[string]interface{}{
		"listing_id": listing.ListingID,
		"agent_id":   listing.AgentID,
		"agent_name": listing.AgentName,
		"category":   listing.Category,
	}))

	return nil
}

// indexListing adds a listing to the search index
func (s *Service) indexListing(listing *AgentListing) {
	s.mu.Lock()
	defer s.mu.Unlock()

	keywords := []string{
		listing.AgentName,
		listing.Description,
		listing.Category,
	}
	keywords = append(keywords, listing.Tags...)
	keywords = append(keywords, listing.Capabilities...)

	for _, keyword := range keywords {
		keyword = toLower(keyword)
		s.searchIndex[keyword] = append(s.searchIndex[keyword], listing.ListingID)
	}
}

// toLower converts string to lowercase
func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			result[i] = s[i] + 32
		} else {
			result[i] = s[i]
		}
	}
	return string(result)
}

// GetListing retrieves a listing by ID
func (s *Service) GetListing(listingID string) (*AgentListing, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	listing, exists := s.listings[listingID]
	if !exists {
		return nil, fmt.Errorf("listing not found: %s", listingID)
	}

	return listing, nil
}

// GetListingByAgentID retrieves a listing by agent ID
func (s *Service) GetListingByAgentID(agentID string) (*AgentListing, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, listing := range s.listings {
		if listing.AgentID == agentID {
			return listing, nil
		}
	}

	return nil, fmt.Errorf("listing not found for agent: %s", agentID)
}

// UpdateListing updates an existing listing
func (s *Service) UpdateListing(listingID string, updates map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	listing, exists := s.listings[listingID]
	if !exists {
		return fmt.Errorf("listing not found: %s", listingID)
	}

	// Apply updates
	if name, ok := updates["name"].(string); ok {
		listing.AgentName = name
	}
	if desc, ok := updates["description"].(string); ok {
		listing.Description = desc
	}
	if category, ok := updates["category"].(string); ok {
		listing.Category = category
	}
	if tags, ok := updates["tags"].([]string); ok {
		listing.Tags = tags
	}
	if price, ok := updates["price"].(float64); ok {
		listing.Price = price
	}
	if status, ok := updates["status"].(ListingStatus); ok {
		listing.Status = status
		if status == ListingStatusPublished && listing.PublishedAt == nil {
			now := time.Now().UTC()
			listing.PublishedAt = &now
		}
	}

	listing.UpdatedAt = time.Now().UTC()

	// Re-index (release lock first)
	s.mu.Unlock()
	s.indexListing(listing)
	s.mu.Lock()

	return nil
}

// PublishListing publishes a listing
func (s *Service) PublishListing(listingID string) error {
	return s.UpdateListing(listingID, map[string]interface{}{
		"status": ListingStatusPublished,
	})
}

// SuspendListing suspends a listing
func (s *Service) SuspendListing(listingID string, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	listing, exists := s.listings[listingID]
	if !exists {
		return fmt.Errorf("listing not found: %s", listingID)
	}

	listing.Status = ListingStatusSuspended
	listing.UpdatedAt = time.Now().UTC()
	if listing.Metadata == nil {
		listing.Metadata = make(map[string]interface{})
	}
	listing.Metadata["suspend_reason"] = reason

	// Publish event
	s.bus.Publish(bus.NewEvent("marketplace.listing.suspended", "", map[string]interface{}{
		"listing_id": listingID,
		"reason":     reason,
	}))

	return nil
}

// GetFeaturedListings returns featured listings
func (s *Service) GetFeaturedListings(limit int) []*AgentListing {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var featured []*AgentListing
	for _, listing := range s.listings {
		if listing.Featured && listing.Status == ListingStatusPublished {
			featured = append(featured, listing)
			if len(featured) >= limit {
				break
			}
		}
	}

	return featured
}

// GetListingsByCategory returns listings in a category
func (s *Service) GetListingsByCategory(category string, limit int) []*AgentListing {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var listings []*AgentListing
	for _, listing := range s.listings {
		if listing.Category == category && listing.Status == ListingStatusPublished {
			listings = append(listings, listing)
			if len(listings) >= limit {
				break
			}
		}
	}

	return listings
}

// SearchListings searches for listings
func (s *Service) SearchListings(ctx context.Context, query *SearchQuery) (*SearchResult, error) {
	query.QueryID = uuid.New().String()
	query.CreatedAt = time.Now().UTC()

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}
	if query.PageSize > 100 {
		query.PageSize = 100
	}

	s.mu.RLock()
	matchingIDs := s.searchListings(query)
	s.mu.RUnlock()

	// Filter and sort listings
	var listings []*AgentListing
	s.mu.RLock()
	for _, id := range matchingIDs {
		if listing, ok := s.listings[id]; ok {
			if s.matchesFilters(listing, query) {
				listings = append(listings, listing)
			}
		}
	}
	s.mu.RUnlock()

	// Sort results
	s.sortListings(listings, query.SortBy, query.SortOrder)

	// Calculate pagination
	totalCount := len(listings)
	totalPages := (totalCount + query.PageSize - 1) / query.PageSize
	start := (query.Page - 1) * query.PageSize
	end := start + query.PageSize
	if end > totalCount {
		end = totalCount
	}

	paginatedListings := listings
	if start < totalCount {
		if end > totalCount {
			end = totalCount
		}
		paginatedListings = listings[start:end]
	}

	// Generate facets
	facets := s.generateFacets(listings)

	result := &SearchResult{
		QueryID:    query.QueryID,
		Listings:   paginatedListings,
		TotalCount: totalCount,
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: totalPages,
		Facets:     facets,
		ReturnedAt: time.Now().UTC(),
	}

	// Publish event
	s.bus.Publish(bus.NewEvent("marketplace.search.completed", "", map[string]interface{}{
		"query_id":     query.QueryID,
		"requester_id": query.RequesterID,
		"keywords":     query.Keywords,
		"result_count": len(paginatedListings),
	}))

	return result, nil
}

// searchListings performs keyword search
func (s *Service) searchListings(query *SearchQuery) []string {
	if query.Keywords == "" {
		// Return all published listings
		var allIDs []string
		s.mu.RLock()
		for id, listing := range s.listings {
			if listing.Status == ListingStatusPublished {
				allIDs = append(allIDs, id)
			}
		}
		s.mu.RUnlock()
		return allIDs
	}

	keywords := toLower(query.Keywords)
	var matchingIDs []string
	keywordMap := make(map[string]bool)

	for _, keyword := range []string{keywords} {
		parts := splitIntoWords(keyword)
		for _, part := range parts {
			keywordMap[part] = true
		}
	}

	s.mu.RLock()
	for keyword, ids := range s.searchIndex {
		for _, part := range splitIntoWords(keyword) {
			if keywordMap[part] {
				matchingIDs = append(matchingIDs, ids...)
			}
		}
	}
	s.mu.RUnlock()

	return matchingIDs
}

// splitIntoWords splits text into words
func splitIntoWords(text string) []string {
	var words []string
	word := make([]byte, 0, len(text))
	for i := 0; i < len(text); i++ {
		if text[i] >= 'a' && text[i] <= 'z' || text[i] >= 'A' && text[i] <= 'Z' || text[i] >= '0' && text[i] <= '9' {
			word = append(word, text[i])
		} else {
			if len(word) > 0 {
				words = append(words, string(word))
				word = make([]byte, 0, len(text))
			}
		}
	}
	if len(word) > 0 {
		words = append(words, string(word))
	}
	return words
}

// matchesFilters checks if a listing matches query filters
func (s *Service) matchesFilters(listing *AgentListing, query *SearchQuery) bool {
	if listing.Status != ListingStatusPublished {
		return false
	}

	if query.Category != "" && listing.Category != query.Category {
		return false
	}

	if query.MinRating > 0 && listing.Rating < query.MinRating {
		return false
	}

	if query.MaxPrice > 0 && listing.Price > query.MaxPrice {
		return false
	}

	if query.PriceModel != "" && listing.PriceModel != query.PriceModel {
		return false
	}

	if len(query.Tags) > 0 {
		hasAllTags := true
		for _, tag := range query.Tags {
			found := false
			for _, listingTag := range listing.Tags {
				if listingTag == tag {
					found = true
					break
				}
			}
			if !found {
				hasAllTags = false
				break
			}
		}
		if !hasAllTags {
			return false
		}
	}

	return true
}

// sortListings sorts listings by the specified criteria
func (s *Service) sortListings(listings []*AgentListing, sortBy, sortOrder string) {
	reverse := sortOrder == "desc"

	switch sortBy {
	case "rating":
		if reverse {
			sortByRatingDesc(listings)
		} else {
			sortByRatingAsc(listings)
		}
	case "downloads":
		if reverse {
			sortByDownloadsDesc(listings)
		} else {
			sortByDownloadsAsc(listings)
		}
	case "price":
		if reverse {
			sortByPriceDesc(listings)
		} else {
			sortByPriceAsc(listings)
		}
	case "newest":
		sortByNewest(listings)
	case "name":
		sortByName(listings)
	}
}

func sortByRatingDesc(listings []*AgentListing) []*AgentListing {
	for i := 0; i < len(listings)-1; i++ {
		for j := i + 1; j < len(listings); j++ {
			if listings[i].Rating < listings[j].Rating {
				listings[i], listings[j] = listings[j], listings[i]
			}
		}
	}
	return listings
}

func sortByRatingAsc(listings []*AgentListing) []*AgentListing {
	for i := 0; i < len(listings)-1; i++ {
		for j := i + 1; j < len(listings); j++ {
			if listings[i].Rating > listings[j].Rating {
				listings[i], listings[j] = listings[j], listings[i]
			}
		}
	}
	return listings
}

func sortByDownloadsDesc(listings []*AgentListing) []*AgentListing {
	for i := 0; i < len(listings)-1; i++ {
		for j := i + 1; j < len(listings); j++ {
			if listings[i].DownloadCount < listings[j].DownloadCount {
				listings[i], listings[j] = listings[j], listings[i]
			}
		}
	}
	return listings
}

func sortByDownloadsAsc(listings []*AgentListing) []*AgentListing {
	for i := 0; i < len(listings)-1; i++ {
		for j := i + 1; j < len(listings); j++ {
			if listings[i].DownloadCount > listings[j].DownloadCount {
				listings[i], listings[j] = listings[j], listings[i]
			}
		}
	}
	return listings
}

func sortByPriceDesc(listings []*AgentListing) []*AgentListing {
	for i := 0; i < len(listings)-1; i++ {
		for j := i + 1; j < len(listings); j++ {
			if listings[i].Price < listings[j].Price {
				listings[i], listings[j] = listings[j], listings[i]
			}
		}
	}
	return listings
}

func sortByPriceAsc(listings []*AgentListing) []*AgentListing {
	for i := 0; i < len(listings)-1; i++ {
		for j := i + 1; j < len(listings); j++ {
			if listings[i].Price > listings[j].Price {
				listings[i], listings[j] = listings[j], listings[i]
			}
		}
	}
	return listings
}

func sortByNewest(listings []*AgentListing) []*AgentListing {
	for i := 0; i < len(listings)-1; i++ {
		for j := i + 1; j < len(listings); j++ {
			if listings[i].PublishedAt != nil && listings[j].PublishedAt != nil {
				if listings[i].PublishedAt.Before(*listings[j].PublishedAt) {
					listings[i], listings[j] = listings[j], listings[i]
				}
			}
		}
	}
	return listings
}

func sortByName(listings []*AgentListing) []*AgentListing {
	for i := 0; i < len(listings)-1; i++ {
		for j := i + 1; j < len(listings); j++ {
			if listings[i].AgentName > listings[j].AgentName {
				listings[i], listings[j] = listings[j], listings[i]
			}
		}
	}
	return listings
}

// generateFacets generates search facets from listings
func (s *Service) generateFacets(listings []*AgentListing) SearchFacets {
	facets := SearchFacets{
		Categories:  make(map[string]int),
		Tags:        make(map[string]int),
		PriceModels: make(map[string]int),
	}

	for _, listing := range listings {
		facets.Categories[listing.Category]++
		for _, tag := range listing.Tags {
			facets.Tags[tag]++
		}
		facets.PriceModels[listing.PriceModel]++
	}

	return facets
}

// AddReview adds a review to a listing
func (s *Service) AddReview(ctx context.Context, review *Review) error {
	review.ReviewID = uuid.New().String()
	review.CreatedAt = time.Now().UTC()
	review.UpdatedAt = time.Now().UTC()
	review.HelpfulCount = 0

	s.mu.Lock()
	s.reviews[review.ReviewID] = review
	s.mu.Unlock()

	// Update listing review count and rating
	s.updateListingRating(review.ListingID)

	// Publish event
	s.bus.Publish(bus.NewEvent("marketplace.review.added", "", map[string]interface{}{
		"review_id":   review.ReviewID,
		"listing_id":  review.ListingID,
		"rating":      review.Rating,
		"reviewer_id": review.ReviewerID,
	}))

	return nil
}

// updateListingRating recalculates listing rating based on reviews
func (s *Service) updateListingRating(listingID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var totalRating float64
	count := 0

	for _, review := range s.reviews {
		if review.ListingID == listingID {
			totalRating += float64(review.Rating)
			count++
		}
	}

	if listing, ok := s.listings[listingID]; ok {
		listing.ReviewCount = count
		if count > 0 {
			listing.Rating = totalRating / float64(count)
		}
	}
}

// GetReviews retrieves reviews for a listing
func (s *Service) GetReviews(listingID string, limit int) []*Review {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var reviews []*Review
	for _, review := range s.reviews {
		if review.ListingID == listingID {
			reviews = append(reviews, review)
			if len(reviews) >= limit {
				break
			}
		}
	}

	return reviews
}

// MarkReviewHelpful increments the helpful count for a review
func (s *Service) MarkReviewHelpful(reviewID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	review, exists := s.reviews[reviewID]
	if !exists {
		return fmt.Errorf("review not found: %s", reviewID)
	}

	review.HelpfulCount++

	return nil
}

// AddReviewResponse adds an author's response to a review
func (s *Service) AddReviewResponse(reviewID string, response *ReviewResponse) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	review, exists := s.reviews[reviewID]
	if !exists {
		return fmt.Errorf("review not found: %s", reviewID)
	}

	response.ResponseID = uuid.New().String()
	response.ReviewID = reviewID
	response.CreatedAt = time.Now().UTC()

	review.Response = response
	review.UpdatedAt = time.Now().UTC()

	return nil
}

// CreateCategory creates a new category
func (s *Service) CreateCategory(category *Category) error {
	category.CategoryID = uuid.New().String()
	category.CreatedAt = time.Now().UTC()
	category.UpdatedAt = time.Now().UTC()

	s.mu.Lock()
	s.categories[category.CategoryID] = category
	s.mu.Unlock()

	return nil
}

// GetCategory retrieves a category
func (s *Service) GetCategory(categoryID string) (*Category, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	category, exists := s.categories[categoryID]
	if !exists {
		return nil, fmt.Errorf("category not found: %s", categoryID)
	}

	return category, nil
}

// GetAllCategories returns all categories
func (s *Service) GetAllCategories() []*Category {
	s.mu.RLock()
	defer s.mu.RUnlock()

	categories := make([]*Category, 0, len(s.categories))
	for _, category := range s.categories {
		categories = append(categories, category)
	}

	return categories
}

// RecordInstall records an installation event
func (s *Service) RecordInstall(listingID, userID, version, source string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	event := InstallEvent{
		EventID:   uuid.New().String(),
		ListingID: listingID,
		UserID:    userID,
		Version:   version,
		Source:    source,
		CreatedAt: time.Now().UTC(),
	}

	s.installLog = append(s.installLog, event)

	if listing, ok := s.listings[listingID]; ok {
		listing.InstallCount++
	}

	// Publish event
	s.bus.Publish(bus.NewEvent("marketplace.install.recorded", "", map[string]interface{}{
		"listing_id": listingID,
		"user_id":    userID,
		"version":    version,
		"source":     source,
	}))

	return nil
}

// RecordDownload records a download event
func (s *Service) RecordDownload(listingID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if listing, ok := s.listings[listingID]; ok {
		listing.DownloadCount++
	}

	return nil
}

// GetStats returns marketplace statistics
func (s *Service) GetStats() ServiceStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := ServiceStats{
		TotalListings:   len(s.listings),
		TotalReviews:    len(s.reviews),
		TotalCategories: len(s.categories),
	}

	for _, listing := range s.listings {
		if listing.Status == ListingStatusPublished {
			stats.PublishedListings++
		}
		stats.TotalDownloads += listing.DownloadCount
		stats.TotalInstalls += listing.InstallCount
	}

	return stats
}

// PrettyPrint prints listing info
func (l *AgentListing) PrettyPrint() string {
	return fmt.Sprintf("Listing{Name: %s, Rating: %.2f, Downloads: %d, Status: %s}",
		l.AgentName, l.Rating, l.DownloadCount, l.Status)
}

// PrettyPrint prints review info
func (r *Review) PrettyPrint() string {
	return fmt.Sprintf("Review{Rating: %d, Title: %s, Reviewer: %s}",
		r.Rating, r.Title, r.ReviewerName)
}

// MarshalJSON for AgentListing
func (l *AgentListing) MarshalJSON() ([]byte, error) {
	type Alias AgentListing
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(l),
	})
}

// MarshalJSON for Review
func (r *Review) MarshalJSON() ([]byte, error) {
	type Alias Review
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	})
}
