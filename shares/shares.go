package shares

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/rollkit/rollkit/appconsts"
	appns "github.com/rollkit/rollkit/namespace"
)

// Share contains the raw share data (including namespace ID).
type Share struct {
	data []byte
}

func (s *Share) Namespace() (appns.Namespace, error) {
	if len(s.data) < appns.NamespaceSize {
		panic(fmt.Sprintf("share %s is too short to contain a namespace", s))
	}
	return appns.From(s.data[:appns.NamespaceSize])
}

func (s *Share) InfoByte() (InfoByte, error) {
	if len(s.data) < appns.NamespaceSize+appconsts.ShareInfoBytes {
		return 0, fmt.Errorf("share %s is too short to contain an info byte", s)
	}
	// the info byte is the first byte after the namespace
	unparsed := s.data[appns.NamespaceSize]
	return ParseInfoByte(unparsed)
}

func NewShare(data []byte) (*Share, error) {
	if err := validateSize(data); err != nil {
		return nil, err
	}
	return &Share{data}, nil
}

func (s *Share) Validate() error {
	return validateSize(s.data)
}

func validateSize(data []byte) error {
	if len(data) != appconsts.ShareSize {
		return fmt.Errorf("share data must be %d bytes, got %d", appconsts.ShareSize, len(data))
	}
	return nil
}

func (s *Share) Len() int {
	return len(s.data)
}

func (s *Share) Version() (uint8, error) {
	infoByte, err := s.InfoByte()
	if err != nil {
		return 0, err
	}
	return infoByte.Version(), nil
}

func (s *Share) DoesSupportVersions(supportedShareVersions []uint8) error {
	ver, err := s.Version()
	if err != nil {
		return err
	}
	if !bytes.Contains(supportedShareVersions, []byte{ver}) {
		return fmt.Errorf("unsupported share version %v is not present in the list of supported share versions %v", ver, supportedShareVersions)
	}
	return nil
}

// IsSequenceStart returns true if this is the first share in a sequence.
func (s *Share) IsSequenceStart() (bool, error) {
	infoByte, err := s.InfoByte()
	if err != nil {
		return false, err
	}
	return infoByte.IsSequenceStart(), nil
}

// IsCompactShare returns true if this is a compact share.
func (s Share) IsCompactShare() (bool, error) {
	ns, err := s.Namespace()
	if err != nil {
		return false, err
	}
	isCompact := ns.IsTx() || ns.IsPayForBlob()
	return isCompact, nil
}

// SequenceLen returns the sequence length of this *share and optionally an
// error. It returns 0, nil if this is a continuation share (i.e. doesn't
// contain a sequence length).
func (s *Share) SequenceLen() (sequenceLen uint32, err error) {
	isSequenceStart, err := s.IsSequenceStart()
	if err != nil {
		return 0, err
	}
	if !isSequenceStart {
		return 0, nil
	}

	start := appconsts.NamespaceSize + appconsts.ShareInfoBytes
	end := start + appconsts.SequenceLenBytes
	if len(s.data) < end {
		return 0, fmt.Errorf("share %s with length %d is too short to contain a sequence length",
			s, len(s.data))
	}
	return binary.BigEndian.Uint32(s.data[start:end]), nil
}

// IsPadding returns whether this *share is padding or not.
func (s *Share) IsPadding() (bool, error) {
	isNamespacePadding, err := s.isNamespacePadding()
	if err != nil {
		return false, err
	}
	isTailPadding, err := s.isTailPadding()
	if err != nil {
		return false, err
	}
	isReservedPadding, err := s.isReservedPadding()
	if err != nil {
		return false, err
	}
	return isNamespacePadding || isTailPadding || isReservedPadding, nil
}

func (s *Share) isNamespacePadding() (bool, error) {
	isSequenceStart, err := s.IsSequenceStart()
	if err != nil {
		return false, err
	}
	sequenceLen, err := s.SequenceLen()
	if err != nil {
		return false, err
	}

	return isSequenceStart && sequenceLen == 0, nil
}

func (s *Share) isTailPadding() (bool, error) {
	ns, err := s.Namespace()
	if err != nil {
		return false, err
	}
	return ns.IsTailPadding(), nil
}

func (s *Share) isReservedPadding() (bool, error) {
	ns, err := s.Namespace()
	if err != nil {
		return false, err
	}
	return ns.IsReservedPadding(), nil
}

func (s *Share) ToBytes() []byte {
	return s.data
}

// RawData returns the raw share data. The raw share data does not contain the
// namespace ID, info byte, sequence length, or reserved bytes.
func (s *Share) RawData() (rawData []byte, err error) {
	if len(s.data) < s.rawDataStartIndex() {
		return rawData, fmt.Errorf("share %s is too short to contain raw data", s)
	}

	return s.data[s.rawDataStartIndex():], nil
}

func (s *Share) rawDataStartIndex() int {
	isStart, err := s.IsSequenceStart()
	if err != nil {
		panic(err)
	}
	isCompact, err := s.IsCompactShare()
	if err != nil {
		panic(err)
	}
	if isStart && isCompact {
		return appconsts.NamespaceSize + appconsts.ShareInfoBytes + appconsts.SequenceLenBytes + appconsts.CompactShareReservedBytes
	} else if isStart && !isCompact {
		return appconsts.NamespaceSize + appconsts.ShareInfoBytes + appconsts.SequenceLenBytes
	} else if !isStart && isCompact {
		return appconsts.NamespaceSize + appconsts.ShareInfoBytes + appconsts.CompactShareReservedBytes
	} else if !isStart && !isCompact {
		return appconsts.NamespaceSize + appconsts.ShareInfoBytes
	} else {
		panic(fmt.Sprintf("unable to determine the rawDataStartIndex for share %s", s.data))
	}
}

func ToBytes(shares []Share) (bytes [][]byte) {
	bytes = make([][]byte, len(shares))
	for i, share := range shares {
		bytes[i] = []byte(share.data)
	}
	return bytes
}

func FromBytes(bytes [][]byte) (shares []Share, err error) {
	for _, b := range bytes {
		share, err := NewShare(b)
		if err != nil {
			return nil, err
		}
		shares = append(shares, *share)
	}
	return shares, nil
}
