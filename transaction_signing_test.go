package overflow

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionSigning(t *testing.T) {
	o, err := OverflowTesting()
	assert.NoError(t, err)

	t.Run("one signer pay/propose/authorize", func(t *testing.T) {

		//The Tx dsl can also contain an inline transaction
		result := o.Tx(`
		import Debug from "../contracts/Debug.cdc"
		transaction {
		  prepare(account1: AuthAccount) {
				Debug.log("account1:".concat(account1.address.toString())) 
			} 
		}`, SignAs("first"))

		account, err := o.AccountE("first")
		assert.NoError(t, err)
		result.AssertSuccess(t).AssertDebugLog(t, fmt.Sprintf("account1:0x%s", account.Address().Hex()))
	})
	t.Run("one signer seperate proposer", func(t *testing.T) {

		//The Tx dsl can also contain an inline transaction
		result := o.Tx(`
		import Debug from "../contracts/Debug.cdc"
		transaction {
		  prepare(account1: AuthAccount) {
				Debug.log("account1:".concat(account1.address.toString())) 
			} 
		}`, ProposeAs("second"), SignAs("first"))

		account, err := o.AccountE("first")
		assert.NoError(t, err)
		result.AssertSuccess(t).AssertDebugLog(t, fmt.Sprintf("account1:0x%s", account.Address().Hex()))
	})

	t.Run("two signers, first pay/propose/authorize", func(t *testing.T) {

		//the payer is always the last auth account
		//The Tx dsl can also contain an inline transaction
		result := o.Tx(`
		import Debug from "../contracts/Debug.cdc"
		transaction {
			prepare(account1: AuthAccount, account2:AuthAccount) {
				Debug.log("account1:".concat(account1.address.toString())) 
				Debug.log("account2:".concat(account2.address.toString())) 
			} 
		}`, SignAs("first"), SignAs("second"))

		account1, err := o.AccountE("first")
		assert.NoError(t, err)

		//the payer is always last
		account2, err := o.AccountE("second")
		assert.NoError(t, err)
		result.AssertSuccess(t).AssertDebugLog(t, fmt.Sprintf("account1:0x%s", account2.Address().Hex()))
		result.AssertSuccess(t).AssertDebugLog(t, fmt.Sprintf("account2:0x%s", account1.Address().Hex()))
	})

	/*
		t.Run("one signer pay/propose/authorize", func(t *testing.T) {

			//The Tx dsl can also contain an inline transaction
			o.TxFileNameFN(`
			import Debug from "../contracts/Debug.cdc"
			transaction( {
			  prepare(account1: AuthAccount, account2: AuthAccount) {
					Debug.log("account1:".concat(account1.toString()))
					Debug.log("account2:".concat(account2.toString()))
				}
			}`, SignAs("first"))
		})
	*/

}
