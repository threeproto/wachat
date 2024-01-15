package params

type Message struct {
	Id            int    `json:"id" db:"message_id"`
	Hash          string `json:"hash" db:"message_hash"`
	Content       string `json:"content"`
	Name          string `json:"name" db:"user_name"`
	Timestamp     uint64 `json:"timestamp"`
	IsStored      bool   `json:"isStored" db:"is_stored"`
	WakuTimestamp uint64 `json:"wakuTimestamp" db:"waku_timestamp"`
}

type User struct {
	Id       int    `json:"id" db:"user_id"`
	Name     string `json:"name" db:"name"`
	Selected bool   `json:"selected" db:"selected"`
}
