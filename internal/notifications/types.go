package notifications

// NotifType identifies the kind of notification.
type NotifType string

const (
	NotifLowHunger     NotifType = "low_hunger"
	NotifLowHappiness  NotifType = "low_happiness"
	NotifHighBoredom   NotifType = "high_boredom"
	NotifLowHygiene    NotifType = "low_hygiene"
	NotifSickness      NotifType = "sickness"
	NotifDeath         NotifType = "death"
	NotifFriendRequest NotifType = "friend_request"
	NotifChatMessage   NotifType = "chat_message"
	NotifBedtime       NotifType = "bedtime"
	NotifWakeUp        NotifType = "wake_up"
	NotifEvent         NotifType = "event"
)

// PushMessage is the payload sent to a device.
type PushMessage struct {
	To    string `json:"to"`
	Title string `json:"title"`
	Body  string `json:"body"`
	Data  any    `json:"data,omitempty"`
	Sound string `json:"sound,omitempty"`
}

// PushToken is a stored device token.
type PushToken struct {
	TokenID  int    `db:"TokenId"`
	UserID   int    `db:"UserId"`
	Token    string `db:"Token"`
	Platform string `db:"Platform"`
}

// ThrottleEntry tracks when a notification type was last sent to a user.
type ThrottleEntry struct {
	LogID       int    `db:"LogId"`
	UserID      int    `db:"UserId"`
	NotifType   string `db:"NotifType"`
	SentAt      string `db:"SentAt"`
	EscalationN int    `db:"EscalationN"`
}
