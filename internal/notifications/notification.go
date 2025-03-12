package notifications

type Notification struct {
	Title   string
	Message string
	Level   string
}

func (n *Notification) IsSame(other Notification) bool {
	return n.Title == other.Title && n.Message == other.Message && n.Level == other.Level
}
