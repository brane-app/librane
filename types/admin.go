package types

import (
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"

	"encoding/json"
	"time"
)

type Ban struct {
	ID      string `json:"id" db:"id"`
	Banner  string `json:"banner" db:"banner"`
	Banned  string `json:"banned" db:"banned"`
	Reason  string `json:"reason" db:"reason"`
	Created int64  `json:"created" db:"created"`
	Expires int64  `json:"expires" db:"expires"`
	Forever bool   `json:"forever" db:"forever"`
}

func (ban Ban) Map() (data map[string]interface{}) {
	data = map[string]interface{}{
		"id":      ban.ID,
		"banned":  ban.Banned,
		"banner":  ban.Banner,
		"reason":  ban.Reason,
		"expires": ban.Expires,
		"created": ban.Created,
		"forever": ban.Forever,
	}

	return
}

func (ban Ban) JSON() (data []byte, err error) {
	data, err = json.Marshal(ban)
	return
}

func (it *Ban) FromMap(data map[string]interface{}) (err error) {
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

func NewBan(banner, banned, reason string, duration int64, forever bool) (ban Ban) {
	var now int64 = time.Now().Unix()

	ban = Ban{
		Banner:  banner,
		Banned:  banned,
		Reason:  reason,
		Forever: forever,

		ID:      uuid.New().String(),
		Created: now,
		Expires: now + duration,
	}

	return
}

type Report struct {
	ID         string `json:"id" db:"id"`
	Reporter   string `json:"reporter" db:"reporter"`
	Reported   string `json:"reported" db:"reported"`
	Type       string `json:"type" db:"type"`
	Reason     string `json:"reason" db:"reason"`
	Created    int64  `json:"created" db:"created"`
	Resolved   bool   `json:"resolved" db:"resolved"`
	Resolution string `json:"resolution" db:"resolution"`
}

func (report Report) Map() (data map[string]interface{}) {
	data = map[string]interface{}{
		"id":         report.ID,
		"reporter":   report.Reporter,
		"reported":   report.Reported,
		"type":       report.Type,
		"reason":     report.Reason,
		"created":    report.Created,
		"resolved":   report.Resolved,
		"resolution": report.Resolution,
	}

	return
}

func (report Report) JSON() (data []byte, err error) {
	data, err = json.Marshal(report)
	return
}

func (it *Report) FromMap(data map[string]interface{}) (err error) {
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

func NewReport(reporter, reported, report_type, reason string) (report Report) {
	report = Report{
		Reporter: reporter,
		Reported: reported,
		Type:     report_type,
		Reason:   reason,

		ID:      uuid.New().String(),
		Created: time.Now().Unix(),
	}

	return
}
