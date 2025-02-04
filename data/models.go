package data

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRole defines different user permission levels
type UserRole string

const (
	RoleUser      UserRole = "user"
	RoleCreator   UserRole = "creator"
	RoleAdmin     UserRole = "admin"
	RoleModerator UserRole = "moderator"
)

// UserStatus represents the current status of a user account
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusBanned    UserStatus = "banned"
)

// User represents a marketplace user
type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username      string             `bson:"username" json:"username" validate:"required,min=3,max=50"`
	Email         string             `bson:"email" json:"email" validate:"required,email"`
	PasswordHash  string             `bson:"password_hash" json:"-"`
	Role          UserRole           `bson:"role" json:"role"`
	Status        UserStatus         `bson:"status" json:"status"`
	Profile       UserProfile        `bson:"profile" json:"profile"`
	Credits       Credits            `bson:"credits" json:"credits"`
	Notifications []Notification     `bson:"notifications" json:"notifications"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

// UserProfile contains additional user information
type UserProfile struct {
	DisplayName string   `bson:"display_name" json:"display_name"`
	Avatar      string   `bson:"avatar" json:"avatar"`
	Bio         string   `bson:"bio" json:"bio"`
	Interests   []string `bson:"interests" json:"interests"`
	SocialLinks []string `bson:"social_links" json:"social_links"`
	Skills      []string `bson:"skills" json:"skills"`
	Links       []string `bson:"links" json:"links"`
}

// Credits manages user's digital marketplace credit system
type Credits struct {
	Balance      float64             `bson:"balance" json:"balance"`
	Transactions []CreditTransaction `bson:"transactions" json:"transactions"`
}

// CreditTransactionType defines different credit transaction types
type CreditTransactionType string

const (
	CreditTransactionPurchase   CreditTransactionType = "purchase"
	CreditTransactionRefund     CreditTransactionType = "refund"
	CreditTransactionDeposit    CreditTransactionType = "deposit"
	CreditTransactionWithdrawal CreditTransactionType = "withdrawal"
	CreditTransactionGift       CreditTransactionType = "gift"
)

// CreditTransaction represents credit-based transactions
type CreditTransaction struct {
	ID            primitive.ObjectID    `bson:"_id,omitempty" json:"id"`
	Type          CreditTransactionType `bson:"type" json:"type"`
	Amount        float64               `bson:"amount" json:"amount"`
	Status        string                `bson:"status" json:"status"`
	Description   string                `bson:"description" json:"description"`
	RelatedItemID primitive.ObjectID    `bson:"related_item_id,omitempty" json:"related_item_id"`
	Timestamp     time.Time             `bson:"timestamp" json:"timestamp"`
}

// ProductStatus represents the current status of a digital product
type ProductStatus string

const (
	ProductStatusDraft    ProductStatus = "draft"
	ProductStatusActive   ProductStatus = "active"
	ProductStatusInactive ProductStatus = "inactive"
	ProductStatusArchived ProductStatus = "archived"
)

// ProductCategory defines product categorization
type ProductCategory string

const (
	CategoryTemplate ProductCategory = "template"
	CategoryPlugin   ProductCategory = "plugin"
	CategoryAsset    ProductCategory = "asset"
	CategoryCourse   ProductCategory = "course"
	CategoryGuide    ProductCategory = "guide"
	CategorySource   ProductCategory = "source_code"
)

// Product represents a digital item in the marketplace
type Product struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatorID       primitive.ObjectID `bson:"creator_id" json:"creator_id"`
	Title           string             `bson:"title" json:"title" validate:"required,min=3,max=200"`
	Description     string             `bson:"description" json:"description"`
	Price           float64            `bson:"price" json:"price" validate:"gte=0"`
	DiscountedPrice float64            `bson:"discounted_price" json:"discounted_price"`
	Category        ProductCategory    `bson:"category" json:"category"`
	Status          ProductStatus      `bson:"status" json:"status"`
	Tags            []string           `bson:"tags" json:"tags"`
	Technologies    []string           `bson:"technologies" json:"technologies"`
	Images          []string           `bson:"images" json:"images"`
	Specifications  map[string]string  `bson:"specifications" json:"specifications"`
	Versions        []ProductVersion   `bson:"versions" json:"versions"`
	DownloadCount   int                `bson:"download_count" json:"download_count"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}

