package nidhi

type CollectionOptions struct {
	Fields []string
}

type CreateOptions struct {
	CreateMetadata  []MetadataMarshaler
	ReplaceMetadata []MetadataMarshaler
	// Replace will replace the document if it exists otherwise it will throw an error.
	Replace bool
}

type CreateOption func(*CreateOptions)

func WithCreateOptions(o CreateOptions) CreateOption {
	return func(opt *CreateOptions) {
		*opt = o
	}
}

func WithCreateCreateMetadata(mm ...MetadataMarshaler) CreateOption {
	return func(co *CreateOptions) {
		co.CreateMetadata = append(co.CreateMetadata, mm...)
	}
}

func WithCreateReplaceMetadata(mm ...MetadataMarshaler) CreateOption {
	return func(co *CreateOptions) {
		co.ReplaceMetadata = append(co.ReplaceMetadata, mm...)
	}
}

type DeleteOptions struct {
	Metadata  []MetadataMarshaler
	Permanent bool
}

type DeleteOption func(*DeleteOptions)

func WithDeleteOptions(o DeleteOptions) DeleteOption {
	return func(opt *DeleteOptions) {
		*opt = o
	}
}

func WithDeleteMetadata(mm ...MetadataMarshaler) DeleteOption {
	return func(do *DeleteOptions) {
		do.Metadata = append(do.Metadata, mm...)
	}
}

type QueryOptions struct {
	PaginationOptions *PaginationOptions
	ViewMask          []string
	CreateMetadata    []CreateMetadataFunc
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

func WithQueryCreateMetadata(mcf ...CreateMetadataFunc) QueryOption {
	return func(qo *QueryOptions) {
		qo.CreateMetadata = append(qo.CreateMetadata, mcf...)
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

	Metadata []MetadataUnmarshaler
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

func WithGetMetadata(mm ...MetadataUnmarshaler) GetOption {
	return func(opt *GetOptions) {
		opt.Metadata = append(opt.Metadata, mm...)
	}
}

type ReplaceOptions struct {
	Metadata []MetadataMarshaler
	Revision int64
}

type ReplaceOption func(*ReplaceOptions)

func WithReplaceOptions(o ReplaceOptions) ReplaceOption {
	return func(opt *ReplaceOptions) {
		*opt = o
	}
}

func WithReplaceMetadata(mm ...MetadataMarshaler) ReplaceOption {
	return func(ro *ReplaceOptions) {
		ro.Metadata = append(ro.Metadata, mm...)
	}
}

type UpdateOptions struct {
	Metadata []MetadataMarshaler
}

type UpdateOption func(*UpdateOptions)

func WithUpdateOptions(o UpdateOptions) UpdateOption {
	return func(uo *UpdateOptions) {
		*uo = o
	}
}

func WithUpdateMetadata(mm ...MetadataMarshaler) UpdateOption {
	return func(uo *UpdateOptions) {
		uo.Metadata = append(uo.Metadata, mm...)
	}
}
