package proto

import "io"

// VersionDir represents the version AND direction of a CQL frame.
type VersionDir uint8

const (
	directionMask VersionDir = 0x80
)

type versioner struct {
	framers map[VersionDir]Framer
}

func NewVersioner() *versioner {
	return &versioner{
		framers: map[VersionDir]Framer{},
	}
}

func (v *versioner) Version(in io.Reader) (Framer, error) {
	var version VersionDir
	if err := readVersionDir(in, &version); err != nil {
		return nil, err
	}

	framer, found := v.framers[version]
	if !found {
		return nil, errUnsupportedProtocolVersion
	}

	return framer, nil
}

func (v *versioner) SetRequestFramer(version Version, framer Framer) {
	v.framers[VersionDir(version)] = framer
}

func (v *versioner) SetResponseFramer(version Version, framer Framer) {
	v.framers[VersionDir(version)|directionMask] = framer
}
