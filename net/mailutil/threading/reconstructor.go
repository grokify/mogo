package threading

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
	"time"
)

// Config contains configuration options for thread reconstruction.
type Config struct {
	// MaxParentAge is the maximum age difference for subject-based matching.
	// Default: 7 days.
	MaxParentAge time.Duration

	// RequireParticipantOverlap requires messages to share at least one
	// participant for subject-based matching. Default: true.
	RequireParticipantOverlap bool

	// SubjectNormalizer is a custom function to normalize subjects.
	// If nil, uses DefaultSubjectNormalizer.
	SubjectNormalizer func(string) string
}

// DefaultConfig returns the default reconstruction configuration.
func DefaultConfig() Config {
	return Config{
		MaxParentAge:              7 * 24 * time.Hour,
		RequireParticipantOverlap: true,
		SubjectNormalizer:         nil,
	}
}

// Reconstructor builds thread relationships across messages.
type Reconstructor struct {
	config      Config
	messages    []ThreadableMessage
	byMessageID map[string]ThreadableMessage
	bySubject   map[string][]ThreadableMessage
	threads     map[string]*Thread
	threadInfos map[string]ThreadingInfo
}

// NewReconstructor creates a new thread reconstructor with default config.
func NewReconstructor() *Reconstructor {
	return NewReconstructorWithConfig(DefaultConfig())
}

// NewReconstructorWithConfig creates a new thread reconstructor with custom config.
func NewReconstructorWithConfig(config Config) *Reconstructor {
	return &Reconstructor{
		config:      config,
		messages:    make([]ThreadableMessage, 0),
		byMessageID: make(map[string]ThreadableMessage),
		bySubject:   make(map[string][]ThreadableMessage),
		threads:     make(map[string]*Thread),
		threadInfos: make(map[string]ThreadingInfo),
	}
}

// AddMessage adds a message to the reconstruction pool.
func (r *Reconstructor) AddMessage(msg ThreadableMessage) {
	r.messages = append(r.messages, msg)

	msgID := msg.GetMessageID()
	if msgID != "" {
		r.byMessageID[msgID] = msg
	}

	normSubject := r.normalizeSubject(msg.GetSubject())
	if normSubject != "" {
		r.bySubject[normSubject] = append(r.bySubject[normSubject], msg)
	}
}

// AddMessages adds multiple messages to the reconstruction pool.
func (r *Reconstructor) AddMessages(msgs []ThreadableMessage) {
	for _, msg := range msgs {
		r.AddMessage(msg)
	}
}

// Reconstruct performs thread reconstruction across all messages.
func (r *Reconstructor) Reconstruct() {
	// Sort messages by date for deterministic processing
	sort.Slice(r.messages, func(i, j int) bool {
		return r.messages[i].GetDate().Before(r.messages[j].GetDate())
	})

	// Initialize threading info for all messages
	for _, msg := range r.messages {
		r.threadInfos[msg.GetMessageID()] = ThreadingInfo{}
	}

	// Phase 1: Use existing In-Reply-To/References headers
	r.useExistingHeaders()

	// Phase 2: Match by embedded message hints
	r.matchByEmbeddedHints()

	// Phase 3: Match by subject and date proximity
	r.matchBySubject()

	// Phase 4: Build threads and assign thread IDs
	r.buildThreads()

	// Phase 5: Generate reference chains
	r.generateReferences()

	// Phase 6: Notify messages of their threading info
	r.notifyMessages()
}

// useExistingHeaders uses In-Reply-To and References headers if present.
func (r *Reconstructor) useExistingHeaders() {
	for _, msg := range r.messages {
		msgID := msg.GetMessageID()
		info := r.threadInfos[msgID]

		inReplyTo := msg.GetInReplyTo()
		if inReplyTo != "" {
			if _, ok := r.byMessageID[inReplyTo]; ok {
				info.ParentID = inReplyTo
				r.threadInfos[msgID] = info
			}
		}
	}
}

