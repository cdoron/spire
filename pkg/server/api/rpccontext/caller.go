package rpccontext

import (
	"context"

	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/spire/proto/spire-next/types"
	"github.com/spiffe/spire/proto/spire/common"
)

type callerIDKey struct{}
type callerEntriesKey struct{}
type attestedNodeKey struct{}

func WithCallerID(ctx context.Context, id spiffeid.ID) context.Context {
	return context.WithValue(ctx, callerIDKey{}, id)
}

func CallerID(ctx context.Context) (spiffeid.ID, bool) {
	id, ok := ctx.Value(callerIDKey{}).(spiffeid.ID)
	return id, ok
}

func WithCallerEntry(ctx context.Context, entry *types.Entry) context.Context {
	return context.WithValue(ctx, callerEntriesKey{}, append(callerEntries(ctx), entry))
}

func CallerDownstreamEntries(ctx context.Context) []*types.Entry {
	var out []*types.Entry
	for _, entry := range callerEntries(ctx) {
		if entry.Downstream {
			out = append(out, entry)
		}
	}
	return out
}

func CallerAdminEntries(ctx context.Context) []*types.Entry {
	var out []*types.Entry
	for _, entry := range callerEntries(ctx) {
		if entry.Admin {
			out = append(out, entry)
		}
	}
	return out
}

func CallerIsDownstream(ctx context.Context) bool {
	for _, entry := range callerEntries(ctx) {
		if entry.Downstream {
			return true
		}
	}
	return false
}

func CallerIsAdmin(ctx context.Context) bool {
	for _, entry := range callerEntries(ctx) {
		if entry.Admin {
			return true
		}
	}
	return false
}

func CallerIsAgent(ctx context.Context) bool {
	_, ok := ctx.Value(attestedNodeKey{}).(*common.AttestedNode)
	return ok
}

func WithAttestedNode(ctx context.Context, attestedNode *common.AttestedNode) context.Context {
	return context.WithValue(ctx, attestedNodeKey{}, attestedNode)
}

func AttestedNode(ctx context.Context) (*common.AttestedNode, bool) {
	attestedNode, ok := ctx.Value(attestedNodeKey{}).(*common.AttestedNode)
	return attestedNode, ok
}

func callerEntries(ctx context.Context) []*types.Entry {
	entries, _ := ctx.Value(callerEntriesKey{}).([]*types.Entry)
	return entries
}
