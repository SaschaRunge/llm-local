package communication

type Role string

const (
	RoleSystem    Role = "system"
	RoleAssistant Role = "assistant"
	RoleUser      Role = "user"
)

type Message struct {
	AuthorName string
	Role       Role
	Reasoning  string
	Content    string
}

func (r Role) IsValid() bool {
	switch r {
	case RoleSystem, RoleAssistant, RoleUser:
		return true
	default:
		return false
	}
}
