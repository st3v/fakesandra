package v3

import (
	"encoding/binary"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"
)

type consistency uint16

const (
	cAny consistency = iota
	cOne
	cTwo
	cThree
	cQuorum
	cAll
	cLocalQuorum
	cEachQuorum
	cSerial
	cLocalSerial
	cLocalOne
)

func (c consistency) String() string {
	switch c {
	case cAny:
		return "ANY"
	case cOne:
		return "ONE"
	case cTwo:
		return "TWO"
	case cThree:
		return "THREE"
	case cQuorum:
		return "QUORUM"
	case cAll:
		return "ALL"
	case cLocalQuorum:
		return "LOCAL_QUORUM"
	case cEachQuorum:
		return "EACH_QUORUM"
	case cSerial:
		return "SERIAL"
	case cLocalSerial:
		return "LOCAL_SERIAL"
	case cLocalOne:
		return "LOCAL_SERIAL"
	default:
		return "UNKNOWN"
	}
}

type queryFlagSet uint8
type queryFlag uint8

const (
	qryValues queryFlag = 1 << iota
	qrySkipMeta
	qryPageSize
	qryPagingState
	qrySerialConsistency
	qryDefaultTimestamp
	qryNames
)

func (fs queryFlagSet) Contains(m queryFlag) bool {
	return byte(fs)&byte(m) == byte(m)
}

func (fs queryFlagSet) Flags() map[string]queryFlag {
	flags := map[string]queryFlag{}
	for f, n := range queryFlagNames {
		if fs.Contains(f) {
			flags[n] = f
		}
	}
	return flags
}

func (fs queryFlagSet) String() string {
	names := []string{}

	for k, _ := range fs.Flags() {
		names = append(names, k)
	}

	return strings.Join(names, ", ")
}

func (f queryFlag) String() string {
	name, found := queryFlagNames[f]
	if !found {
		return "UNKNOWN"
	}
	return name
}

var queryFlagNames = map[queryFlag]string{
	qryValues:            "VALUES",
	qrySkipMeta:          "SKIP_METADATA",
	qryPageSize:          "PAGE_SIZE",
	qryPagingState:       "WITH_PAGING_STATE",
	qrySerialConsistency: "WITH_SERIAL_CONSISTENCY",
	qryDefaultTimestamp:  "WITH_DEFAULT_TIMESTAMP",
	qryNames:             "WITH_NAMES",
}

type query struct {
	stmt              string
	consistency       consistency
	flagSet           queryFlagSet
	values            [][]byte
	valueNames        []string
	pageSize          int32
	pagingState       []byte
	serialConsistency consistency
	defaultTimestamp  time.Time
}

func (q query) Values() ([][]byte, bool) {
	return q.values, q.flagSet.Contains(qryValues)
}

func (q query) NamedValues() (map[string][]byte, bool) {
	nv := map[string][]byte{}

	for i, name := range q.valueNames {
		nv[name] = q.values[i]
	}

	return nv, q.flagSet.Contains(qryNames) && q.flagSet.Contains(qryValues)
}

func (q query) PageSize() (int32, bool) {
	return q.pageSize, q.flagSet.Contains(qryPageSize)
}

func (q query) PagingState() ([]byte, bool) {
	return q.pagingState, q.flagSet.Contains(qryPagingState)
}

func (q query) SerialConsistency() (consistency, bool) {
	return q.serialConsistency, q.flagSet.Contains(qrySerialConsistency)
}

func (q query) DefaultTimestamp() (time.Time, bool) {
	return q.defaultTimestamp, q.flagSet.Contains(qryDefaultTimestamp)
}

func (q query) String() string {
	format := `Query [ Statement: "%s", Consistency: "%s", Flags: "%s"%s ]`

	optional := []string{}

	if ps, set := q.PageSize(); set {
		optional = append(optional, fmt.Sprintf(`PageSize: "%d"`, ps))
	}

	if ps, set := q.PagingState(); set {
		optional = append(optional, fmt.Sprintf(`PageStateLen: "%d"`, len(ps)))
	}

	if sc, set := q.SerialConsistency(); set {
		optional = append(optional, fmt.Sprintf(`SerialConsistency: "%s"`, sc))
	}

	if ts, set := q.DefaultTimestamp(); set {
		optional = append(optional, fmt.Sprintf(`DefaultTimestamp: "%s"`, ts))
	}

	newlines := regexp.MustCompile(`[\r\n]`)
	stmt := newlines.ReplaceAllString(q.stmt, " ")

	spaces := regexp.MustCompile(`[\s\t]+`)
	stmt = spaces.ReplaceAllString(stmt, " ")

	stmt = strings.Trim(stmt, " ")

	options := strings.Join(optional, ", ")
	if options != "" {
		options = ", " + options
	}

	return fmt.Sprintf(format, stmt, q.consistency, q.flagSet, options)
}

