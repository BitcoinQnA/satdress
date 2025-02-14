package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/fiatjaf/go-lnurl"
	"github.com/fiatjaf/makeinvoice"
)

func metaData(params *Params) lnurl.Metadata {

	//addImageToMetaData(w.telegram, &metadata, username, user.Telegram)
	return lnurl.Metadata{
		Description:      fmt.Sprintf("Pay to %s@%s", params.Name, params.Domain),
		LightningAddress: fmt.Sprintf("%s@%s", params.Name, params.Domain),
	}
}

func makeInvoice(
	params *Params,
	msat int,
	pin *string,
) (bolt11 string, err error) {
	// prepare params
	var backend makeinvoice.BackendParams
	switch params.Kind {
	case "sparko":
		backend = makeinvoice.SparkoParams{
			Host: params.Host,
			Key:  params.Key,
		}
	case "lnd":
		backend = makeinvoice.LNDParams{
			Host:     params.Host,
			Macaroon: params.Key,
		}
	case "lnbits":
		backend = makeinvoice.LNBitsParams{
			Host: params.Host,
			Key:  params.Key,
		}
	case "lnpay":
		backend = makeinvoice.LNPayParams{
			PublicAccessKey:  params.Pak,
			WalletInvoiceKey: params.Waki,
		}
	case "eclair":
		backend = makeinvoice.EclairParams{
			Host:     params.Host,
			Password: "",
		}
	case "commando":
		backend = makeinvoice.CommandoParams{
			Host:   params.Host,
			NodeId: params.NodeId,
			Rune:   params.Rune,
		}
	}

	mip := makeinvoice.Params{
		Msatoshi: int64(msat),
		Backend:  backend,

		Label: params.Domain + "/" + strconv.FormatInt(time.Now().Unix(), 16),
	}

	if pin != nil {
		// use this as the description for new accounts
		mip.Description = fmt.Sprintf("%s's PIN for '%s@%s' lightning address: %s", params.Domain, params.Name, params.Domain, *pin)
	} else {
		//use zapEventSerializedStr if nip57
		if params.zapEventSerializedStr != "" {
			mip.Description = params.zapEventSerializedStr
		} else {
			// make the lnurlpay description_hash
			mip.Description = metaData(params).Encode()
		}
		mip.UseDescriptionHash = true
	}

	// actually generate the invoice
	bolt11, err = makeinvoice.MakeInvoice(mip)

	log.Debug().Int("msatoshi", msat).
		Interface("backend", backend).
		Str("bolt11", bolt11).Err(err).Str("Description", mip.Description).
		Msg("invoice generation")

	return bolt11, err
}