// ProductVersion tracks different versions of a digital product
type ProductVersion struct {
	Version       string    `bson:"version" json:"version"`
	ReleaseNotes  string    `bson:"release_notes" json:"release_notes"`
	DownloadURL   string    `bson:"download_url" json:"download_url"`
	Compatibility []string  `bson:"compatibility" json:"compatibility"`
	FileSize      int64     `bson:"file_size" json:"file_size"`
	ReleasedAt    time.Time `bson:"released_at" json:"released_at"`
}

// PurchaseStatus defines different states of a digital product purchase
type PurchaseStatus string

const (
	PurchaseStatusPending   PurchaseStatus = "pending"
	PurchaseStatusCompleted PurchaseStatus = "completed"
	PurchaseStatusRefunded  PurchaseStatus = "refunded"
	PurchaseStatusFailed    PurchaseStatus = "failed"
)

// Purchase represents a customer's digital product acquisition
type Purchase struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID         primitive.ObjectID `bson:"user_id" json:"user_id"`
	ProductID      primitive.ObjectID `bson:"product_id" json:"product_id"`
	ProductVersion string             `bson:"product_version" json:"product_version"`
	Price          float64            `bson:"price" json:"price"`
	Status         PurchaseStatus     `bson:"status" json:"status"`
	DownloadLink   string             `bson:"download_link" json:"download_link"`
	LicenseKey     string             `bson:"license_key" json:"license_key"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	ExpiresAt      time.Time          `bson:"expires_at" json:"expires_at"`
}

// ReviewStatus defines the status of a review
type ReviewStatus string

const (
	ReviewStatusPending  ReviewStatus = "pending"
	ReviewStatusApproved ReviewStatus = "approved"
	ReviewStatusRejected ReviewStatus = "rejected"
)

// Review represents a product review
type Review struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	ProductID primitive.ObjectID `bson:"product_id" json:"product_id"`
	Rating    int                `bson:"rating" json:"rating" validate:"gte=1,lte=5"`
	Comment   string             `bson:"comment" json:"comment"`
	Status    ReviewStatus       `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

// Notification represents user notifications
type Notification struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Type      string             `bson:"type" json:"type"`
	Message   string             `bson:"message" json:"message"`
	RelatedID primitive.ObjectID `bson:"related_id,omitempty" json:"related_id"`
	IsRead    bool               `bson:"is_read" json:"is_read"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

// Dispute represents a dispute between buyer and creator
type Dispute struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PurchaseID  primitive.ObjectID `bson:"purchase_id" json:"purchase_id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	CreatorID   primitive.ObjectID `bson:"creator_id" json:"creator_id"`
	Reason      string             `bson:"reason" json:"reason"`
	Description string             `bson:"description" json:"description"`
	Status      string             `bson:"status" json:"status"`
	Resolution  string             `bson:"resolution" json:"resolution"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// MarketplaceSettings represents global marketplace configuration
type MarketplaceSettings struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CommissionRate    float64            `bson:"commission_rate" json:"commission_rate"`
	MinimumWithdrawal float64            `bson:"minimum_withdrawal" json:"minimum_withdrawal"`
	CreatorPayoutRate float64            `bson:"creator_payout_rate" json:"creator_payout_rate"`
	RefundPeriodDays  int                `bson:"refund_period_days" json:"refund_period_days"`
	MaintenanceMode   bool               `bson:"maintenance_mode" json:"maintenance_mode"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at"`
}

// CreatorAnalytics provides comprehensive performance tracking for creators
// CreatorAnalytics represents comprehensive analytics and performance metrics for a creator in the marketplace.
// It tracks various aspects including:
//   - Overall performance metrics (revenue, sales, ratings)
//   - Product-specific performance data
//   - Time-based revenue analytics (daily, weekly, monthly)
//   - Audience insights and demographics
//   - Engagement metrics including review statistics
//   - Payout and financial information
//
// The analytics data is stored in MongoDB and can be serialized to/from JSON.
// All monetary values are stored as float64 in the platform's base currency.
// Time-based fields use time.Time and are stored in UTC.
type CreatorAnalytics struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatorID primitive.ObjectID `bson:"creator_id" json:"creator_id"`

	// Overall Performance Metrics
	TotalRevenue      float64 `bson:"total_revenue" json:"total_revenue"`
	TotalProductsSold int     `bson:"total_products_sold" json:"total_products_sold"`
	AverageRating     float64 `bson:"average_rating" json:"average_rating"`
	TotalViews        int     `bson:"total_views" json:"total_views"`

	// Product-Level Performance
	ProductPerformance []ProductPerformance `bson:"product_performance" json:"product_performance"`

	// Time-Based Analytics
	DailyRevenue   []DailyMetric `bson:"daily_revenue" json:"daily_revenue"`
	WeeklyRevenue  []DailyMetric `bson:"weekly_revenue" json:"weekly_revenue"`
	MonthlyRevenue []DailyMetric `bson:"monthly_revenue" json:"monthly_revenue"`

	// Audience Insights
	AudienceBreakdown AudienceInsights `bson:"audience_insights" json:"audience_insights"`

	// Engagement Metrics
	TotalReviews    int `bson:"total_reviews" json:"total_reviews"`
	PositiveReviews int `bson:"positive_reviews" json:"positive_reviews"`
	NeutralReviews  int `bson:"neutral_reviews" json:"neutral_reviews"`
	NegativeReviews int `bson:"negative_reviews" json:"negative_reviews"`

	// Payout Information
	PendingPayout     float64   `bson:"pending_payout" json:"pending_payout"`
	LastPayoutDate    time.Time `bson:"last_payout_date" json:"last_payout_date"`
	TotalPayoutAmount float64   `bson:"total_payout_amount" json:"total_payout_amount"`

	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// ProductPerformance tracks individual product metrics
