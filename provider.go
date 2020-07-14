package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func tryParseTime(value string, layouts ...string) (t time.Time, err error) {
	for _, layout := range layouts {
		t, err = time.Parse(layout, value)
		if err == nil {
			return
		}
	}

	return
}

func tryParseTimeInLocation(value string, location *time.Location, layouts ...string) (t time.Time, err error) {
	for _, layout := range layouts {
		t, err = time.ParseInLocation(layout, value, location)
		if err == nil {
			return
		}
	}

	return
}

func main() {
	// arg is our input text:
	arg := ""
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}

	o := json.NewEncoder(os.Stdout)

	encoding := base64.StdEncoding
	b64enc := encoding.EncodeToString([]byte(arg))
	b64dec, err := encoding.DecodeString(arg)
	if err != nil {
		b64dec = []byte(fmt.Sprintf("(%v)", err))
	}

	u, _ := uuid.NewRandom()

	var t_utcstr string
	var t_eststr string = ""
	var t_cststr string = ""
	var t_ms string = ""
	est, _ := time.LoadLocation("America/New_York")
	cst, _ := time.LoadLocation("America/Chicago")
	// try parsing with time.RFC3339 with nanos, without nanos, with 'T', and with ' ':
	layouts := []string{
		"2006-01-02T15:04:05.999999999Z07:00",
		"2006-01-02T15:04:05.999Z07:00",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05.999999999Z07:00",
		"2006-01-02 15:04:05.999Z07:00",
		"2006-01-02 15:04:05Z07:00",
	}
	t, terr := tryParseTime(arg, layouts...)
	if terr != nil {
		layouts := []string{
			"2006-01-02T15:04:05.999999999",
			"2006-01-02T15:04:05.999",
			"2006-01-02T15:04:05",
			"2006-01-02 15:04:05.999999999",
			"2006-01-02 15:04:05.999",
			"2006-01-02 15:04:05",
		}
		t, terr = tryParseTimeInLocation(arg, cst, layouts...)
		if terr != nil {
			t_utcstr = "failed parsing time"
		}
	}
	if terr == nil {
		t_utcstr = t.UTC().Format(time.RFC3339Nano)
		t_eststr = t.In(est).Format(time.RFC3339Nano)
		t_cststr = t.In(cst).Format(time.RFC3339Nano)
		t_ms = fmt.Sprintf("%v", t.UTC().Unix() * 1000)
	}

	var ts int64
	var ts_err error
	ts_str_sec := ""
	ts_str_msec := ""
	ts_str_nsec := ""
	ts, ts_err = strconv.ParseInt(arg, 10, 64)
	if ts_err != nil {
		ts_str_sec = "not an int64 timestamp"
	}
	if ts_err == nil {
		tm := time.Unix(ts, 0)
		ts_str_sec = tm.UTC().Format(time.RFC3339Nano)
		tm = time.Unix(ts / 1000, (ts % 1000) * 1_000_000)
		ts_str_msec = tm.UTC().Format(time.RFC3339Nano)
		tm = time.Unix(ts / 1_000_000_000, ts % 1_000_000_000)
		ts_str_nsec = tm.UTC().Format(time.RFC3339Nano)
	}

	err = o.Encode([]map[string]interface{}{
		{"type": "text", "value": "--- base64:"},
		{"type": "text", "value": b64enc},
		{"type": "text", "value": string(b64dec)},
		{"type": "text", "value": "--- uuid:"},
		{"type": "text", "value": u.String()},
		{"type": "text", "value": "--- time parse RFC3339 (UTC, EST, CST):"},
		{"type": "text", "value": t_utcstr},
		{"type": "text", "value": t_eststr},
		{"type": "text", "value": t_cststr},
		{"type": "text", "value": t_ms},
		{"type": "text", "value": "--- time parse int64 timestamp (sec, msec, nsec):"},
		{"type": "text", "value": ts_str_sec},
		{"type": "text", "value": ts_str_msec},
		{"type": "text", "value": ts_str_nsec},
		{"type": "text", "value": ts},
	})
	if err != nil {
		log.Fatal(err)
	}
}
