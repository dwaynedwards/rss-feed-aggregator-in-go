package rssfeed

import "time"

type User struct {
	ID int64 `json:"id"`

	Name string `json:"name"`

	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
}
