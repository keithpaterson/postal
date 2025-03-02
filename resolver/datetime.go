package resolver

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	defaultDateTimeFormat = time.RFC3339
	defaultDateFormat     = time.DateOnly
	defaultTimeFormat     = time.TimeOnly
)

const (
	datespecNow = "now"
)

var (
	// parses: "now.(format) + 2D" into ["now", "format", "+" and "2D"]
	regTokenValue = regexp.MustCompile(`(?U)^(.+)(?:\.\((.+)\))?(?:\s+([\+\-])\s+(\S*))?$`)
	// parses: "2D" into ["2", "D"]
	regField = regexp.MustCompile(`^(\d+)([YMWDhms])$`)
)

// supported operation types

const (
	opAdd operation = iota
	opSubtract
)

var (
	opType = map[string]operation{"+": opAdd, "-": opSubtract}
)

type operation int

// supported modifier fields

const (
	fieldYear field = iota
	fieldMonth
	fieldWeek
	fieldDay
	fieldHour
	fieldMinute
	fieldSecond
)

var (
	fieldType = map[string]field{
		"Y": fieldYear, "M": fieldMonth, "W": fieldWeek, "D": fieldDay,
		"h": fieldHour, "m": fieldMinute, "s": fieldSecond}
)

type field int

type timeValueModifier struct {
	op    operation
	field field
	delta int
}

func (m timeValueModifier) adjust(from time.Time) time.Time {
	if m.delta == 0 {
		return from
	}

	delta := m.delta
	if m.op == opSubtract {
		delta *= -1
	}

	switch m.field {
	case fieldYear:
		return from.AddDate(delta, 0, 0)
	case fieldMonth:
		return from.AddDate(0, delta, 0)
	case fieldWeek:
		return from.AddDate(0, 0, 7*delta)
	case fieldDay:
		return from.AddDate(0, 0, delta)
	case fieldHour:
		return from.Add(time.Hour * time.Duration(delta))
	case fieldMinute:
		return from.Add(time.Minute * time.Duration(delta))
	case fieldSecond:
		return from.Add(time.Second * time.Duration(delta))
	}

	// this is an error condition (how did we get here?)
	return from
}

type timeValue struct {
	time   time.Time
	format string
}

type dateTimeResolver struct {
	resolverImpl

	nowFn func() time.Time
}

func newDateTimeResolver(root *rootResolver) *dateTimeResolver {
	return &dateTimeResolver{
		resolverImpl{root: root, log: root.log.Named("datetime")},
		time.Now,
	}
}

func (r *dateTimeResolver) resolve(name string, value string) (string, bool) {
	var err error
	var result string
	changed := true
	switch name {
	case "datetime":
		result, err = r.resolveDateTimeValue(value)
	case "date":
		result, err = r.resolveDateValue(value)
	case "time":
		result, err = r.resolveTimeValue(value)
	case "epoch":
		result, err = r.resolveEpochValue(value)
	default:
		err = fmt.Errorf("invalid date/time token name '%s'", name)
	}
	if err != nil {
		// log a warning?
		changed = false
		result = value
	}
	return result, changed
}

func (r *dateTimeResolver) now() time.Time {
	return r.nowFn()
}

func (r *dateTimeResolver) parseTokenValue(value string) (spec string, format string, op string, delta string) {
	// different possible strings to process, generally combinations of the string:
	// 'now.(YYMMDD) + 2D'
	//    spec='now', format='YYMMDD', op='+', delta='2D'
	// 'now' is required, everything else is optional, and the regex captures are nicely consistent
	matches := regTokenValue.FindStringSubmatch(value)
	spec = strings.TrimSpace(matches[1])
	format = strings.TrimSpace(matches[2])
	op = strings.TrimSpace(matches[3])
	delta = strings.TrimSpace(matches[4])
	return
}

func (r *dateTimeResolver) resolveDateTimeValue(value string) (string, error) {
	return r.resolveDateAndTime(value, defaultDateTimeFormat)
}

func (r *dateTimeResolver) resolveDateValue(value string) (string, error) {
	return r.resolveDateAndTime(value, defaultDateFormat)
}

func (r *dateTimeResolver) resolveTimeValue(value string) (string, error) {
	return r.resolveDateAndTime(value, defaultTimeFormat)
}

func (r *dateTimeResolver) resolveDateAndTime(value string, defaultFormat string) (string, error) {
	// split into components:
	// "datetime-spec" "op" "modifier"
	// e.g. the token "${datetime:now.(UnixDate) + 1M}" becomes "UnixDate", "+" "1M"
	//      the token "${date:now.(DateOnly) + 1M}" becomes "DateOnly", "+" "1M"
	//      the token "${time:now.(TimeOnly) + 1h}" becomes "TimeOnly", "+" "1M"
	//      the token "${epoch:now + 30s}" becomes "epoch" "+" "30s"
	// first we need to process any embedded tokens in the value itself.
	var err error
	value = r.resolveValue(value)

	spec, format, op, delta := r.parseTokenValue(value)
	if format == "" {
		format = defaultFormat
	}

	var modifier timeValueModifier
	if modifier, err = r.toModifier(op, delta); err != nil {
		return value, err
	}

	var dateValue timeValue
	if dateValue, err = r.toDateTime(spec, format, modifier); err != nil {
		return value, err
	}

	return r.format(dateValue), nil
}

