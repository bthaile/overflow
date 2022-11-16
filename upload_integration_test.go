package overflow

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionUpload(t *testing.T) {
	g := NewTestingEmulator().Start()

	t.Run("Upload image file invalid file", func(t *testing.T) {

		err := g.UploadImageAsDataUrl("testFile2.txt", "first")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "could not read imageFile testFile2.txt")
	})

	t.Run("Upload test file invalid file", func(t *testing.T) {

		err := g.UploadFile("testFile2.txt", "first")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "could not read file testFile2.txt")
	})

	t.Run("Upload test file", func(t *testing.T) {

		//need flow in account to upload
		g.TransactionFromFile("mint_tokens").
			SignProposeAndPayAsService().
			Args(g.Arguments().
				Account("first").
				UFix64(100.0)).
			Test(t).
			AssertSuccess()

		err := g.UploadFile("testdata/testFile.txt", "first")
		assert.NoError(t, err)
		g.Transaction(`
import Debug from "../contracts/Debug.cdc"
transaction {
  prepare(account: AuthAccount) {
    var content= account.load<String>(from: /storage/upload) ?? panic("could not load content")
	Debug.log(content)
 }
}`).
			SignProposeAndPayAs("first").
			Test(t).
			AssertSuccess().
			AssertDebugLog("VGhpcyBpcyBhIGZpbGU=")
	})

	t.Run("Upload test image", func(t *testing.T) {

		//need flow in account to upload
		g.TransactionFromFile("mint_tokens").
			SignProposeAndPayAsService().
			Args(g.Arguments().
				Account("first").
				UFix64(100.0)).
			Test(t).
			AssertSuccess()

		err := g.UploadImageAsDataUrl("testdata/pig.png", "first")
		assert.NoError(t, err)
		g.Transaction(`
import Debug from "../contracts/Debug.cdc"
transaction {
  prepare(account: AuthAccount) {
    var content= account.load<String>(from: /storage/upload) ?? panic("could not load content")
	Debug.log(content)
 }
}`).
			SignProposeAndPayAs("first").
			Test(t).
			AssertSuccess().
			AssertDebugLog("data:image/png;base64, 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")
	})

	t.Run("Download and upload file fails wrong url", func(t *testing.T) {

		err := g.DownloadAndUploadFile("https://foo.bar", "first")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "dial tcp: lookup foo.ba")

	})

	t.Run("Download and upload imagefile fails wrong url", func(t *testing.T) {

		err := g.DownloadImageAndUploadAsDataUrl("https://foo.bar", "first")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "dial tcp: lookup foo.ba")

	})

	t.Run("Download and upload file", func(t *testing.T) {
		//need flow in account to upload image
		g.TransactionFromFile("mint_tokens").
			SignProposeAndPayAsService().
			Args(g.Arguments().
				Account("first").
				UFix64(100.0)).
			Test(t).
			AssertSuccess()

		err := g.DownloadAndUploadFile("https://httpbin.org/image/png", "first")
		assert.NoError(t, err)

		g.Transaction(`
import Debug from "../contracts/Debug.cdc"
transaction {
  prepare(account: AuthAccount) {
    var content= account.load<String>(from: /storage/upload) ?? panic("could not load content")
	Debug.log(content)
 }
}`).
			SignProposeAndPayAs("first").
			Test(t).
			AssertSuccess()

	})
	t.Run("Download image and upload", func(t *testing.T) {
		//need flow in account to upload image
		g.TransactionFromFile("mint_tokens").
			SignProposeAndPayAsService().
			Args(g.Arguments().
				Account("first").
				UFix64(100.0)).
			Test(t).
			AssertSuccess()

		err := g.DownloadImageAndUploadAsDataUrl("https://httpbin.org/image/png", "first")
		assert.NoError(t, err)

		g.Transaction(`
import Debug from "../contracts/Debug.cdc"
transaction {
  prepare(account: AuthAccount) {
    var content= account.load<String>(from: /storage/upload) ?? panic("could not load content")
	Debug.log(content)
 }
}`).
			SignProposeAndPayAs("first").
			Test(t).
			AssertSuccess()

	})
}
