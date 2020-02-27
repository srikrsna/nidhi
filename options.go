package nidhi

type CreateOptions struct {
}

type CreateOption func(*CreateOptions)

func WithCreateOptions(o CreateOptions) CreateOption {
	return func(opt *CreateOptions) {
		*opt = o
	}
}

type DeleteOptions struct {
	Permanent bool
}

type DeleteOption func(*DeleteOptions)

func WithDeleteOptions(o DeleteOptions) DeleteOption {
	return func(opt *DeleteOptions) {
		*opt = o
	}
}

type QueryOptions struct {
}

type QueryOption func(*QueryOptions)

func WithQueryOptions(o QueryOptions) QueryOption {
	return func(opt *QueryOptions) {
		*opt = o
	}
}

type ReplaceOptions struct {
	Revision int64
}

type ReplaceOption func(*ReplaceOptions)

func WithReplaceOptions(o ReplaceOptions) ReplaceOption {
	return func(opt *ReplaceOptions) {
		*opt = o
	}
}
