package params

type Message struct {
	Id        int    `json:"id" db:"message_id"`
	Hash      string `json:"hash" db:"message_hash"`
	Content   string `json:"content"`
	Name      string `json:"name" db:"user_name"`
	Timestamp uint64 `json:"timestamp"`
	IsStored  bool   `json:"is_stored" db:"is_stored"`
}
