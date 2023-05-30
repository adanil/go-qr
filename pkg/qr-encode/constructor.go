package qr_encode

type EncoderOptions func(*Encoder)

func WithVersionRange(minVersion, maxVersion int) EncoderOptions {
	return func(e *Encoder) {
		e.minVersion = minVersion
		e.maxVersion = maxVersion
	}
}

func NewEncoder(correctionLevel CodeLevel, options ...EncoderOptions) *Encoder {
	encoder := new(Encoder)
	encoder.level = correctionLevel
	encoder.minVersion = 0
	encoder.maxVersion = 40
	for _, option := range options {
		option(encoder)
	}
	return encoder
}
