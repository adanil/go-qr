package qr

const (
	defaultLevel                   Correction = M
	defaultMinVer, defaultMaxVer   int        = 0, 40
	defaultMinMask, defaultMaxMask int        = 0, 8
)

// EncoderOptions is a functional object that can be provided to NewEncoder to specify parameters of Encoder
type EncoderOptions func(*Encoder)

// WithCorrectionLevel is a Encoder option that allows to specify a level of correction to be used in QR
func WithCorrectionLevel(level Correction) EncoderOptions {
	return func(e *Encoder) {
		e.level = level
	}
}

// WithVersionRange is a Encoder option that allows to specify a range of versions can be used to encode data
func WithVersionRange(minVersion, maxVersion int) EncoderOptions {
	return func(e *Encoder) {
		e.minVersion = minVersion
		e.maxVersion = maxVersion
	}
}

// WithMaskRange is a Encoder option that allows to specify a range of mask can be used to optimize mask
func WithMaskRange(minMask, maxMask int) EncoderOptions {
	return func(e *Encoder) {
		e.minMask = minMask
		e.maxMask = maxMask
	}
}

// NewEncoder returns a new Encoder with default options if none are provided
func NewEncoder(options ...EncoderOptions) *Encoder {
	encoder := &Encoder{
		level:      defaultLevel,
		minVersion: defaultMinVer,
		maxVersion: defaultMaxVer,
		minMask:    defaultMinMask,
		maxMask:    defaultMaxMask,
	}

	for _, option := range options {
		option(encoder)
	}

	return encoder
}
