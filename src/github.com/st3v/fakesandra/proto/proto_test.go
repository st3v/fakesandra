package proto_test

import (
	"bytes"
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/st3v/fakesandra/proto"
)

var _ = Describe("proto.Read", func() {
	var (
		input  *proto.Frame
		output proto.Frame
		err    error
	)

	JustBeforeEach(func() {
		buf := new(bytes.Buffer)

		if input != nil {
			Expect(proto.Write(buf, *input)).ToNot(HaveOccurred())
		}

		output = proto.Frame{}
		err = proto.Read(buf, &output)
	})

	Context("when there is no frame to read", func() {
		BeforeEach(func() {
			input = nil
		})

		It("returns an EOF error", func() {
			Expect(err).To(Equal(io.EOF))
		})
	})

	Context("when there is an empty frame to read", func() {
		BeforeEach(func() {
			input = new(proto.Frame)
		})

		It("does not return an error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns the correct frame", func() {
			Expect(output).To(Equal(proto.NewFrame(
				proto.FrameHeader{},
				proto.FrameBody{},
			)))
		})
	})

	Context("when there is a non-empty frame to read", func() {
		BeforeEach(func() {
			body := []byte("Foobar")
			f := proto.NewFrame(
				proto.FrameHeader{
					Version: 2,
					Flags:   3,
					Stream:  4,
					OpCode:  5,
					Length:  uint32(len(body)),
				},
				body,
			)
			input = &f
		})

		It("does not return an error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns the correct frame", func() {
			Expect(output).To(Equal(*input))
		})
	})
})
