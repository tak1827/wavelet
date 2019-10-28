package api

import (
	"encoding/hex"

	"github.com/perlin-network/wavelet"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
)

type Account struct {
	ID         wavelet.AccountID `json:"id"`
	Balance    uint64            `json:"balance"`
	GasBalance uint64            `json:"gas_balance"`
	Stake      uint64            `json:"stake"`
	Reward     uint64            `json:"reward"`
	Nonce      uint64            `json:"nonce"`
	IsContract bool              `json:"is_contract"`
	NumPages   uint64            `json:"num_pages"`
}

var _ JSONObject = (*Account)(nil)

func (g *Gateway) getAccount(ctx *fasthttp.RequestCtx) {
	param, ok := ctx.UserValue("id").(string)
	if !ok {
		g.renderError(ctx, ErrBadRequest(errors.New("id must be a string")))
		return
	}

	slice, err := hex.DecodeString(param)
	if err != nil {
		g.renderError(ctx, ErrBadRequest(errors.Wrap(
			err, "account ID must be presented as valid hex")))
		return
	}

	if len(slice) != wavelet.SizeAccountID {
		g.renderError(ctx, ErrBadRequest(errors.Errorf(
			"account ID must be %d bytes long", wavelet.SizeAccountID)))
		return
	}

	var id wavelet.AccountID
	copy(id[:], slice)

	snapshot := g.Ledger.Snapshot()

	balance, _ := wavelet.ReadAccountBalance(snapshot, id)
	gasBalance, _ := wavelet.ReadAccountContractGasBalance(snapshot, id)
	stake, _ := wavelet.ReadAccountStake(snapshot, id)
	reward, _ := wavelet.ReadAccountReward(snapshot, id)
	nonce, _ := wavelet.ReadAccountNonce(snapshot, id)
	_, isContract := wavelet.ReadAccountContractCode(snapshot, id)
	numPages, _ := wavelet.ReadAccountContractNumPages(snapshot, id)

	g.render(ctx, &Account{
		ID:         id,
		Balance:    balance,
		GasBalance: gasBalance,
		Stake:      stake,
		Reward:     reward,
		Nonce:      nonce,
		IsContract: isContract,
		NumPages:   numPages,
	})
}

func (s *Account) MarshalArena(arena *fastjson.Arena) ([]byte, error) {
	o := arena.NewObject()

	arenaSet(arena, o, "id", s.ID)
	arenaSet(arena, o, "balance", s.Balance)
	arenaSet(arena, o, "gas_balance", s.GasBalance)
	arenaSet(arena, o, "stake", s.Stake)
	arenaSet(arena, o, "reward", s.Reward)
	arenaSet(arena, o, "nonce", s.Nonce)
	arenaSet(arena, o, "is_contract", s.IsContract)

	if s.NumPages != 0 {
		arenaSet(arena, o, "num_pages", s.NumPages)
	}

	return o.MarshalTo(nil), nil
}

func (s *Account) UnmarshalValue(v *fastjson.Value) error {
	if err := valueHex(v, s.ID[:], "id"); err != nil {
		return err
	}

	s.Balance = v.GetUint64("balance")
	s.GasBalance = v.GetUint64("gas_balance")
	s.Stake = v.GetUint64("stake")
	s.Reward = v.GetUint64("reward")
	s.Nonce = v.GetUint64("nonce")
	s.IsContract = v.GetBool("is_contract")
	s.NumPages = v.GetUint64("num_pages")
	return nil
}