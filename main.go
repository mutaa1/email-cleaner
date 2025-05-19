package main

import (
	"fmt"
	"log"

	"github.com/emersion/go-imap/client"

	"github.com/mutaa1/email-cleaner/imapclient"
	"github.com/mutaa1/email-cleaner/scanner"
)

func main() {
	email := "denzmutash@gmail.com"
	password := "oeno bgav mdbg jsbw"

	// Connect to server
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Logout()

	// Login
	if err := c.Login(email, password); err != nil {
		log.Fatal(err)
	}

	// Select mailbox
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Flags for INBOX: %v\n", mbox.Flags)

	// Scan inbox for emails with unsubscribe info (your scanner.ScanInbox returns []Email with UIDs)
	scanned, err := scanner.ScanInbox(c)
	if err != nil {
		log.Fatal(err)
	}

	var uidsToDelete []uint32
	for _, email := range scanned {
		fmt.Printf("From: %s\nSubject: %s\nUnsubscribe: %s\n\n", email.From, email.Subject, email.Unsubscribe)
		// Decide if you want to delete this email; here we delete if Unsubscribe link exists
		if email.Unsubscribe != "" {
			uidsToDelete = append(uidsToDelete, email.UID) // make sure your Email struct has UID field
		}
	}

	if len(uidsToDelete) == 0 {
		fmt.Println("No emails to delete.")
		return
	}

	// Delete emails in bulk
	if err := imapclient.DeleteEmails(c, uidsToDelete); err != nil {
		log.Fatalf("Failed to delete emails: %v", err)
	}

	fmt.Printf("🧹 Deleted %d emails.\n", len(uidsToDelete))
}
