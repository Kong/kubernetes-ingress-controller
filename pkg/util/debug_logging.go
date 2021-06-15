package util

import (
	"io"
	"time"

	"github.com/sirupsen/logrus"
)

// -----------------------------------------------------------------------------
// Public - Reduced Redudancy Debug Logging
// -----------------------------------------------------------------------------

// MakeDebugLoggerWithReducedRedudancy is a logrus.Logger that "stifles" repetitive logs.
//
// The "stifling" mechanism is triggered by one of two conditions the result of which is
// that the "stifled" log entry will be dropped entirely.
//
// The conditions checked are:
//
//  1. This logger will drop log entries where an identical log entry has posted within the
//     last "redundantLogEntryBackoff". For example, you could set this to "time.Second * 3"
//     and the result would be that if the logger had already logged an identical message
//     within the previous 3 seconds it will be dropped.
//
//  2. This logger will "stifle" redudant entries which are logged consecutively a number of
//     times equal to the provided "redudantLogEntryAllowedConsecutively" number. For example,
//     you could set this to 3 and then if the last 3 log entries emitted were the same message
//     further entries of the same message would be dropped.
//
// The caller can choose to set either argument to "0" to disable that check, but setting both
// to zero will result in no redundancy reduction.
//
// NOTE: Please consider this logger a "debug" only logging implementation.
//       This logger was originally created to help reduce the noise coming from the controller
//       during integration tests for better human readability, so keep in mind it was built for
//       testing environments if you're currently reading this and you're considering using it
//       somewhere that would produce production environment logs: there's significant
//       performance overhead triggered by the logging hooks this adds.
func MakeDebugLoggerWithReducedRedudancy(writer io.Writer, formatter logrus.Formatter,
	redudantLogEntryAllowedConsecutively int, redundantLogEntryBackoff time.Duration,
) *logrus.Logger {
	// setup the logger with debug level logging
	log := logrus.New()
	log.Level = logrus.DebugLevel
	log.Formatter = formatter
	log.Out = writer

	// set up the reduced redudancy logging hook
	log.Hooks.Add(newReducedRedundancyLogHook(redundantLogEntryBackoff, redudantLogEntryAllowedConsecutively, formatter))
	return log
}

// -----------------------------------------------------------------------------
// Private - Reduced Redudancy Debug Logging
// -----------------------------------------------------------------------------

// reducedRedudancyLogHook is a logrus.Hook that reduces redudant log entries.
type reducedRedudancyLogHook struct {
	backoff            time.Duration
	consecutiveAllowed int
	consecutivePosted  int
	formatter          logrus.Formatter
	lastMessage        string
	timeWindow         map[string]bool
	timeWindowStart    time.Time
}

func newReducedRedundancyLogHook(
	backoff time.Duration,
	consecutive int,
	formatter logrus.Formatter,
) *reducedRedudancyLogHook {
	return &reducedRedudancyLogHook{
		backoff:            backoff,
		consecutiveAllowed: consecutive,
		formatter:          formatter,
		timeWindowStart:    time.Now(),
		timeWindow:         map[string]bool{},
	}
}

func (r *reducedRedudancyLogHook) Fire(entry *logrus.Entry) error {
	defer func() { r.lastMessage = entry.Message }()

	// to make this hook work we override the logger formatter to the nilFormatter
	// for some entries, but we also need to reset it here to ensure the default.
	entry.Logger.Formatter = r.formatter

	// if the current entry has the exact same message as the last entry, check the
	// consecutive posting rules for this entry to see whether it should be dropped.
	if r.consecutiveAllowed > 0 && entry.Message == r.lastMessage {
		r.consecutivePosted++
		if r.consecutivePosted >= r.consecutiveAllowed {
			entry.Logger.SetFormatter(&nilFormatter{})
			return nil
		}
	} else {
		r.consecutivePosted = 0
	}

	// determine whether or not the previous time window is still valid and if not create
	// a new time window and return.
	if time.Now().After(r.timeWindowStart.Add(r.backoff)) {
		r.timeWindow = map[string]bool{}
		r.timeWindowStart = time.Now()
		return nil
	}

	// if we're here then the time window is still valid, we need to determine if the
	// current entry would be considered redundant during this time window.
	// if the entry has not yet been seen during this time window, we record it so that
	// future checks can find it.
	if _, ok := r.timeWindow[entry.Message]; ok {
		entry.Logger.SetFormatter(&nilFormatter{})
	}
	r.timeWindow[entry.Message] = true

	return nil
}

func (r *reducedRedudancyLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// -----------------------------------------------------------------------------
// Private - Nil Logging Formatter
// -----------------------------------------------------------------------------

type nilFormatter struct{}

func (n *nilFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return nil, nil
}
