package qr_encode

type (
	Option interface {
		isOptional()
	}

	Versions struct {
		minVersion int
		maxVersion int
	}
	//Logger struct {
	//	logger zap.Logger
	//}
)

func (Versions) isOptional() {}

func WithVersionRange(minVersion, maxVersion int) Option {
	return Versions{
		minVersion: minVersion,
		maxVersion: maxVersion,
	}
}

func NewEncoder(correctionLevel CodeLevel, options ...Option) *Encoder {
	encoder := new(Encoder)
	encoder.level = correctionLevel
	encoder.minVersion = 0
	encoder.maxVersion = 40
	for _, option := range options {
		switch op := option.(type) {
		case Versions:
			encoder.minVersion = op.minVersion
			encoder.maxVersion = op.maxVersion
		default:
			panic("wrong type of option")
		}
	}
	return encoder
}
