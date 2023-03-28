package types

import (
	"errors"

	"github.com/tendermint/tendermint/crypto/ed25519"
)

// ValidateBasic performs basic validation of a block.
func (b *Block) ValidateBasic() error {
	err := b.SignedHeader.ValidateBasic()
	if err != nil {
		return err
	}

	err = b.Data.ValidateBasic()
	if err != nil {
		return err
	}

	return nil
}

// ValidateBasic performs basic validation of a header.
func (h *Header) ValidateBasic() error {
	if len(h.ProposerAddress) == 0 {
		return errors.New("no proposer address")
	}

	return nil
}

// ValidateBasic performs basic validation of block data.
// Actually it's a placeholder, because nothing is checked.
func (d *Data) ValidateBasic() error {
	return nil
}

// ValidateBasic performs basic validation of a commit.
func (c *Commit) ValidateBasic() error {
	if len(c.Signatures) == 0 {
		return errors.New("no signatures")
	}
	return nil
}

// ValidateBasic performs basic validation of a signed header.
func (h *SignedHeader) ValidateBasic() error {
	err := h.Header.ValidateBasic()
	if err != nil {
		return err
	}

	err = h.Commit.ValidateBasic()
	if err != nil {
		return err
	}

	err = h.Validators.ValidateBasic()
	if err != nil {
		return err
	}

	// Make sure there is exactly one signature
	if len(h.Commit.Signatures) != 1 {
		return errors.New("expected exactly one signature")
	}

	signature := h.Commit.Signatures[0]
	proposer := h.Validators.GetProposer()
	var pubKey ed25519.PubKey = proposer.PubKey.Bytes()
	msg, err := h.Header.MarshalBinary()
	if err != nil {
		return errors.New("signature verification failed, unable to marshal header")
	}
	if !pubKey.VerifySignature(msg, signature) {
		return errors.New("signature verification failed")
	}

	return nil
}
