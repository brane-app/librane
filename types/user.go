package types

import (
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"

	"encoding/json"
	"time"
)

var (
	userPrivateFields []string = []string{
		"email",
		"created",
	}
)

type User struct {
	ID                string `json:"id" db:"id"`
	Email             string `json:"email" db:"email"`
	Nick              string `json:"nick" db:"nick"`
	Bio               string `json:"bio" db:"bio"`
	SubscriberCount   int    `json:"subscriber_count" db:"subscriber_count"`
	SubscriptionCount int    `json:"subscription_count" db:"subscription_count"`
	PostCount         int    `json:"post_count" db:"post_count"`
	Created           int64  `json:"created" db:"created"`
	Moderator         bool   `json:"moderator" db:"moderator"`
	Admin             bool   `json:"admin" db:"admin"`
}

func (it *User) FromMap(data map[string]interface{}) (err error) {
	var config mapstructure.DecoderConfig = mapstructure.DecoderConfig{
		Metadata: nil,
		TagName:  "json",
		Result:   &it,
	}

	var decoder *mapstructure.Decoder
	if decoder, err = mapstructure.NewDecoder(&config); err == nil {
		err = decoder.Decode(data)
	}

	return
}

func (user User) Map() (data map[string]interface{}) {
	var bytes []byte
	bytes, _ = json.Marshal(user)
	json.Unmarshal(bytes, &data)
	return
}

func (user User) PublicMap() (data map[string]interface{}) {
	data = user.Map()

	var field string
	for _, field = range userPrivateFields {
		delete(data, field)
	}

	return
}

func (user User) JSON() (data []byte, err error) {
	data, err = json.Marshal(user)
	return
}

func NewUser(nick, bio, email string) (made User) {
	made = User{
		Nick:  nick,
		Email: email,
		Bio:   bio,

		ID:      uuid.New().String(),
		Created: time.Now().Unix(),
	}

	return
}

func UserFromMap(data map[string]interface{}) (it User, err error) {
	var config mapstructure.DecoderConfig = mapstructure.DecoderConfig{
		Metadata: nil,
		TagName:  "json",
		Result:   &it,
	}

	var decoder *mapstructure.Decoder
	if decoder, err = mapstructure.NewDecoder(&config); err == nil {
		err = decoder.Decode(data)
	}

	return
}
