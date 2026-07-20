package main

import (
	"fmt"
	"net"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

// Sync journal style: real work is one readable line; repeated "nothing to do"
// is counted and printed as a summary when something actually changes (or the
// first skip of a quiet streak, so idle still proves sync is alive).

type syncSkipBucket struct {
	mu      sync.Mutex
	count   int
	reasons map[string]int
	since   time.Time
	// firstLogged is true after the opening line of this streak.
	firstLogged bool
}

var syncSkips = &syncSkipBucket{reasons: map[string]int{}}

func (b *syncSkipBucket) note(reason string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.count == 0 {
		b.since = time.Now()
		b.firstLogged = false
		b.reasons = map[string]int{}
	}
	b.count++
	b.reasons[reason]++
	if !b.firstLogged {
		b.firstLogged = true
		fmt.Fprintf(os.Stderr, "writerdeck-server: sync: nothing to do (%s) — notes match last sync\n", reason)
		return
	}
	// Later skips stay quiet until flush (real sync) or a periodic reminder.
}

// flush writes a summary of coalesced skips (beyond the first) and resets.
// Call before logging a real sync result so the journal reads as a story.
func (b *syncSkipBucket) flush() {
	b.mu.Lock()
	defer b.mu.Unlock()
	extra := b.count - 1 // first already printed
	if extra <= 0 {
		b.count = 0
		b.reasons = map[string]int{}
		b.firstLogged = false
		return
	}
	fmt.Fprintf(os.Stderr, "writerdeck-server: sync: nothing to do ×%d more (%s) since %s\n",
		extra, formatReasonCounts(b.reasons), b.since.Format("15:04:05"))
	b.count = 0
	b.reasons = map[string]int{}
	b.firstLogged = false
}

func formatReasonCounts(m map[string]int) string {
	if len(m) == 0 {
		return ""
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		n := m[k]
		if n == 1 {
			parts = append(parts, k)
		} else {
			parts = append(parts, fmt.Sprintf("%s×%d", k, n))
		}
	}
	return strings.Join(parts, ", ")
}

func logSyncIdle(reason string) {
	syncSkips.note(reason)
}

func logSyncChanged(reason string, details []string) {
	syncSkips.flush()
	if len(details) == 0 {
		fmt.Fprintf(os.Stderr, "writerdeck-server: sync: checked GitHub (%s) — still nothing to change\n", reason)
		return
	}
	const maxShow = 8
	show := details
	more := 0
	if len(show) > maxShow {
		more = len(show) - maxShow
		show = show[:maxShow]
	}
	msg := strings.Join(show, "; ")
	if more > 0 {
		msg += fmt.Sprintf("; …+%d more", more)
	}
	fmt.Fprintf(os.Stderr, "writerdeck-server: sync: %s (%s) — %d file%s\n",
		msg, reason, len(details), plural(len(details)))
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

// --- WebSocket connect coalesce (many tabs from one phone/laptop) ---

const connectLogWindow = 1500 * time.Millisecond

type connectBurst struct {
	mu    sync.Mutex
	host  string
	count int
	timer *time.Timer
}

var connectLog = &connectBurst{}

func logClientConnected(remoteAddr string) {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		host = remoteAddr
	}
	connectLog.mu.Lock()
	defer connectLog.mu.Unlock()
	if connectLog.host != "" && connectLog.host != host {
		connectLog.flushLocked()
	}
	if connectLog.host == "" {
		connectLog.host = host
		connectLog.count = 1
		connectLog.timer = time.AfterFunc(connectLogWindow, func() {
			connectLog.mu.Lock()
			defer connectLog.mu.Unlock()
			connectLog.flushLocked()
		})
		return
	}
	connectLog.count++
	if connectLog.timer != nil {
		connectLog.timer.Reset(connectLogWindow)
	}
}

func (c *connectBurst) flushLocked() {
	if c.count <= 0 || c.host == "" {
		c.host = ""
		c.count = 0
		return
	}
	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
	}
	if c.count == 1 {
		fmt.Fprintf(os.Stderr, "writerdeck-server: client connected %s\n", c.host)
	} else {
		fmt.Fprintf(os.Stderr, "writerdeck-server: client connected %s (%d tabs)\n", c.host, c.count)
	}
	c.host = ""
	c.count = 0
}
