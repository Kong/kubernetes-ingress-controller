package conditions

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// Option is a functional option for specifying conditions to match.
type Option func(condition metav1.Condition) bool

func Not(option Option) Option {
	return func(condition metav1.Condition) bool {
		return !option(condition)
	}
}

// WithType returns a ConditionOption that matches conditions with the given type.
func WithType(conditionType string) Option {
	return func(condition metav1.Condition) bool {
		return condition.Type == conditionType
	}
}

// WithStatus returns a ConditionOption that matches conditions with the given status.
func WithStatus(status metav1.ConditionStatus) Option {
	return func(condition metav1.Condition) bool {
		return condition.Status == status
	}
}

// WithReason returns a ConditionOption that matches conditions with the given reason.
func WithReason(reason string) Option {
	return func(condition metav1.Condition) bool {
		return condition.Reason == reason
	}
}

// WithMessage returns a ConditionOption that matches conditions with the given message.
func WithMessage(message string) Option {
	return func(condition metav1.Condition) bool {
		return condition.Message == message
	}
}

// WithLastTransitionTime returns a ConditionOption that matches conditions with the given last transition time.
func WithLastTransitionTime(lastTransitionTime metav1.Time) Option {
	return func(condition metav1.Condition) bool {
		return condition.LastTransitionTime == lastTransitionTime
	}
}

// WithCondition returns a ConditionOption that matches conditions with the given condition.
func WithCondition(condition metav1.Condition) Option {
	return func(c metav1.Condition) bool {
		return c == condition
	}
}

// Contain returns true if the given conditions slice contains a condition that matches all of the given options.
func Contain(conditions []metav1.Condition, options ...Option) bool {
	for _, condition := range conditions {
		matches := true
		for _, option := range options {
			if !option(condition) {
				matches = false
				break
			}
		}
		if matches {
			return true
		}
	}
	return false
}
