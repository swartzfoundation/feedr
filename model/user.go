package model

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
)

// UserContextKey is the key used to store the user in the context.
const UserContextKey ContextKey = "user"
const UserTableName = "users"

type User struct {
	// ID is the unique ID for the user.
	// required: true
	ID string `json:"id" gorm:"primaryKey"`

	// FirstName is the first name of the user.
	// required: false
	FirstName string `json:"first_name" validate:"omitempty"`

	// LastName is the last name of the user.
	// required: false
	LastName string `json:"last_name" validate:"omitempty"`

	// Username is the unique username of the user.
	// required: true
	Username string `json:"username" gorm:"uniqueIndex:idx_users_username; not null; default:null;" validate:"required"`

	// Password is the password of the user.
	// required: false
	PasswordHash string `json:"-" validate:"required"`

	// PasswordHashType is the type of hash used to store the password.
	// required: false
	PasswordHashType string `json:"-" gorm:"not null;default:'bcrypt';"`

	// Email is the email of the user.
	// required: true
	Email string `json:"email" gorm:"uniqueIndex:idx_users_email; not null; default:null;" validate:"required,email"`

	// EmailVerified is true if the user has verified their email.
	// required: true
	EmailVerified bool `json:"email_verified,omitempty"`

	isAnon bool `json:"-" gorm:"-"`

	// IsActive is true if the user is active.
	// required: true
	IsActive bool `json:"is_active" gorm:"default:true"`

	// IsAdmin is true if the user is an admin.
	// required: true
	IsAdmin bool `json:"is_admin" gorm:"default:false"`

	// IsSuperAdmin is true if the user is a superadmin.
	// required: true
	IsSuperAdmin bool `json:"is_super_admin" gorm:"default:false"`

	// AllowAnalytics is true if the user allows analytics.
	// required: true
	AllowAnalytics bool `json:"allow_analytics,omitempty" gorm:"default:true"`

	// LoginType is the type of login used for the user.
	// required: true
	LoginType string `json:"login_type,omitempty" gorm:"default:'email'"`

	// LastPasswordUpdateAt is the unix timestamp of the last password update.
	// required: true
	LastPasswordUpdateAt int64 `json:"last_password_update,omitempty" gorm:"default:0"`

	// FailedAttempts is the number of failed login attempts.
	// required: true
	FailedAttempts int `json:"failed_attempts,omitempty" gorm:"default:0"`

	// Locale is the locale of the user.
	// required: true
	Locale string `json:"locale" gorm:"default:'en'"`

	// LastActivityAt is the unix timestamp of the last activity.
	// required: true
	LastActivityAt int64 `json:"last_activity_at,omitempty"`

	// LastLoginAt is the unix timestamp of the last login.
	// required: true
	LastLoginAt int64 `json:"last_login_at,omitempty"`
	// LastLoginIP is the IP address of the last login.
	// required: true
	LastLoginIP string `json:"last_login_ip,omitempty"`

	// CreatedAt is the unix timestamp of the creation date.
	// required: true
	CreatedAt int64 `json:"created_at" gorm:"autoCreateTime"`

	//UpdatedAt is the unix timestamp of the last update.
	// required: true
	UpdatedAt int64 `json:"updated_at"`

	// DeletedAt is the unix timestamp of the deletion date.
	// required: true
	DeletedAt int64 `json:"-" gorm:"index"`
}

func (u *User) TableName() string {
	return UserTableName
}

type UserPatch struct {
	// FirstName is the first name of the user.
	// required: false
	FirstName *string `json:"first_name" validate:"omitempty"`

	// LastName is the last name of the user.
	// required: false
	LastName *string `json:"last_name" validate:"omitempty"`

	// Username is the unique username of the user.
	// required: true
	Username *string `json:"username" validate:"required"`

	// Email is the email of the user.
	// required: false
	Email *string `json:"email" validate:"required,email"`

	// AllowAnalytics is true if the user allows analytics.
	// required: false
	AllowAnalytics *bool `json:"allow_analytics,omitempty"`
}

// PreSave will set the ID and Username if missing.  It will also fill
// in the CreateAt, UpdateAt times.  It will also hash the password.  It should
// be run before saving the user to the db.
func (u *User) PreSave() {
	if u.Username == "" {
		u.Username = strings.ToLower(NewID())
	}
	u.Email = strings.ToLower(u.Email)
	u.Email = SanitizeUnicode(u.Email)
	u.Username = SanitizeUnicode(u.Username)

}

func (u *User) IsValid() error {
	if u.isAnon {
		return errors.New("user is anon")
	}
	return nil
}

// BeforeCreate will set the ID and Username if missing.  It will also fill
func (u *User) BeforeCreate(db *gorm.DB) error {
	u.ID = NewID()
	u.PreSave()
	return u.IsValid()
}

func (u *User) BeforeUpdate(db *gorm.DB) error {
	u.PreSave()
	return u.IsValid()
}

func (u *User) CheckPassword(password string) bool {
	return CheckPassword(u.PasswordHash, password)
}
func (u *User) SetPassword(password string) error {
	hash, err := HashPassword(password)
	if err != nil {
		return err
	}
	u.PasswordHash = hash
	return nil
}

func GetUserByID(ctx context.Context, id string) (*User, error) {
	var u User
	if result := db.WithContext(ctx).First(&u, "id = ?", id); result.Error != nil {
		return nil, result.Error
	}

	return &u, nil
}
func GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	if result := db.WithContext(ctx).First(&u, "email = ?", email); result.Error != nil {
		return nil, result.Error
	}
	return &u, nil
}

func WithUserContext(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, UserContextKey, u)
}

func (up *UserPatch) Patch(user *User) *User {
	if up.FirstName != nil {
		user.FirstName = *up.FirstName
	}
	if up.LastName != nil {
		user.LastName = *up.LastName
	}
	if up.Username != nil {
		user.Username = *up.Username
	}
	if up.Email != nil {
		user.Email = *up.Email
	}
	if up.AllowAnalytics != nil {
		user.AllowAnalytics = *up.AllowAnalytics
	}
	return user
}

func (u *User) IsAdminUser() bool {
	return (u.IsActive && u.IsAdmin)
}

// Anon returns a new anonymous user.
func UserAnon() *User {
	return &User{
		ID:     "",
		isAnon: true,
	}
}

// IsAnon returns true if the user is anonymous.
func (u *User) IsAnon() bool {
	return u.ID == "" && u.isAnon
}

func CreateUser(ctx context.Context, user *User) *gorm.DB {
	return db.Create(user)
}
