package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/kwunyeung/desmos/x/magpie/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	coinKeeper bank.Keeper
	storeKey   sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc        *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the magpie Keeper
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

func (k Keeper) SetPost(ctx sdk.Context, post types.Post) sdk.Error {
	if post.Owner.Empty() {
		return sdk.ErrInvalidAddress("No address found.")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(post.ID), k.cdc.MustMarshalBinaryBare(post))

	return nil
}

func (k Keeper) GetPost(ctx sdk.Context, id string) types.Post {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(id)) {
		return types.NewPost()
	}

	bz := store.Get([]byte(id))
	var post types.Post
	k.cdc.MustUnmarshalBinaryBare(bz, &post)
	return post
}

func (k Keeper) EditPost(ctx sdk.Context, id string, message string) sdk.Error {
	post := k.GetPost(ctx, id)
	post.Message = message

	if err := k.SetPost(ctx, post); err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetPostOwner(ctx sdk.Context, id string) sdk.AccAddress {
	return k.GetPost(ctx, id).Owner
}

func (k Keeper) GetPostLikes(ctx sdk.Context, id string) uint {
	return k.GetPost(ctx, id).Likes
}

func (k Keeper) AddPostLike(ctx sdk.Context, id string) sdk.Error {
	post := k.GetPost(ctx, id)
	post.Likes = post.Likes + 1
	return k.SetPost(ctx, post)
}

func (k Keeper) SetLike(ctx sdk.Context, id string, like types.Like) sdk.Error {
	if like.Owner.Empty() || (len(like.PostID) == 0) {
		return sdk.ErrUnauthorized("Liker and post id must exist.")
	}

	post := k.GetPost(ctx, like.PostID)
	if len(post.ID) == 0 {
		return sdk.ErrUnknownRequest("The post requested does not exist.")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(id), k.cdc.MustMarshalBinaryBare(like))

	return k.AddPostLike(ctx, like.PostID)
}

func (k Keeper) GetLike(ctx sdk.Context, id string) types.Like {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(id)) {
		return types.NewLike()
	}

	bz := store.Get([]byte(id))
	var like types.Like
	k.cdc.MustUnmarshalBinaryBare(bz, &like)
	return like
}

func (k Keeper) GetPostsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte{})
}

func (k Keeper) SetSession(ctx sdk.Context, session types.Session) sdk.Error {
	if session.Owner.Empty() {
		return sdk.ErrInvalidAddress("No address found.")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(session.ID), k.cdc.MustMarshalBinaryBare(session))

	return nil
}

func (k Keeper) GetSession(ctx sdk.Context, id string) types.Session {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(id)) {
		return types.NewSession()
	}

	bz := store.Get([]byte(id))
	var session types.Session
	k.cdc.MustUnmarshalBinaryBare(bz, &session)
	return session
}

func (k Keeper) EditSession(ctx sdk.Context, id string, expiry time.Time) sdk.Error {
	session := k.GetSession(ctx, id)
	session.Expiry = expiry

	return k.SetSession(ctx, session)
}
