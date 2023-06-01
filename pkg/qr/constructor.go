package qr

const (
	defaultLevel                   Correction = M
	defaultMinVer, defaultMaxVer   int        = 0, 40
	defaultMinMask, defaultMaxMask int        = 0, 8
)

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
