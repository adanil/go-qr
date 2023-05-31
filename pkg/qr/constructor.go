package qr

type EncoderOptions func(*Encoder)

func WithCorrectionLevel(level Correction) EncoderOptions {
	return func(e *Encoder) {
		e.level = level
	}
}

func WithVersionRange(minVersion, maxVersion int) EncoderOptions {
	return func(e *Encoder) {
		e.minVersion = minVersion
		e.maxVersion = maxVersion
	}
}

func WithMaskRange(minMask, maxMask int) EncoderOptions {
	return func(e *Encoder) {
		e.minMask = minMask
		e.maxMask = maxMask
	}
}

func NewEncoder(options ...EncoderOptions) *Encoder {
	encoder := &Encoder{
		level:      M,
		minVersion: 0, maxVersion: 40,
		minMask: 0, maxMask: 8,
	}

	for _, option := range options {
		option(encoder)
	}

	return encoder
}
