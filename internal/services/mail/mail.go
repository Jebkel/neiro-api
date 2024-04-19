package mail

import (
	"neiro-api/pkg/queue"
)

type ServiceMail struct {
	data Data
}

type Data struct {
	To      string
	Subject string
	Header  string
	From    string
	Lines   []string
}

func NewMailService() *ServiceMail {
	return &ServiceMail{}
}

func (m *ServiceMail) New() *ServiceMail {
	return NewMailService()
}

func (m *ServiceMail) To(recipient string) *ServiceMail {
	m.data.To = recipient
	return m
}

func (m *ServiceMail) From(from string) *ServiceMail {
	m.data.From = from
	return m
}

func (m *ServiceMail) Line(line string) *ServiceMail {
	m.data.Lines = append(m.data.Lines, line)
	return m
}

func (m *ServiceMail) Header(header string) *ServiceMail {
	m.data.Header = header
	return m
}

func (m *ServiceMail) Subject(subject string) *ServiceMail {
	m.data.Subject = subject
	return m
}

func (m *ServiceMail) Send() {
	queueService := queue.GetQueue()
	queueService.AddTask("mail", m.data)
}