func (r *dateTimeResolver) resolveEpochValue(value string) (result string, err error) {
	// split into components:
	// "epoch-spec" "op" "modifier"
	// e.g. the token "${epoch:now + 30s}" becomes "now" "+" "30s"
	spec, _, op, delta := r.parseTokenValue(value)

	var modifier timeValueModifier
	if modifier, err = r.toModifier(op, delta); err != nil {
		return "", err
	}

	var epochValue int64
	if epochValue, err = r.toEpoch(strings.ToLower(spec), modifier); err != nil {
		return
	}

	result = strconv.FormatInt(epochValue, 10)
	return
}

func (r *dateTimeResolver) format(value timeValue) string {
	return value.time.Format(r.toLayout(value.format))
}

func (r *dateTimeResolver) toDateTime(timeval string, format string, modifier timeValueModifier) (result timeValue, err error) {
	// date spec can be several things:
	// - "now" (case insensitive)
	// - a collection of numbers and characters that can be converted to a date+time format
	//   (either supported by a Layout or some combination oof YYMMDDhhmmss
	//   So in the non-now case the format will be required (or assumed)
	//
	// optionally with .(format), where suffix can be:
	//   - any valid time package const
	//   - "ISO8601" => "RFC3339"
	result = timeValue{}

	var actual time.Time
	if actual, err = r.parseDateTime(timeval, format); err != nil {
		r.log.Error("error", fmt.Sprintf("failed to parse '%s': %s\n", timeval, err.Error()))
		return
	}

	result.time = modifier.adjust(actual)
	result.format = format
	return result, nil
}

func (r *dateTimeResolver) toEpoch(timeval string, modifier timeValueModifier) (result int64, err error) {
	var epoch time.Time
	if timeval == "now" {
		epoch = r.now()
	} else {
		var seconds int64
		if seconds, err = strconv.ParseInt(timeval, 10, 32); err != nil {
			return
		}
		epoch = time.Unix(seconds, 0)
	}

	epoch = modifier.adjust(epoch)
	result = epoch.Unix()
	return
}

func (r *dateTimeResolver) toModifier(opval string, modval string) (timeValueModifier, error) {
	modifier := timeValueModifier{}
	var err error

	if opval == "" && modval == "" {
		// valid: delta-zero modification
		return modifier, nil
	}
	if opval == "" || modval == "" {
		// error: both must be present
		return modifier, fmt.Errorf("invalid date/time operation: ('%s', '%s')", opval, modval)
	}

	var ok bool
	var op operation
	if op, ok = opType[opval]; !ok {
		return modifier, fmt.Errorf("invalid operation '%s'", opval)
	}

	modval = strings.TrimSpace(modval)
	matches := regField.FindStringSubmatch(modval)
	if len(matches) != 3 {
		// didn't match enough information; this is an invalid specification
		return modifier, fmt.Errorf("invalid date/time modifier '%s'", modval)
	}

	var delta int
	if delta, err = strconv.Atoi(matches[1]); err != nil {
		return modifier, fmt.Errorf("invalid date/time modifier delta '%s'", matches[1])
	}

	var field field
	if field, ok = fieldType[matches[2]]; !ok {
		return modifier, fmt.Errorf("invlaid date/time modifier field'%s'", matches[2])
	}

	modifier.op = op
	modifier.field = field
	modifier.delta = delta

	return modifier, nil
}

func (r *dateTimeResolver) parseDateTime(datespec string, format string) (result time.Time, err error) {
	if strings.ToLower(datespec) == datespecNow {
		return r.now(), nil
	}

	// otherwise it's some kind of string
	return time.Parse(r.toLayout(format), datespec)
}

var (
	dateFormatReplacements = []string{
		// time constants
		"ANSIC", time.ANSIC,
		"UnixDate", time.UnixDate,
		"RubyDate", time.RubyDate,
		"RFC822", time.RFC822,
		"RFC822Z", time.RFC822Z,
		"RFC850", time.RFC850,
		"RFC1123", time.RFC1123,
		"RFC1123Z", time.RFC1123Z,
		"ISO8601", time.RFC3339,
		"RFC3339", time.RFC3339,
		"RFC3339Nano", time.RFC3339Nano,
		"Kitchen", time.Kitchen,
		"Stamp", time.Stamp,
		"StampMilli", time.StampMilli,
		"StampMicro", time.StampMicro,
		"StampNano", time.StampNano,
		"DateTime", time.DateTime,
		"DateOnly", time.DateOnly,
		"TimeOnly", time.TimeOnly,
		// custom constants
		"Year", "2006",
		"YYYY", "2006",
		"YY", "06",
		// month
		"Month", "January",
		"Mon", "Jan",
		"MM", "01",
		"_M", "_1",
		// day
		"Weekday", "Monday",
		"Day", "Mon",
		"DD", "02",
		"_D", "_1",
		// hour
		"HH", "15", // 24h clock
		"hh", "03", // 12h clock
		"_h", "_3", // 12h clock
		// minute
		"mm", "04",
		// second
		"ss", "05",
		"AM", "PM",
		// timezone
		"TimeZone", "MST",
		"TZ", "-0700",
		"ZZ", "Z07:00",
	}
)

func (r *dateTimeResolver) toLayout(format string) string {
	if format == "" {
		// should log an error here, but this would be an internal screwup if it happened
		return defaultDateTimeFormat
	}

	// process the format string by looking for our custom format specifiers.
	// for each one, convert it to the time.Layout equivalent.
	// probably simpler to do a mass replace for our special strings
	replacer := strings.NewReplacer(dateFormatReplacements...)
	return replacer.Replace(format)
}