// matchByEmbeddedHints matches messages using embedded message hints.
func (r *Reconstructor) matchByEmbeddedHints() {
	for _, msg := range r.messages {
		msgID := msg.GetMessageID()
		info := r.threadInfos[msgID]

		if info.ParentID != "" {
			continue // Already have parent
		}

		hints := msg.GetEmbeddedMessageHints()
		for _, hint := range hints {
			parent := r.findParentByHint(msg, hint)
			if parent != nil {
				info.ParentID = parent.GetMessageID()
				r.threadInfos[msgID] = info
				break
			}
		}
	}
}

// findParentByHint finds a parent message based on an embedded hint.
func (r *Reconstructor) findParentByHint(msg ThreadableMessage, hint EmbeddedHint) ThreadableMessage {
	if hint.SenderPattern == "" {
		return nil
	}

	normSubject := r.normalizeSubject(msg.GetSubject())
	candidates := r.bySubject[normSubject]

	senderLower := strings.ToLower(hint.SenderPattern)

	for _, candidate := range candidates {
		if candidate.GetMessageID() == msg.GetMessageID() {
			continue
		}

		// Must be before the current message
		if !candidate.GetDate().Before(msg.GetDate()) {
			continue
		}

		// Check if any participant matches the sender pattern
		for _, participant := range candidate.GetParticipants() {
			participantLower := strings.ToLower(participant)
			if strings.Contains(participantLower, senderLower) ||
				strings.Contains(senderLower, participantLower) {
				// Check date proximity if hint has date
				if !hint.Date.IsZero() {
					dateDiff := candidate.GetDate().Sub(hint.Date)
					if dateDiff < 0 {
						dateDiff = -dateDiff
					}
					if dateDiff < 24*time.Hour {
						return candidate
					}
				} else {
					return candidate
				}
			}
		}
	}

	return nil
}

// matchBySubject matches messages by subject line and date proximity.
func (r *Reconstructor) matchBySubject() {
	for _, msg := range r.messages {
		msgID := msg.GetMessageID()
		info := r.threadInfos[msgID]

		if info.ParentID != "" {
			continue // Already have parent
		}

		// Only match if this looks like a reply
		if !isReplySubject(msg.GetSubject()) {
			continue
		}

		parent := r.findParentBySubject(msg)
		if parent != nil {
			info.ParentID = parent.GetMessageID()
			r.threadInfos[msgID] = info
		}
	}
}

// findParentBySubject finds a parent message based on subject matching.
func (r *Reconstructor) findParentBySubject(msg ThreadableMessage) ThreadableMessage {
	normSubject := r.normalizeSubject(msg.GetSubject())
	candidates := r.bySubject[normSubject]

	var bestMatch ThreadableMessage
	var bestTimeDiff time.Duration = r.config.MaxParentAge

	msgParticipants := toSet(msg.GetParticipants())

	for _, candidate := range candidates {
		if candidate.GetMessageID() == msg.GetMessageID() {
			continue
		}

		// Must be before the current message
		if !candidate.GetDate().Before(msg.GetDate()) {
			continue
		}

		// Check participant overlap if required
		if r.config.RequireParticipantOverlap {
			candidateParticipants := toSet(candidate.GetParticipants())
			if !hasOverlap(msgParticipants, candidateParticipants) {
				continue
			}
		}

		timeDiff := msg.GetDate().Sub(candidate.GetDate())
		if timeDiff < bestTimeDiff {
			bestTimeDiff = timeDiff
			bestMatch = candidate
		}
	}

	return bestMatch
}

// buildThreads groups messages into threads and assigns thread IDs.
func (r *Reconstructor) buildThreads() {
	// Find root messages (no parent) and children
	children := make(map[string][]string) // parentID -> child message IDs

	for _, msg := range r.messages {
		msgID := msg.GetMessageID()
		info := r.threadInfos[msgID]

		if info.ParentID != "" {
			children[info.ParentID] = append(children[info.ParentID], msgID)
		}
	}

	// Build threads from roots
	visited := make(map[string]bool)

	for _, msg := range r.messages {
		msgID := msg.GetMessageID()

		if visited[msgID] {
			continue
		}

		// Find root of this message's thread
		root := r.findRoot(msgID)
		if visited[root] {
			continue
		}

		// Build thread from root
		thread := r.buildThreadFromRoot(root, children, visited)
		if thread != nil {
			r.threads[thread.ID] = thread

			// Update thread ID for all messages in thread
			for _, mid := range thread.MessageIDs {
				info := r.threadInfos[mid]
				info.ThreadID = thread.ID
				r.threadInfos[mid] = info
			}
		}
	}
}

