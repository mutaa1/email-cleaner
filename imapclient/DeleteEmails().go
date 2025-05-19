package imapclient

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

func DeleteEmails(c *client.Client, uids []uint32) error {
	seqset := new(imap.SeqSet)
	for _, uid := range uids {
		seqset.AddNum(uid)
	}

	item := imap.FormatFlagsOp(imap.AddFlags, true)
	flags := []interface{}{imap.DeletedFlag}

	// Mark messages as deleted
	if err := c.UidStore(seqset, item, flags, nil); err != nil {
		return err
	}

	// Permanently remove deleted messages
	return c.Expunge(nil)
}
