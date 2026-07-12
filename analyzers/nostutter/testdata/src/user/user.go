package user

type UserService struct{} // want "type UserService stutters with package name user; consider renaming to Service"

type Service struct{}

// Allowed: exact package name match (user.User is idiomatic)
type User struct{}

type UserID string // want "type UserID stutters with package name user; consider renaming to ID"
