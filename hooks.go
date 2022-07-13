package nidhi

import "context"

// Hooks are callbacks that are called on store operations.
//
// The passed values can if modified will change for the op.
type Hooks struct {
	// OnCreate is called on [Store.Create].
	OnCreate OnCreateHook
	// OnGet is called on [Store.Get].
	OnGet OnGetHook
	// OnQuery is called on [Store.Query].
	OnQuery OnQueryHook
	// OnDelete is called on [Store.Delete].
	OnDelete OnDeleteHook
	// OnDeleteMany is called on [Store.DeleteMany].
	OnDeleteMany OnDeleteManyHook
	// OnReplace is called on [Store.Replace].
	OnReplace OnReplaceHook
	// OnUpdate is called on [Store.Update].
	OnUpdate OnUpdateHook
	// OnUpdateMany is called on [Store.UpdateMany].
	OnUpdateMany OnUpdateManyHook
}

// HookContext is the context passed to [Hooks]
type HookContext struct {
	context.Context

	idFn    func(any) string
	setIdFn func(any, string)
}

// NewHookContext returns a [HookContext].
func NewHookContext[T any](ctx context.Context, store *Store[T]) *HookContext {
	return &HookContext{
		Context: ctx,
		idFn: func(t any) string {
			return store.idFn(t.(*T))
		},
		setIdFn: func(t any, s string) {
			store.setIdFn(t.(*T), s)
		},
	}
}

// SetId sets the id field of the document.
func (s *HookContext) SetId(v any, id string) {
	s.setIdFn(v, id)
}

// Id returns the id of the document
func (s *HookContext) Id(v any) string {
	return s.idFn(v)
}
