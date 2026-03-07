package permission

// Default always returns Member.
func Default() Permission {
	return Member{}
}

// Member is the default permission level of a server, 0 is the permission level.
type Member struct{}

// Level ...
func (p Member) Level() int32 {
	return 0
}

// Name ...
func (p Member) Name() string {
	return "member"
}

// Operator represent a higher level of permissions, 1 is the permission level.
type Operator struct{}

// Name ...
func (p Operator) Name() string {
	return "operator"
}

// Level ...
func (p Operator) Level() int32 {
	return 1
}

// Visitor is a special type of permission that revokes every natural minecraft function to the player, for example breaking, or placing
// blocks, 2 is the permission level.
type Visitor struct{}

// Name ...
func (p Visitor) Name() string {
	return "visitor"
}

// Level ...
func (p Visitor) Level() int32 {
	return 2
}

// Custom is a customizable type of permission level, here you can customize if the player can place of break blocks for example, the
// permission level is 3.
type Custom struct{}

// Name ...
func (p Custom) Name() string {
	return "custom"
}

// Level ...
func (p Custom) Level() int32 {
	return 3
}
