package main

type Queue struct {
}

func (q *Queue) Emit(message interface{}) error {
	return nil
}

func (q *Queue) Peek(topic string) (int64, []byte, error) {
	return 0, nil, nil
}

func (q *Queue) Delete(messageId int64) error {
	return nil
}
