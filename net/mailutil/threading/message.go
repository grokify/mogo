// Package threading provides generic email thread reconstruction.
// It can reconstruct threading relationships even when In-Reply-To
// and References headers are missing, using subject matching,
// date proximity, and embedded message hints.
package threading

import "time"

// ThreadableMessage is the interface that messages must implement
// for thread reconstruction.
type ThreadableMessage interface {
	// GetMessageID returns the unique message identifier.
	GetMessageID() string

	// GetDate returns the message date.
	GetDate() time.Time

	// GetSubject returns the message subject.
	GetSubject() string

	// GetInReplyTo returns the In-Reply-To header value (may be empty).
	GetInReplyTo() string

	// GetReferences returns the References header values (may be empty).
	GetReferences() []string

	// GetParticipants returns all email addresses involved in the message
	// (From, To, Cc, Bcc).
	GetParticipants() []string

	// GetEmbeddedMessageHints returns hints about embedded/quoted messages
	// that can be used for threading when headers are missing.
	GetEmbeddedMessageHints() []EmbeddedHint

	// SetThreadingInfo is called after reconstruction to provide
	// the computed threading information back to the message.
	SetThreadingInfo(info ThreadingInfo)
}

// EmbeddedHint represents information about a message embedded in the body,
// such as a quoted reply or forwarded message.
type EmbeddedHint struct {
	// SenderPattern is a pattern to match against participant addresses
	// (e.g., "john.smith" or "john.smith@enron.com").
	SenderPattern string

	// Date is the date of the embedded message (if parseable).
	Date time.Time

	// Subject is the subject of the embedded message (if available).
	Subject string

	// Type indicates the type of embedding: "reply", "forward", "quoted".
	Type string
}

// ThreadingInfo contains the computed threading information for a message.
type ThreadingInfo struct {
	// ThreadID is a unique identifier for the thread this message belongs to.
	ThreadID string

	// ParentID is the MessageID of the parent message in the thread.
	// Empty if this is a root message.
	ParentID string

	// References is the reconstructed chain of message IDs leading to this message.
	References []string

	// Depth is the nesting depth in the thread (0 for root messages).
	Depth int
}

// Thread represents a collection of related messages.
type Thread struct {
	// ID is a unique identifier for the thread.
	ID string `json:"id"`

	// Subject is the normalized subject of the thread.
	Subject string `json:"subject"`

	// RootMessageID is the MessageID of the first message in the thread.
	RootMessageID string `json:"root_message_id"`

	// MessageIDs contains all message IDs in the thread, sorted by date.
	MessageIDs []string `json:"message_ids"`

	// Participants contains all unique email addresses in the thread.
	Participants []string `json:"participants"`

	// StartDate is the date of the first message.
	StartDate time.Time `json:"start_date"`

	// EndDate is the date of the last message.
	EndDate time.Time `json:"end_date"`

	// Size is the number of messages in the thread.
	Size int `json:"size"`
}
