package scanner

import (
	"io"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/mutaa1/email-cleaner/models"
)

func ScanInbox(c *client.Client) ([]models.Email, error) {
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		return nil, err
	}

	seqSet := new(imap.SeqSet)
	from := uint32(1)
	to := mbox.Messages
	if mbox.Messages > 10 {
		from = mbox.Messages - 9
	}
	seqSet.AddRange(from, to)

	section := &imap.BodySectionName{}
	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	go func() {
		done <- c.Fetch(seqSet, []imap.FetchItem{imap.FetchEnvelope, section.FetchItem()}, messages)
	}()

	var emails []models.Email

	for msg := range messages {
		r := msg.GetBody(section)
		if r == nil {
			continue
		}

		mr, err := mail.CreateReader(r)
		if err != nil {
			continue
		}

		header := mr.Header
		from, _ := header.AddressList("From")
		subject, _ := header.Subject()
		unsubscribe, _ := header.Text("List-Unsubscribe")

		emails = append(emails, models.Email{
			From:        from[0].Address,
			Subject:     subject,
			Unsubscribe: strings.Trim(unsubscribe, "<>"),
		})

		for {
			_, err := mr.NextPart()
			if err == io.EOF {
				break
			}
		}
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return emails, nil
}