// findRoot finds the root message ID for a given message.
func (r *Reconstructor) findRoot(msgID string) string {
	seen := make(map[string]bool)
	current := msgID

	for {
		if seen[current] {
			return current // Cycle detected
		}
		seen[current] = true

		info := r.threadInfos[current]
		if info.ParentID == "" {
			return current
		}
		current = info.ParentID
	}
}

// buildThreadFromRoot builds a thread starting from a root message.
func (r *Reconstructor) buildThreadFromRoot(rootID string, children map[string][]string, visited map[string]bool) *Thread {
	rootMsg, ok := r.byMessageID[rootID]
	if !ok {
		return nil
	}

	normSubject := r.normalizeSubject(rootMsg.GetSubject())
	threadID := generateThreadID(normSubject, rootMsg.GetDate())

	thread := &Thread{
		ID:            threadID,
		Subject:       normSubject,
		RootMessageID: rootID,
		MessageIDs:    []string{},
		StartDate:     rootMsg.GetDate(),
		EndDate:       rootMsg.GetDate(),
	}

	participantSet := make(map[string]bool)

	// Collect all messages in thread using BFS
	var collectMessages func(msgID string, depth int)
	collectMessages = func(msgID string, depth int) {
		if visited[msgID] {
			return
		}
		visited[msgID] = true

		msg, ok := r.byMessageID[msgID]
		if !ok {
			return
		}

		thread.MessageIDs = append(thread.MessageIDs, msgID)

		// Update depth
		info := r.threadInfos[msgID]
		info.Depth = depth
		r.threadInfos[msgID] = info

		// Update dates
		msgDate := msg.GetDate()
		if msgDate.Before(thread.StartDate) {
			thread.StartDate = msgDate
		}
		if msgDate.After(thread.EndDate) {
			thread.EndDate = msgDate
		}

		// Collect participants
		for _, p := range msg.GetParticipants() {
			participantSet[strings.ToLower(p)] = true
		}

		// Process children
		for _, childID := range children[msgID] {
			collectMessages(childID, depth+1)
		}
	}

	collectMessages(rootID, 0)

	// Sort message IDs by date
	sort.Slice(thread.MessageIDs, func(i, j int) bool {
		mi := r.byMessageID[thread.MessageIDs[i]]
		mj := r.byMessageID[thread.MessageIDs[j]]
		return mi.GetDate().Before(mj.GetDate())
	})

	// Collect participants
	for p := range participantSet {
		thread.Participants = append(thread.Participants, p)
	}
	sort.Strings(thread.Participants)

	thread.Size = len(thread.MessageIDs)

	return thread
}

// generateReferences generates References header chains for each message.
func (r *Reconstructor) generateReferences() {
	for _, msg := range r.messages {
		msgID := msg.GetMessageID()
		refs := r.buildReferenceChain(msgID)

		info := r.threadInfos[msgID]
		info.References = refs
		r.threadInfos[msgID] = info
	}
}

// buildReferenceChain builds the chain of message IDs leading to this message.
func (r *Reconstructor) buildReferenceChain(msgID string) []string {
	var refs []string
	seen := make(map[string]bool)

	current := msgID
	for {
		info := r.threadInfos[current]
		if info.ParentID == "" {
			break
		}
		if seen[info.ParentID] {
			break // Avoid cycles
		}
		seen[info.ParentID] = true
		refs = append([]string{info.ParentID}, refs...)
		current = info.ParentID
	}

	return refs
}

// notifyMessages calls SetThreadingInfo on each message.
func (r *Reconstructor) notifyMessages() {
	for _, msg := range r.messages {
		info := r.threadInfos[msg.GetMessageID()]
		msg.SetThreadingInfo(info)
	}
}

