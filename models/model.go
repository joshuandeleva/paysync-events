package models

import (
	"time"

	"gorm.io/gorm"
)

type RoleName string

const (
	RoleAdmin    RoleName = "ADMIN"
	RoleArtist   RoleName = "ARTIST"
	RolePromoter RoleName = "PROMOTER"
	RoleNormal   RoleName = "NORMAL"
)

type UserRole struct {
	gorm.Model
	RoleName RoleName `json:"role_name" gorm:"unique" validate:"required,oneof=ADMIN ARTIST PROMOTER NORMAL"`
}

type User struct {
	gorm.Model
	Username       string   `json:"username" gorm:"unique" validate:"required,min=3,max=50"`
	HashedPassword string   `json:"hashed_password" validate:"required,min=8"`
	Email          string   `json:"email" gorm:"unique" validate:"required,email"`
	RoleID         uint     `json:"role_id"`
	Role           UserRole `json:"role" gorm:"foreignKey:RoleID"`
}

type Promoter struct {
	gorm.Model
	UserID      uint   `json:"user_id" gorm:"unique"`
	User        User   `json:"user" gorm:"foreignKey:UserID"`
	CompanyName string `json:"company_name" validate:"required,min=3,max=50"`
}

type Artist struct {
	gorm.Model
	UserID     uint    `json:"user_id" gorm:"unique"`
	User       User    `json:"user" gorm:"foreignKey:UserID"`
	ArtistName string  `json:"artist_name" validate:"required,min=3,max=50"`
	BookingFee float64 `json:"booking_fee" validate:"required,min=1"`
}

type Venue struct {
	gorm.Model
	VenueName string   `json:"venue_name" validate:"required,min=3,max=50"`
	Location  string   `json:"location" validate:"required,min=3,max=50"`
	Capacity  int      `json:"capacity" validate:"required,min=1"`
	Events    []Event  `gorm:"many2many:event_venues;"`
}

type Event struct {
	gorm.Model
	VenueID    uint      `json:"venue_id"`
	Venue      Venue     `json:"venue" gorm:"foreignKey:VenueID"`
	EventName  string    `json:"event_name" validate:"required,min=3,max=50"`
	Date       string    `json:"date" validate:"required"`
	StartTime  string    `json:"start_time" validate:"required"`
	EndTime    string    `json:"end_time" validate:"required"`
	PromoterID uint      `json:"promoter_id"`
	Promoter   Promoter  `json:"promoter" gorm:"foreignKey:PromoterID"`
	Artists    []Artist  `gorm:"many2many:event_artists;"`
}

type EventArtist struct {
	gorm.Model
	EventID  uint   `json:"event_id"`
	Event    Event  `json:"event" gorm:"foreignKey:EventID"`
	ArtistID uint   `json:"artist_id"`
	Artist   Artist `json:"artist" gorm:"foreignKey:ArtistID"`
}

type EventVenue struct {
	gorm.Model
	EventID uint  `json:"event_id"`
	Event   Event `json:"event" gorm:"foreignKey:EventID"`
	VenueID uint  `json:"venue_id"`
	Venue   Venue `json:"venue" gorm:"foreignKey:VenueID"`
}

type Booking struct {
	gorm.Model
	EventID  uint      `json:"event_id"`
	Event    Event     `json:"event" gorm:"foreignKey:EventID"`
	UserID   uint      `json:"user_id"`
	User     User      `json:"user" gorm:"foreignKey:UserID"`
	ArtistID uint      `json:"artist_id"`
	Artist   Artist    `json:"artist" gorm:"foreignKey:ArtistID"`
	BookedAt time.Time `json:"booked_at" validate:"required"`
}

type Ticket struct {
	gorm.Model
	EventID    uint      `json:"event_id"`
	Event      Event     `json:"event" gorm:"foreignKey:EventID"`
	UserID     uint      `json:"user_id"`
	User       User      `json:"user" gorm:"foreignKey:UserID"`
	Quantity   int       `json:"quantity" validate:"required,min=1"`
	PurchaseAt time.Time `json:"purchase_at" validate:"required"`
	TicketType string    `json:"ticket_type" validate:"required,oneof=ADULT CHILDREN SENIOR"`
	Price      float64   `json:"price" validate:"required,min=1"`
}
