package constant

type From int

const (
	ADMIN From = iota
	KANBOARD
)

type Gender uint

const (
	MALE   Gender = 1
	FEMALE Gender = 2
)

const (
	NOT_ADMIN = false
	IS_ADMIN  = true
)

const (
	NOT_LOGINABLE = false
	IS_LOGINABLE  = true
)

const (
	ADMIN_SUBJECT    = "ADMIN"
	KANBOARD_SUBJECT = "KANBOARD"
)

const (
	TASK_STATUS_UNDO = iota
	TASK_STATUS_IN_PROGRESS
	TASK_STATUS_DONE
)

const (
	TASK_PRIORITY_LOW = iota - 1
	TASK_PRIORITY_MEDIUM
	TASK_PRIORITY_HIGH
)

const (
	KANBOARD_MESSAGE_CHANNEL = "KANBOARD_NOTIFICATION"
	ADMIN_MESSAGE_CHANNEL    = "ADMIN_NOTIFICATION"
)

const (
	NEW_MESSAGE     = "new_message"
	NEW_TASK_STATUS = "new_task_status"
	UNREAD_MESSAGE  = "unread_message"
	PUBLISH_MESSAGE = "publish_message"
	UPDATE_USER     = "update_user"
)

type EventType string

const (
	PROJECT_EVENT     EventType = "project_event"
	TASK_EVENT        EventType = "task_event"
	UPDATE_TASK_EVENT EventType = "update_task_event"
)
