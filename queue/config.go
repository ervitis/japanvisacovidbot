package queue

var (
	Queue *gobbit
)

const (
	NewCovidEntryEvent = "NEW_COVID_ENTRY_EVENT"
)

func init() {
	Queue = New()
}

func AllTopics() []string {
	return []string{
		NewCovidEntryEvent,
	}
}
