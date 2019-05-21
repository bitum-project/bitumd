// Copyright (c) 2019 The Bitum developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package updater

import (
	"github.com/bitum-project/bitumd/wire"
)

// SplitMsgInv splits an inventory message into non-updater and updater
// inventory messages.
func SplitMsgInv(msg *wire.MsgInv) (*wire.MsgInv, *wire.MsgInv, error) {
	nonUpdaterMsg := wire.NewMsgInvSizeHint(uint(len(msg.InvList)))
	updaterMsg := wire.NewMsgInvSizeHint(uint(len(msg.InvList)))
	for _, invVect := range msg.InvList {
		if invVect.Type == wire.InvTypeCodechainEntry ||
			invVect.Type == wire.InvTypePatch {
			if err := updaterMsg.AddInvVect(invVect); err != nil {
				return nil, nil, err
			}
		} else {
			if err := nonUpdaterMsg.AddInvVect(invVect); err != nil {
				return nil, nil, err
			}
		}
	}
	return nonUpdaterMsg, updaterMsg, nil
}