// normalizeSubject normalizes a subject line for comparison.
func (r *Reconstructor) normalizeSubject(subject string) string {
	if r.config.SubjectNormalizer != nil {
		return r.config.SubjectNormalizer(subject)
	}
	return DefaultSubjectNormalizer(subject)
}

// GetThreads returns all reconstructed threads.
func (r *Reconstructor) GetThreads() []*Thread {
	threads := make([]*Thread, 0, len(r.threads))
	for _, thread := range r.threads {
		threads = append(threads, thread)
	}

	sort.Slice(threads, func(i, j int) bool {
		return threads[i].StartDate.Before(threads[j].StartDate)
	})

	return threads
}

// GetThreadingInfo returns the threading info for a message.
func (r *Reconstructor) GetThreadingInfo(msgID string) (ThreadingInfo, bool) {
	info, ok := r.threadInfos[msgID]
	return info, ok
}

// Stats returns threading statistics.
func (r *Reconstructor) Stats() Stats {
	stats := Stats{
		TotalMessages:  len(r.messages),
		TotalThreads:   len(r.threads),
		UniqueSubjects: len(r.bySubject),
	}

	for _, info := range r.threadInfos {
		if info.ParentID != "" {
			stats.MessagesWithParent++
		}
		if len(info.References) > 0 {
			stats.MessagesWithRefs++
		}
	}

	for _, thread := range r.threads {
		switch {
		case thread.Size == 1:
			stats.SingleMessageThreads++
		case thread.Size <= 5:
			stats.SmallThreads++
		case thread.Size <= 20:
			stats.MediumThreads++
		default:
			stats.LargeThreads++
		}
	}

	return stats
}

// Stats contains statistics about thread reconstruction.
type Stats struct {
	TotalMessages        int `json:"total_messages"`
	TotalThreads         int `json:"total_threads"`
	UniqueSubjects       int `json:"unique_subjects"`
	MessagesWithParent   int `json:"messages_with_parent"`
	MessagesWithRefs     int `json:"messages_with_refs"`
	SingleMessageThreads int `json:"single_message_threads"`
	SmallThreads         int `json:"small_threads"`  // 2-5 messages
	MediumThreads        int `json:"medium_threads"` // 6-20 messages
	LargeThreads         int `json:"large_threads"`  // 21+ messages
}

// DefaultSubjectNormalizer removes common reply/forward prefixes.
func DefaultSubjectNormalizer(subject string) string {
	s := strings.TrimSpace(subject)
	for {
		trimmed := s
		for _, prefix := range []string{
			"Re: ", "RE: ", "re: ", "Re:", "RE:", "re:",
			"Fw: ", "FW: ", "fw: ", "Fw:", "FW:", "fw:",
			"Fwd: ", "FWD: ", "fwd: ", "Fwd:", "FWD:", "fwd:",
		} {
			if strings.HasPrefix(trimmed, prefix) {
				trimmed = strings.TrimSpace(trimmed[len(prefix):])
				break
			}
		}
		if trimmed == s {
			break
		}
		s = trimmed
	}
	return s
}

// isReplySubject checks if a subject indicates a reply or forward.
func isReplySubject(subject string) bool {
	lower := strings.ToLower(strings.TrimSpace(subject))
	return strings.HasPrefix(lower, "re:") ||
		strings.HasPrefix(lower, "fw:") ||
		strings.HasPrefix(lower, "fwd:")
}

// generateThreadID creates a unique thread ID.
func generateThreadID(subject string, date time.Time) string {
	data := subject + date.Format("2006-01-02")
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:8])
}

// toSet converts a slice to a set.
func toSet(items []string) map[string]bool {
	set := make(map[string]bool)
	for _, item := range items {
		set[strings.ToLower(item)] = true
	}
	return set
}

// hasOverlap checks if two sets have any common elements.
func hasOverlap(a, b map[string]bool) bool {
	for k := range a {
		if b[k] {
			return true
		}
	}
	return false
}
