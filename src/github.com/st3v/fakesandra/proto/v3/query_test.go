package v3

import (
	"bytes"
	"io"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("readQuery", func() {
	var (
		buf   *bytes.Buffer
		query Query
		err   error
	)

	JustBeforeEach(func() {
		err = readQuery(buf, &query)
	})

	Context("when there is a query to read", func() {
		var (
			stmt        = "SOME STATEMENT"
			consistency = LocalOne
		)

		BeforeEach(func() {
			// initialize empty buffer
			buf = bytes.NewBuffer([]byte{})

			// write statement
			err := writeLongString(buf, stmt)
			Expect(err).ToNot(HaveOccurred())

			// write consistency
			err = writeShort(buf, uint16(consistency))
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when there are no query flags", func() {
			BeforeEach(func() {
				// write empty flag set
				err = writeByte(buf, uint8(0))
				Expect(err).ToNot(HaveOccurred())
			})

			It("returns no error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("correctly parses the statement", func() {
				Expect(query.TrimmedStatement()).To(Equal(stmt))
			})

			It("correctly parses the consistency", func() {
				Expect(query.Consistency).To(Equal(consistency))
			})

			It("does not set any query values", func() {
				_, set := query.Values()
				Expect(set).To(BeFalse())
			})

			It("does not set any query value names", func() {
				_, set := query.NamedValues()
				Expect(set).To(BeFalse())
			})

			It("does not set the skip metadata option", func() {
				skip := query.SkipMetadata()
				Expect(skip).To(BeFalse())
			})

			It("does not set the page size option", func() {
				_, set := query.PageSize()
				Expect(set).To(BeFalse())
			})

			It("does not set the paging state option", func() {
				_, set := query.PagingState()
				Expect(set).To(BeFalse())
			})

			It("does not set the serial consistency option", func() {
				_, set := query.SerialConsistency()
				Expect(set).To(BeFalse())
			})

			It("does not set the default timestamp", func() {
				_, set := query.DefaultTimestamp()
				Expect(set).To(BeFalse())
			})
		})

		Context("when skip metadata is set", func() {
			BeforeEach(func() {
				err = writeByte(buf, uint8(qrySkipMeta))
				Expect(err).ToNot(HaveOccurred())
			})

			It("correctly parses the skip metadata option", func() {
				skip := query.SkipMetadata()
				Expect(skip).To(BeTrue())
			})
		})

		Context("when page size is set", func() {
			var pageSize = int32(12345)

			BeforeEach(func() {
				err = writeByte(buf, uint8(qryPageSize))
				Expect(err).ToNot(HaveOccurred())

				err = writeInt(buf, pageSize)
				Expect(err).ToNot(HaveOccurred())
			})

			It("correctly parses the page size option", func() {
				ps, set := query.PageSize()
				Expect(set).To(BeTrue())
				Expect(ps).To(Equal(pageSize))
			})
		})

		Context("when paging state is set", func() {
			var pagingState = []byte("FOOBAR")

			BeforeEach(func() {
				err = writeByte(buf, uint8(qryPagingState))
				Expect(err).ToNot(HaveOccurred())

				err = writeBytes(buf, pagingState)
				Expect(err).ToNot(HaveOccurred())
			})

			It("correctly parses the paging state option", func() {
				ps, set := query.PagingState()
				Expect(set).To(BeTrue())
				Expect(ps).To(Equal(pagingState))
			})
		})

		Context("when serial consistency is set", func() {
			var serialConsistency = LocalSerial

			BeforeEach(func() {
				err = writeByte(buf, uint8(qrySerialConsistency))
				Expect(err).ToNot(HaveOccurred())

				err = writeShort(buf, uint16(serialConsistency))
				Expect(err).ToNot(HaveOccurred())
			})

			It("correctly parses the serial consistency option", func() {
				sc, set := query.SerialConsistency()
				Expect(set).To(BeTrue())
				Expect(sc).To(Equal(serialConsistency))
			})
		})

		Context("when default timestamp is set", func() {
			var (
				defaultTimestamp = time.Now().Add(time.Hour).Round(time.Microsecond)
				microSeconds     = defaultTimestamp.UnixNano() / int64(time.Microsecond)
			)

			BeforeEach(func() {
				err = writeByte(buf, uint8(qryDefaultTimestamp))
				Expect(err).ToNot(HaveOccurred())

				err = writeLong(buf, microSeconds)
				Expect(err).ToNot(HaveOccurred())
			})

			It("correctly parses the serial consistency option", func() {
				ts, set := query.DefaultTimestamp()
				Expect(set).To(BeTrue())
				Expect(ts).To(Equal(defaultTimestamp))
			})
		})

		Context("when there are unnamed query values", func() {
			var (
				flagSet = qryValues
				values  = [][]byte{[]byte("foo"), []byte("bar")}
			)

			BeforeEach(func() {
				// write flagSet
				err = writeByte(buf, uint8(flagSet))
				Expect(err).ToNot(HaveOccurred())

				// write number of values
				err = writeShort(buf, uint16(len(values)))
				Expect(err).ToNot(HaveOccurred())

				// write actual values
				for _, v := range values {
					err = writeBytes(buf, v)
					Expect(err).ToNot(HaveOccurred())
				}
			})

			It("returns no error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("correctly sets the query values", func() {
				v, set := query.Values()
				Expect(set).To(BeTrue())
				Expect(v).To(Equal(values))
			})

			It("does not set any query value names", func() {
				_, set := query.NamedValues()
				Expect(set).To(BeFalse())
			})
		})

		Context("when there are named query values", func() {
			var (
				flagSet = qryValues | qryNames
				values  = [][]byte{[]byte("foo"), []byte("bar")}
				names   = []string{"one", "two"}
			)

			BeforeEach(func() {
				// write flagSet
				err = writeByte(buf, uint8(flagSet))
				Expect(err).ToNot(HaveOccurred())

				// write number of values
				err = writeShort(buf, uint16(len(values)))
				Expect(err).ToNot(HaveOccurred())

				// write names and values
				for i, v := range values {
					err = writeString(buf, names[i])
					Expect(err).ToNot(HaveOccurred())

					err = writeBytes(buf, v)
					Expect(err).ToNot(HaveOccurred())
				}
			})

			It("returns no error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("correctly sets the query values", func() {
				v, set := query.Values()
				Expect(set).To(BeTrue())
				Expect(v).To(Equal(values))
			})

			It("correctly sets the query value names", func() {
				expected := map[string][]byte{}
				for i, v := range values {
					expected[names[i]] = v
				}

				actual, set := query.NamedValues()
				Expect(set).To(BeTrue())
				Expect(actual).To(Equal(expected))
			})
		})
	})

	Context("when there is nothing to read", func() {
		BeforeEach(func() {
			buf = bytes.NewBuffer([]byte{})
		})

		It("returns an EOF error", func() {
			Expect(err).To(Equal(io.EOF))
		})
	})

	Context("when the query format is invalid", func() {
		BeforeEach(func() {
			buf = bytes.NewBuffer([]byte{})

			stmt := "TOO SHORT"

			err = writeShort(buf, uint16(len(stmt)+10))
			Expect(err).ToNot(HaveOccurred())

			_, err = buf.Write([]byte(stmt))
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns an error", func() {
			Expect(err).To(HaveOccurred())
		})
	})
})
