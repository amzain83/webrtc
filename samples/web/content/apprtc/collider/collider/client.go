// Copyright (c) 2014 The WebRTC project authors. All Rights Reserved.
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

package collider

import (
	"errors"
	"io"
	"log"
)

type client struct {
	id string
	// rwc is the interface to access the websocket connection.
	// It is set after the client registers with the server.
	rwc io.ReadWriteCloser
	// msgs is the queued messages sent from this client.
	msgs []string
}

func newClient(id string) *client {
	return &client{id: id}
}

// register binds the ReadWriteCloser to the client if it's not done yet.
func (c *client) register(rwc io.ReadWriteCloser) error {
	if c.rwc != nil {
		log.Printf("Not registering because the client %s already has a connection", c.id)
		return errors.New("Duplicated registration")
	}
	c.rwc = rwc
	return nil
}

// Adds a message to the client's message queue.
func (c *client) enqueue(msg string) {
	c.msgs = append(c.msgs, msg)
}

// sendQueued the queued messages to the other client.
func (c *client) sendQueued(other *client) error {
	if c.id == other.id || other.rwc == nil {
		return errors.New("Invalid client")
	}
	for _, m := range c.msgs {
		sendServerMsg(other.rwc, m)
	}
	c.msgs = nil
	log.Printf("Sent queued messages from %s to %s", c.id, other.id)
	return nil
}

// send sends the message to the other client if the other client has registered,
// or queues the message otherwise.
func (c *client) send(other *client, msg string) error {
	if c.id == other.id {
		return errors.New("Invalid client")
	}
	if other.rwc != nil {
		return sendServerMsg(other.rwc, msg)
	}
	c.enqueue(msg)
	return nil
}

// close closes the ReadWriteCloser if it exists.
func (c *client) close() {
	if c.rwc != nil {
		c.rwc.Close()
		c.rwc = nil
	}
}
