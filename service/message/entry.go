package message

type Entry struct {
	ID                string `json:"id,omitempty"`
	UserID            string `json:"user_id,omitempty"`
	WorkflowID        string `json:"workflow_id,omitempty"`
	Status            int    `json:"status,omitempty"`
	Digest            string `json:"digest,omitempty"`
	CreatedAt         string `json:"created_at,omitempty"`
	Channel           string `json:"channel,omitempty"`
	ChannelMessageID  string `json:"channel_message_id,omitempty"`
	ChannelIdentifier string `json:"channel_identifier,omitempty"`
}

func (e *Entry) TableName() string {
	return "message"
}
