package bugtracker

import (
	"github.com/openela/mothership/base/forge"
)

type Options struct {
	// Labels are the labels to apply to the ticket.
	// All bug trackers might not support labels, but it's a common feature.
	Labels []string
}

type Bugtracker interface {
	// CreateTicket creates a ticket in the bug tracker.
	// Returns the ticket ID or an error.
	CreateTicket(auth *forge.Authenticator, title string, body string, opts Options) (string, error)

	// EditTicket edits the ticket in the bug tracker.
	// Returns an error if the ticket could not be edited.
	EditTicket(auth *forge.Authenticator, ticketID string, title string, body string, opts Options) error

	// CloseTicket closes the ticket in the bug tracker.
	// Returns an error if the ticket could not be closed.
	CloseTicket(auth *forge.Authenticator, ticketID string) error

	// TicketURI returns the URI to the ticket in the bug tracker.
	// Returns an error if the URI could not be generated.
	TicketURI(ticketID string) (string, error)

	// URIToTicket returns the ticket ID from the URI.
	// Returns an error if the ticket ID could not be extracted.
	URIToTicket(uri string) (string, error)

	// GetAuthenticator returns an authenticator for the bug tracker.
	GetAuthenticator() (*forge.Authenticator, error)
}
