package v3

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/st3v/fakesandra/cql/proto"
)

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

type Query struct {
	Statement         string
	Consistency       proto.Consistency
	flagSet           queryFlagSet
	values            [][]byte
	valueNames        []string
	pageSize          int32
	pagingState       []byte
	serialConsistency proto.Consistency
	defaultTimestamp  time.Time
}

func (q Query) Values() ([][]byte, bool) {
	return q.values, q.flagSet.Contains(qryValues)
}

func (q Query) NamedValues() (map[string][]byte, bool) {
	nv := map[string][]byte{}

	for i, name := range q.valueNames {
		nv[name] = q.values[i]
	}

	return nv, q.flagSet.Contains(qryNames) && q.flagSet.Contains(qryValues)
}

func (q Query) SkipMetadata() bool {
	return q.flagSet.Contains(qrySkipMeta)
}

func (q Query) PageSize() (int32, bool) {
	return q.pageSize, q.flagSet.Contains(qryPageSize)
}

func (q Query) PagingState() ([]byte, bool) {
	return q.pagingState, q.flagSet.Contains(qryPagingState)
}

func (q Query) SerialConsistency() (proto.Consistency, bool) {
	return q.serialConsistency, q.flagSet.Contains(qrySerialConsistency)
}

func (q Query) DefaultTimestamp() (time.Time, bool) {
	return q.defaultTimestamp, q.flagSet.Contains(qryDefaultTimestamp)
}

func (q Query) TrimmedStatement() string {
	newlines := regexp.MustCompile(`[\r\n]`)
	stmt := newlines.ReplaceAllString(q.Statement, " ")

	spaces := regexp.MustCompile(`[\s\t]+`)
	stmt = spaces.ReplaceAllString(stmt, " ")

	return strings.Trim(stmt, " ")
}

func (q Query) String() string {
	fields := []string{
		fmt.Sprintf(`Statement: "%s"`, q.TrimmedStatement()),
		fmt.Sprintf(`Consistency: "%s"`, q.Consistency),
		fmt.Sprintf(`Flags: "%s"`, q.flagSet),
		fmt.Sprintf(`SkipMetadata: %t`, q.SkipMetadata()),
	}

	if ps, set := q.PageSize(); set {
		fields = append(fields, fmt.Sprintf(`PageSize: %d`, ps))
	}

	if ps, set := q.PagingState(); set {
		fields = append(fields, fmt.Sprintf(`PagingStateLength: %d`, len(ps)))
	}

	if sc, set := q.SerialConsistency(); set {
		fields = append(fields, fmt.Sprintf(`SerialConsistency: "%s"`, sc))
	}

	if ts, set := q.DefaultTimestamp(); set {
		fields = append(fields, fmt.Sprintf(`DefaultTimestamp: "%s"`, ts))
	}

	return fmt.Sprintf("Query [ %s ]", strings.Join(fields, ", "))
}

func readQuery(r io.Reader, q *Query) error {
	var err error
	if q.Statement, err = proto.ReadLongString(r); err != nil {
		return err
	}

	if err := proto.ReadConsistency(r, &q.Consistency); err != nil {
		return err
	}

	if err := proto.ReadBinary(r, &q.flagSet); err != nil {
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
	if err := proto.ReadShort(r, &numValues); err != nil {
		return errResult(nil)
	}

	var err error
	names := make([]string, numValues)
	values := make([][]byte, numValues)

	for i := uint16(0); i < numValues; i++ {
		if fs.Contains(qryNames) {
			if names[i], err = proto.ReadString(r); err != nil {
				return errResult(err)
			}
		}

		if values[i], err = proto.ReadBytes(r); err != nil {
			return errResult(err)
		}
	}

	return values, names, nil
}

func readPageSize(r io.Reader, fs queryFlagSet, ps *int32) error {
	if !fs.Contains(qryPageSize) {
		return nil
	}
	return proto.ReadInt(r, ps)
}

func readPagingState(r io.Reader, fs queryFlagSet) ([]byte, error) {
	if !fs.Contains(qryPagingState) {
		return []byte{}, nil
	}
	return proto.ReadBytes(r)
}

func readSerialConsistency(r io.Reader, fs queryFlagSet, c *proto.Consistency) error {
	if !fs.Contains(qrySerialConsistency) {
		return nil
	}
	return proto.ReadConsistency(r, c)
}

func readDefaultTimestamp(r io.Reader, fs queryFlagSet) (time.Time, error) {
	ts := time.Unix(0, 0)

	if !fs.Contains(qryDefaultTimestamp) {
		return ts, nil
	}

	var ms int64
	if err := proto.ReadLong(r, &ms); err != nil {
		return ts, err
	}

	return ts.Add(time.Duration(ms) * time.Microsecond), nil
}
