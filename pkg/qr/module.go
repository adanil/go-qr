package qr

type qrModule struct {
	value bool
	isSet bool
}

func (m *qrModule) Set(value bool) {
	m.value = value
	m.isSet = true
}

func (m *qrModule) String() string {
	if !m.value {
		return "██"
	}
	return "  "
}
