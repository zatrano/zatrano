package features

// User is a minimal subject for role / rollout checks (optional for anonymous traffic).
type User struct {
	ID    uint
	Roles []string
}

// For returns an evaluator scoped to the given user (nil = anonymous).
func (r *Registry) For(u *User) *Eval {
	return &Eval{reg: r, user: u}
}

// ForUser is a convenience wrapper when you already have id and role slugs.
func (r *Registry) ForUser(id uint, roles []string) *Eval {
	if id == 0 && len(roles) == 0 {
		return &Eval{reg: r, user: nil}
	}
	return r.For(&User{ID: id, Roles: roles})
}
