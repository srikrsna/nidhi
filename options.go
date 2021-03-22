package nidhi

type CollectionOptions struct {
	Fields []string
}

type CreateOptions struct {
	CreateMetadata  Metadata
	ReplaceMetadata Metadata
	// Replace will replace the document if it exists otherwise it will throw an error.
	Replace bool
}

type CreateOption func(*CreateOptions)

func WithCreateOptions(o CreateOptions) CreateOption {
	return func(opt *CreateOptions) {
		*opt = o
	}
}

func WithCreateMetadataCreateOptions(mm ...MetadataMarshaler) CreateOption {
	return func(co *CreateOptions) {
		co.CreateMetadata = append(co.CreateMetadata, mm...)
	}
}

func WithReplaceMetadataReplaceOptions(mm ...MetadataMarshaler) CreateOption {
	return func(co *CreateOptions) {
		co.ReplaceMetadata = append(co.ReplaceMetadata, mm...)
	}
}

type DeleteOptions struct {
	Metadata  Metadata
	Permanent bool
}

type DeleteOption func(*DeleteOptions)

func WithDeleteOptions(o DeleteOptions) DeleteOption {
	return func(opt *DeleteOptions) {
		*opt = o
	}
}

func WithMetadataDeleteOptions(mm ...MetadataMarshaler) DeleteOption {
	return func(do *DeleteOptions) {
		do.Metadata = append(do.Metadata, mm...)
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
	Cursor string
	// NextCursor will be set by Nidhi
	NextCursor string
	// OrderBy fields, if empty defaults to [{"id", false}]
	OrderBy  []OrderBy
	Limit    uint64
	Backward bool
	// HasMore will be set by Nidhi
	HasMore bool
}

type OrderBy struct {
	Field Orderer
	Desc  bool
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

type ReplaceOptions struct {
	Metadata Metadata
	Revision int64
}

type ReplaceOption func(*ReplaceOptions)

func WithReplaceOptions(o ReplaceOptions) ReplaceOption {
	return func(opt *ReplaceOptions) {
		*opt = o
	}
}

func WithMetadataReplaceOptions(mm ...MetadataMarshaler) ReplaceOption {
	return func(ro *ReplaceOptions) {
		ro.Metadata = append(ro.Metadata, mm...)
	}
}

type UpdateOptions struct {
	Metadata Metadata
}

type UpdateOption func(*UpdateOptions)

func WithUpdateOptions(o UpdateOptions) UpdateOption {
	return func(uo *UpdateOptions) {
		*uo = o
	}
}

func WithMetadataUpdateOptions(mm ...MetadataMarshaler) UpdateOption {
	return func(uo *UpdateOptions) {
		uo.Metadata = append(uo.Metadata, mm...)
	}
}
