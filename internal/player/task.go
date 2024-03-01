package player

type Task[T any] struct {
	val T
	Controller
}

func NewTask[T any](val T, streamer Controller) *Task[T] {
	return &Task[T]{
		Controller: streamer,
		val:        val,
	}
}

func (t *Task[T]) Value() T {
	return t.val
}
