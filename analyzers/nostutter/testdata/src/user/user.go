package user

type UserService struct{} // want "type UserService stutters with package name user; consider renaming to Service"

type Service struct{}

// Allowed: exact package name match (user.User is idiomatic)
type User struct{}

// Allowed: suffix starts with lowercase (natural derivation like io.Reader)
type Userinfo struct{}

type UserID string // want "type UserID stutters with package name user; consider renaming to ID"
