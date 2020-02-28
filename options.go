package nidhi

type CreateOptions struct {
	Metadata                Metadata
	ReplaceIfExists         bool
	ReplaceMetadataIfExists bool
}

type CreateOption func(*CreateOptions)

func WithCreateOptions(o CreateOptions) CreateOption {
	return func(opt *CreateOptions) {
		*opt = o
	}
}

type DeleteOptions struct {
	Permanent bool

	Metadata        Metadata
	ReplaceMetadata bool
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

type GetOptions struct {
	LoadMetadata Metadata
}

type GetOption func(*GetOptions)

func WithGetOptions(o GetOptions) GetOption {
	return func(opt *GetOptions) {
		*opt = o
	}
}

type CountOptions struct {
}

type CountOption func(*CountOptions)

func WithCountOptions(o CountOptions) CountOption {
	return func(opt *CountOptions) {
		*opt = o
	}
}

type ReplaceOptions struct {
	Revision        int64
	Metadata        Metadata
	ReplaceMetadata bool
}

type ReplaceOption func(*ReplaceOptions)

func WithReplaceOptions(o ReplaceOptions) ReplaceOption {
	return func(opt *ReplaceOptions) {
		*opt = o
	}
}
