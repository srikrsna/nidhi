package nidhi

type CollectionOptions struct {
	SubjectFunc SubjectFunc

	Fields []string
}

type CreateOptions struct {
	// Replace will replace the document if it exists otherwise it will throw an error.
	Replace bool
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
	PaginationOptions *PaginationOptions
	ViewMask          []string
}

type QueryOption func(*QueryOptions)

func WithQueryOptions(o QueryOptions) QueryOption {
	return func(opt *QueryOptions) {
		*opt = o
	}
}

func WithQueryViewMask(vm []string) QueryOption {
	return func(qo *QueryOptions) {
		qo.ViewMask = vm
	}
}

func WithPaginationOptions(po *PaginationOptions) QueryOption {
	return func(qo *QueryOptions) {
		qo.PaginationOptions = po
	}
}

type PaginationOptions struct {
	Backward bool
	Cursor   string
	Limit    uint64

	// Will be set by Nidhi
	HasMore bool
}

type GetOptions struct {
	ViewMask []string
}

type GetOption func(*GetOptions)

func WithGetOptions(o GetOptions) GetOption {
	return func(opt *GetOptions) {
		*opt = o
	}
}

func WithGetViewMask(vm []string) GetOption {
	return func(opt *GetOptions) {
		opt.ViewMask = vm
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
	Revision int64
}

type ReplaceOption func(*ReplaceOptions)

func WithReplaceOptions(o ReplaceOptions) ReplaceOption {
	return func(opt *ReplaceOptions) {
		*opt = o
	}
}

type UpdateOptions struct {
}

type UpdateOption func(*UpdateOptions)

func WithUpdateOptions(o UpdateOptions) UpdateOption {
	return func(uo *UpdateOptions) {
		*uo = o
	}
}
