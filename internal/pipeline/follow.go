package pipeline

import (
	"context"
	"io"
	"os"
	"time"
)

// followInterval is the polling interval used when tailing a file.
const followInterval = 250 * time.Millisecond

// follower wraps an *os.File and implements io.Reader by polling for new
// content, similar to `tail -f`. It blocks until new data is available or
// the context is cancelled.
type follower struct {
	file   *os.File
	ctx    context.Context
	offset int64
}

// newFollower creates a follower that starts reading from the current end of
// the file so that only new lines written after startup are emitted.
func newFollower(ctx context.Context, path string) (*follower, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// Seek to end so we tail only new content.
	offset, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		f.Close()
		return nil, err
	}

	return &follower{
		file:   f,
		ctx:    ctx,
		offset: offset,
	}, nil
}

// Read implements io.Reader. It polls the underlying file for new bytes,
// sleeping between attempts. It returns io.EOF only when the context is done.
func (fw *follower) Read(p []byte) (int, error) {
	for {
		// Attempt to read available bytes.
		n, err := fw.file.Read(p)
		if n > 0 {
			fw.offset += int64(n)
			return n, nil
		}

		// Real error — surface it.
		if err != nil && err != io.EOF {
			return 0, err
		}

		// No new data yet — check if we should stop.
		select {
		case <-fw.ctx.Done():
			return 0, io.EOF
		default:
		}

		// Sleep before the next poll.
		select {
		case <-fw.ctx.Done():
			return 0, io.EOF
		case <-time.After(followInterval):
		}

		// Re-seek to our known offset in case the OS buffered position drifted.
		if _, seekErr := fw.file.Seek(fw.offset, io.SeekStart); seekErr != nil {
			return 0, seekErr
		}
	}
}

// Close releases the underlying file descriptor.
func (fw *follower) Close() error {
	return fw.file.Close()
}
