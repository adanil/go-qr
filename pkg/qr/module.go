package qr

// Module is a component of QR
type Module struct {
	value bool
	isSet bool
}

// Sets the Module color
func (m *Module) Set(value bool) {
	m.value = value
	m.isSet = true
}

func (m *Module) String() string {
	if !m.value {
		return "██"
	}
	return "  "
}
