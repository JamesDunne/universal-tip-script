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
	"github.com/karrick/tparse"
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

type JsonObject map[string]interface{}

func main() {
	// arg is our input text:
	arg := ""
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}

	// start to build the output JSON array:
	jsonItems := make([]JsonObject, 0, 20)
	jsonItems = append(jsonItems, JsonObject{"type": "text", "value": arg})

	var t_utcstr string
	var t_eststr string = ""
	var t_cststr string = ""
	ts_str_sec := ""
	ts_str_msec := ""
	ts_str_nsec := ""
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
	}
	if terr != nil {
		var ts int64
		ts, terr = strconv.ParseInt(arg, 10, 64)
		if terr == nil {
			if ts < 99_999_999_999 {
				// seconds
				t = time.Unix(ts, 0)
			} else if ts < 99_999_999_999_999 {
				// milliseconds
				t = time.Unix(ts/1_000, (ts%1_000)*1_000_000)
			} else if ts < 99_999_999_999_999_999 {
				// microseconds
				t = time.Unix(ts/1_000_000, (ts%1_000_000)*1_000)
			} else {
				// nanoseconds
				t = time.Unix(ts/1_000_000_000, (ts % 1_000_000_000))
			}
		}
	}
	if terr != nil {
		// parse as a relative timestamp:
		// arg e.g. "now+1d-3w4mo+7y6h4m"
		var tr time.Time
		tr, terr = tparse.ParseNow(time.RFC3339, arg)
		if terr == nil {
			t = tr
		}
	}
	if terr != nil {
		t_utcstr = "failed parsing time"
	}

	if terr == nil {
		t_utcstr = t.UTC().Format(time.RFC3339Nano)
		t_eststr = t.In(est).Format(time.RFC3339Nano)
		t_cststr = t.In(cst).Format(time.RFC3339Nano)
		// format as unix epoch values in seconds, milliseconds, nanoseconds:
		ts_str_sec = strconv.FormatInt(t.UTC().Unix(), 10)
		ts_str_msec = strconv.FormatInt(t.UTC().UnixNano()/1_000_000, 10)
		ts_str_nsec = strconv.FormatInt(t.UTC().UnixNano(), 10)
		jsonItems = append(jsonItems,
			JsonObject{"type": "text", "value": "--- time RFC3339 (UTC, EST, CST):"},
			JsonObject{"type": "text", "value": t_utcstr},
			JsonObject{"type": "text", "value": t_eststr},
			JsonObject{"type": "text", "value": t_cststr},
			JsonObject{"type": "text", "value": "--- time Unix (sec, msec, nsec):"},
			JsonObject{"type": "text", "value": ts_str_sec},
			JsonObject{"type": "text", "value": ts_str_msec},
			JsonObject{"type": "text", "value": ts_str_nsec},
		)
	}

	encoding := base64.StdEncoding
	b64enc := encoding.EncodeToString([]byte(arg))
	b64dec, err := encoding.DecodeString(arg)
	if err != nil {
		b64dec = []byte(fmt.Sprintf("!error: %v", err))
	}
	jsonItems = append(jsonItems,
		JsonObject{"type": "text", "value": "--- base64 (enc, dec):"},
		JsonObject{"type": "text", "value": b64enc},
		JsonObject{"type": "text", "value": string(b64dec)},
	)

	u, _ := uuid.NewRandom()
	jsonItems = append(jsonItems,
		JsonObject{"type": "text", "value": "--- generated uuid:"},
		JsonObject{"type": "text", "value": u.String()},
	)

	o := json.NewEncoder(os.Stdout)
	err = o.Encode(jsonItems)
	if err != nil {
		log.Fatal(err)
	}
}