type ProductPerformance struct {
	ProductID      primitive.ObjectID `bson:"product_id" json:"product_id"`
	ProductTitle   string             `bson:"product_title" json:"product_title"`
	TotalRevenue   float64            `bson:"total_revenue" json:"total_revenue"`
	UnitsSold      int                `bson:"units_sold" json:"units_sold"`
	AverageRating  float64            `bson:"average_rating" json:"average_rating"`
	TotalViews     int                `bson:"total_views" json:"total_views"`
	ConversionRate float64            `bson:"conversion_rate" json:"conversion_rate"`
	RefundRate     float64            `bson:"refund_rate" json:"refund_rate"`
}

// DailyMetric represents time-series financial data
type DailyMetric struct {
	Date  time.Time `bson:"date" json:"date"`
	Value float64   `bson:"value" json:"value"`
}

// AudienceInsights provides demographic and behavioral analytics
type AudienceInsights struct {
	GeographicDistribution map[string]int `bson:"geographic_distribution" json:"geographic_distribution"`
	AgeGroups              map[string]int `bson:"age_groups" json:"age_groups"`
	PurchaseFrequency      map[string]int `bson:"purchase_frequency" json:"purchase_frequency"`
	PrimaryInterests       []string       `bson:"primary_interests" json:"primary_interests"`
}

// CreatorPayout represents detailed payout information
type CreatorPayout struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatorID       primitive.ObjectID `bson:"creator_id" json:"creator_id"`
	PeriodStart     time.Time          `bson:"period_start" json:"period_start"`
	PeriodEnd       time.Time          `bson:"period_end" json:"period_end"`
	TotalRevenue    float64            `bson:"total_revenue" json:"total_revenue"`
	MarketplaceFee  float64            `bson:"marketplace_fee" json:"marketplace_fee"`
	NetPayout       float64            `bson:"net_payout" json:"net_payout"`
	PayoutMethod    string             `bson:"payout_method" json:"payout_method"`
	Status          string             `bson:"status" json:"status"`
	PayoutDate      time.Time          `bson:"payout_date" json:"payout_date"`
	PurchaseDetails []PurchaseDetail   `bson:"purchase_details" json:"purchase_details"`
}

// PurchaseDetail provides granular information about individual purchases
type PurchaseDetail struct {
	ProductID    primitive.ObjectID `bson:"product_id" json:"product_id"`
	ProductTitle string             `bson:"product_title" json:"product_title"`
	Quantity     int                `bson:"quantity" json:"quantity"`
	UnitPrice    float64            `bson:"unit_price" json:"unit_price"`
	TotalRevenue float64            `bson:"total_revenue" json:"total_revenue"`
	PurchaseDate time.Time          `bson:"purchase_date" json:"purchase_date"`
}

// CreatorApplicationTracking for marketplace onboarding
type CreatorApplicationTracking struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID             primitive.ObjectID `bson:"user_id" json:"user_id"`
	ApplicationStatus  string             `bson:"application_status" json:"application_status"`
	SubmissionDate     time.Time          `bson:"submission_date" json:"submission_date"`
	ReviewDate         time.Time          `bson:"review_date" json:"review_date"`
	DocumentsSubmitted []string           `bson:"documents_submitted" json:"documents_submitted"`
	Notes              string             `bson:"notes" json:"notes"`
	ApprovedAt         time.Time          `bson:"approved_at" json:"approved_at"`
	RejectedAt         time.Time          `bson:"rejected_at" json:"rejected_at"`
}