func readByte(r io.Reader, n *uint8) error {
	return binary.Read(r, binary.BigEndian, n)
}

func readShort(r io.Reader, n *uint16) error {
	return binary.Read(r, binary.BigEndian, n)
}

func readInt(r io.Reader, n *int32) error {
	return binary.Read(r, binary.BigEndian, n)
}

func readLong(r io.Reader, n *int64) error {
	return binary.Read(r, binary.BigEndian, n)
}

func readBytes(r io.Reader) ([]byte, error) {
	var n int32
	if err := readInt(r, &n); err != nil {
		return []byte{}, err
	}

	b := make([]byte, n)
	if _, err := io.ReadFull(r, b); err != nil {
		return []byte{}, err
	}

	return b, nil
}

func readShortBytes(r io.Reader) ([]byte, error) {
	var n uint16
	if err := readShort(r, &n); err != nil {
		return []byte{}, err
	}

	b := make([]byte, n)
	if _, err := io.ReadFull(r, b); err != nil {
		return []byte{}, err
	}

	return b, nil
}

func readString(r io.Reader) (string, error) {
	str, err := readShortBytes(r)
	return string(str), err
}

func readLongString(r io.Reader) (string, error) {
	str, err := readBytes(r)
	return string(str), err
}

func readConsistency(r io.Reader, c *consistency) error {
	return binary.Read(r, binary.BigEndian, c)
}

func readQuery(r io.Reader, q *query) error {
	var err error
	if q.stmt, err = readLongString(r); err != nil {
		return err
	}

	if err := readConsistency(r, &q.consistency); err != nil {
		return err
	}

	if err := binary.Read(r, binary.BigEndian, &q.flagSet); err != nil {
		return err
	}

	if q.values, q.valueNames, err = readValues(r, q.flagSet); err != nil {
		return err
	}

	if err := readPageSize(r, q.flagSet, &q.pageSize); err != nil {
		return err
	}

	if q.pagingState, err = readPagingState(r, q.flagSet); err != nil {
		return err
	}

	if err := readSerialConsistency(r, q.flagSet, &q.serialConsistency); err != nil {
		return err
	}

	if q.defaultTimestamp, err = readDefaultTimestamp(r, q.flagSet); err != nil {
		return err
	}

	return nil
}

func readValues(r io.Reader, fs queryFlagSet) ([][]byte, []string, error) {
	errResult := func(err error) ([][]byte, []string, error) {
		return [][]byte{}, []string{}, err
	}

	if !fs.Contains(qryValues) {
		return errResult(nil)
	}

	var numValues uint16
	if err := readShort(r, &numValues); err != nil {
		return errResult(nil)
	}

	var err error
	names := make([]string, numValues)
	values := make([][]byte, numValues)

	for i := uint16(0); i < numValues; i++ {
		if fs.Contains(qryNames) {
			if names[i], err = readString(r); err != nil {
				return errResult(err)
			}
		}

		if values[i], err = readBytes(r); err != nil {
			return errResult(err)
		}
	}

	return values, names, nil
}

func readPageSize(r io.Reader, fs queryFlagSet, ps *int32) error {
	if !fs.Contains(qryPageSize) {
		return nil
	}
	return readInt(r, ps)
}

func readPagingState(r io.Reader, fs queryFlagSet) ([]byte, error) {
	if !fs.Contains(qryPagingState) {
		return []byte{}, nil
	}
	return readBytes(r)
}

func readSerialConsistency(r io.Reader, fs queryFlagSet, c *consistency) error {
	if !fs.Contains(qrySerialConsistency) {
		return nil
	}
	return readConsistency(r, c)
}

func readDefaultTimestamp(r io.Reader, fs queryFlagSet) (time.Time, error) {
	ts := time.Unix(0, 0)

	if !fs.Contains(qryDefaultTimestamp) {
		return ts, nil
	}

	var ms int64
	if err := readLong(r, &ms); err != nil {
		return ts, err
	}

	return ts.Add(time.Duration(ms) * time.Microsecond), nil
}
