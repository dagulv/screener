package port

import "context"

// Used for acquiring a context from a port. If it's e.g. a Postgres store that is implementing the port,
// the context could be a transaction, ensuring that any queries with the same context are executed in
// one single, atomic transaction.
//
// Using a Port's contexts is optional, but if you once run `AcquireContext()` you MUST release it once you're
// done. Otherwise you might cause a deadlock if it's a transaction that never is either commited or
// rolled back.
//
// It's always safe to run `CommitContext()` and/or `ReleaseContext()` on a context that is already
// committed or released.
//
// The recommended approach is:
//
//		 if ctx, err = port.AcquireContext(ctx); err != nil {
//		     return
//		 }
//
//	     // Always release when finished, no matter what
//		 defer port.ReleaseContext(ctx)
//
//		 // Do some work on ctx and/or abort if you encounter an error
//
//		 // Finish by committing
//		 if err = port.CommitContext(ctx); err != nil {
//		     return
//		 }
type Context interface {
	AcquireContext(ctx context.Context, readOnly ...bool) (newCtx context.Context, err error)
	CommitContext(ctx context.Context) error
	ReleaseContext(ctx context.Context) error
}
