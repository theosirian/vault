// +build integration vault

package test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	provide "github.com/provideplatform/provide-go/api/vault"
	"github.com/provideplatform/vault/common"
	cryptovault "github.com/provideplatform/vault/vault"
)

func keyFactoryWithSeed(token, vaultID, keyType, keyUsage, keySpec, keyName, keyDescription, seedPhrase string) (*provide.Key, error) {

	resp, err := provide.CreateKey(token, vaultID, map[string]interface{}{
		"type":        keyType,
		"usage":       keyUsage,
		"spec":        keySpec,
		"name":        keyName,
		"description": keyDescription,
		"mnemonic":    seedPhrase,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create key error: %s", err.Error())
	}

	key := &provide.Key{}
	respRaw, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to marshall key data: %s", err.Error())
	}
	json.Unmarshal(respRaw, &key)
	return key, nil
}

func init() {
	token, err := userTokenFactory()
	if err != nil {
		log.Printf("failed to create token; %s", err.Error())
		return
	}

	time.Sleep(time.Second * 5)

	//test getting a new unsealer key
	newkeyresp, err := provide.GenerateSeal(*token, map[string]interface{}{})
	if err != nil {
		log.Printf("error generating new unsealer key %s", err.Error())
	}
	log.Printf("newkeyresp: %+v", *newkeyresp)
	log.Printf("newly generated unsealer key %s", *newkeyresp.UnsealerKey)
	log.Printf("newly generated unsealer key hash %s", *newkeyresp.ValidationHash)

	_, err = provide.Unseal(token, map[string]interface{}{
		"key": "traffic charge swing glimpse will citizen push mutual embrace volcano siege identify gossip battle casual exit enrich unlock muscle vast female initial please day",
	})
	if err != nil {
		log.Printf("**vault not unsealed**. error: %s", err.Error())
		return
	}

	// now try it again, and we expect a 204 (no response) when trying to unseal a sealed vault
	_, err = provide.Unseal(token, map[string]interface{}{
		"key": "traffic charge swing glimpse will citizen push mutual embrace volcano siege identify gossip battle casual exit enrich unlock muscle vast female initial please day",
	})
	if err != nil {
		log.Printf("**second unseal attempt failed when it should pass**. error: %s", err.Error())
		return
	}

}

func TestSealUnsealer(t *testing.T) {
	// we're going to unseal the vault,
	// do an operation
	// seal the vault
	// retry operation, expect it to fail
	// unseal the vault, retry operation and expect it to succeed

	//get the vault unsealed to make sure other tests can continue
	defer unsealVault()

	token, err := userTokenFactory()
	if err != nil {
		log.Printf("failed to create token; %s", err.Error())
		return
	}

	_, err = provide.Unseal(token, map[string]interface{}{
		"key": "traffic charge swing glimpse will citizen push mutual embrace volcano siege identify gossip battle casual exit enrich unlock muscle vast female initial please day",
	})
	if err != nil {
		t.Errorf("**vault not unsealed**. error: %s", err.Error())
		return
	}

	_, err = vaultFactory(*token, "vaulty vault", "just a boring vaulty vault")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	_, err = provide.Seal(*token, map[string]interface{}{
		"key": "traffic charge swing glimpse will citizen push mutual embrace volcano siege identify gossip battle casual exit enrich unlock muscle vast female initial please day",
	})
	if err != nil {
		t.Errorf("**vault not sealed**. error: %s", err.Error())
		return
	}

	_, err = vaultFactory(*token, "vaulty vault", "just a boring vaulty vault")
	if err == nil {
		t.Errorf("performed operation while sealed!")
		return
	}

	_, err = provide.Unseal(token, map[string]interface{}{
		"key": "traffic charge swing glimpse will citizen push mutual embrace volcano siege identify gossip battle casual exit enrich unlock muscle vast female initial please day",
	})
	if err != nil {
		t.Errorf("**vault not unsealed**. error: %s", err.Error())
		return
	}

	_, err = vaultFactory(*token, "vaulty vault", "just a boring vaulty vault")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	// now we'll try to seal it badly and expect it to continue working
	_, err = provide.Seal(*token, map[string]interface{}{
		"key": "raffic charge swing glimpse will citizen push mutual embrace volcano siege identify gossip battle casual exit enrich unlock muscle vast female initial please day",
	})
	if err == nil {
		t.Errorf("**vault sealed with bad key**")
		return
	}

	_, err = vaultFactory(*token, "vaulty vault", "just a boring vaulty vault")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	// now we'll seal it and unseal it badly and expect it to fail
	_, err = provide.Seal(*token, map[string]interface{}{
		"key": "traffic charge swing glimpse will citizen push mutual embrace volcano siege identify gossip battle casual exit enrich unlock muscle vast female initial please day",
	})
	if err != nil {
		t.Errorf("**vault not sealed**. error: %s", err.Error())
		return
	}

	_, err = vaultFactory(*token, "vaulty vault", "just a boring vaulty vault")
	if err == nil {
		t.Errorf("performed operation while sealed!")
		return
	}

	_, err = provide.Unseal(token, map[string]interface{}{
		"key": "raffic charge swing glimpse will citizen push mutual embrace volcano siege identify gossip battle casual exit enrich unlock muscle vast female initial please day",
	})
	if err == nil {
		t.Errorf("unsealed vault with bad key.")
		return
	}

	_, err = vaultFactory(*token, "vaulty vault", "just a boring vaulty vault")
	if err == nil {
		t.Errorf("created vault while sealed!")
		return
	}

	//finish up with a valid unseal, before the additional deferred unseal
	_, err = provide.Unseal(token, map[string]interface{}{
		"key": "traffic charge swing glimpse will citizen push mutual embrace volcano siege identify gossip battle casual exit enrich unlock muscle vast female initial please day",
	})
	if err != nil {
		t.Errorf("**vault not unsealed**. error: %s", err.Error())
		return
	}

}

func TestAPICreateVault(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	_, err = vaultFactory(*token, "vaulty vault", "just a boring vaulty vault")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

}

func TestAPICreateKey(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	_, err = provide.CreateKey(*token, vault.ID.String(), map[string]interface{}{
		"type":        "asymmetric",
		"usage":       "sign/verify",
		"spec":        "secp256k1",
		"name":        "integration test ethereum key",
		"description": "organization eth/stablecoin wallet",
	})

	if err != nil {
		t.Errorf("failed to create key error: %s", err.Error())
		return
	}
}

func TestAPIDeleteKey(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	key, err := keyFactory(*token, vault.ID.String(), "asymmetric", "sign/verify", "secp256K1", "namey name", "cute description")
	if err != nil {
		t.Errorf("failed to create key; %s", err.Error())
		return
	}

	err = provide.DeleteKey(*token, vault.ID.String(), key.ID.String())
	if err != nil {
		t.Errorf("failed to delete key for vault: %s", err.Error())
		return
	}
}

func TestAPISign(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	key, err := keyFactory(*token, vault.ID.String(), "asymmetric", "sign/verify", "SECP256K1", "namey name", "cute description")
	if err != nil {
		t.Errorf("failed to create key; %s", err.Error())
		return
	}

	payloadBytes, _ := common.RandomBytes(32)
	payload := hex.EncodeToString(payloadBytes)
	_, err = provide.SignMessage(*token, vault.ID.String(), key.ID.String(), payload, nil)
	if err != nil {
		t.Errorf("failed to sign message %s", err.Error())
		return
	}

	// TODO check for signature in response, not sure if the errors are actually tripping
}

func TestAPIVerifySecp256k1Signature(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	key, err := keyFactory(*token, vault.ID.String(), "asymmetric", "sign/verify", cryptovault.KeySpecECCSecp256k1, "namey name", "cute description")
	if err != nil {
		t.Errorf("failed to create key; %s", err.Error())
		return
	}

	payloadBytes, _ := common.RandomBytes(32)
	messageToSign := hex.EncodeToString(payloadBytes)
	sigresponse, err := provide.SignMessage(*token, vault.ID.String(), key.ID.String(), messageToSign, nil)
	if err != nil {
		t.Errorf("failed to sign message %s", err.Error())
		return
	}

	t.Logf("******* signresponse: %+v", sigresponse)

	//ensure we haven't returned a derivation path
	if sigresponse.DerivationPath != nil {
		t.Errorf("Derivation path present for non-derived key, path %s", *sigresponse.DerivationPath)
		return
	}

	//ensure we haven't returned an address
	if sigresponse.Address != nil {
		t.Errorf("address present for non-derived key, address %s", *sigresponse.Address)
		return
	}

	verifyresponse, err := provide.VerifySignature(*token, vault.ID.String(), key.ID.String(), messageToSign, *sigresponse.Signature, nil)
	if err != nil {
		t.Errorf("failed to verify signature for vault: %s", err.Error())
		return
	}

	if verifyresponse.Verified != true {
		t.Error("failed to verify signature for vault")
		return
	}
}

func TestAPIVerifyEd25519Signature(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	key, err := keyFactory(*token, vault.ID.String(), "asymmetric", "sign/verify", cryptovault.KeySpecECCEd25519, "namey name", "cute description")
	if err != nil {
		t.Errorf("failed to create key; %s", err.Error())
		return
	}

	payloadBytes, _ := common.RandomBytes(1000)
	messageToSign := hex.EncodeToString(payloadBytes)
	sigresponse, err := provide.SignMessage(*token, vault.ID.String(), key.ID.String(), messageToSign, nil)
	if err != nil {
		t.Errorf("failed to sign message %s", err.Error())
		return
	}

	verifyresponse, err := provide.VerifySignature(*token, vault.ID.String(), key.ID.String(), messageToSign, *sigresponse.Signature, nil)
	if err != nil {
		t.Errorf("failed to verify signature for vault: %s", err.Error())
		return
	}

	if verifyresponse.Verified != true {
		t.Error("failed to verify signature for vault")
		return
	}
}

func TestAPIVerifyEd25519NKeySignature(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	key, err := keyFactory(*token, vault.ID.String(), "asymmetric", "sign/verify", cryptovault.KeySpecECCEd25519NKey, "namey name", "cute description")
	if err != nil {
		t.Errorf("failed to create key; %s", err.Error())
		return
	}

	payloadBytes, _ := common.RandomBytes(1000)
	messageToSign := hex.EncodeToString(payloadBytes)
	sigresponse, err := provide.SignMessage(*token, vault.ID.String(), key.ID.String(), messageToSign, nil)
	if err != nil {
		t.Errorf("failed to sign message %s", err.Error())
		return
	}

	verifyresponse, err := provide.VerifySignature(*token, vault.ID.String(), key.ID.String(), messageToSign, *sigresponse.Signature, nil)
	if err != nil {
		t.Errorf("failed to verify signature for vault: %s", err.Error())
		return
	}

	if verifyresponse.Verified != true {
		t.Error("failed to verify signature for vault")
		return
	}
}

func TestAPIVerifyRSA2048PS256Signature(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	key, err := keyFactory(*token, vault.ID.String(), "asymmetric", "sign/verify", cryptovault.KeySpecRSA2048, "namey name", "cute description")
	if err != nil {
		t.Errorf("failed to create key; %s", err.Error())
		return
	}

	payloadBytes, _ := common.RandomBytes(1000)
	messageToSign := hex.EncodeToString(payloadBytes)

	opts := map[string]interface{}{}
	json.Unmarshal([]byte(`{"algorithm":"PS256"}`), &opts)

	sigresponse, err := provide.SignMessage(*token, vault.ID.String(), key.ID.String(), messageToSign, opts)
	if err != nil {
		t.Errorf("failed to sign message %s", err.Error())
		return
	}

	verifyresponse, err := provide.VerifySignature(*token, vault.ID.String(), key.ID.String(), messageToSign, *sigresponse.Signature, opts)
	if err != nil {
		t.Errorf("failed to verify signature for vault: %s", err.Error())
		return
	}

	if verifyresponse.Verified != true {
		t.Error("failed to verify signature for vault")
		return
	}

}

func TestAPIEncrypt(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	key, err := keyFactory(*token, vault.ID.String(), "symmetric", "encrypt/decrypt", "aes-256-GCM", "namey name", "cute description")
	if err != nil {
		t.Errorf("failed to create key; %s", err.Error())
		return
	}

	data := common.RandomString(128)
	nonce := "1"

	_, err = provide.EncryptWithNonce(*token, vault.ID.String(), key.ID.String(), data, nonce)

	if err != nil {
		t.Errorf("failed to encrypt message for vault: %s", vault.ID)
		return
	}
}

func TestAPIChachaDecrypt(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	key, err := keyFactory(*token, vault.ID.String(), "symmetric", "encrypt/decrypt", "chaCha20", "namey name", "cute description")
	if err != nil {
		t.Errorf("failed to create key; %s", err.Error())
		return
	}

	data := common.RandomString(128)
	nonce := "1"

	encryptedDataResponse, err := provide.EncryptWithNonce(*token, vault.ID.String(), key.ID.String(), data, nonce)

	if err != nil {
		t.Errorf("failed to encrypt message for vault: %s", vault.ID)
		return
	}

	decryptedDataResponse, err := provide.Decrypt(*token, vault.ID.String(), key.ID.String(), map[string]interface{}{
		"data": encryptedDataResponse.Data,
	})

	if decryptedDataResponse.Data != data {
		t.Errorf("decrypted data mismatch, expected %s, got %s", data, decryptedDataResponse.Data)
		return
	}
}

func TestAPIDecrypt(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	key, err := keyFactory(*token, vault.ID.String(), "symmetric", "encrypt/decrypt", "aes-256-GCM", "namey name", "cute description")
	if err != nil {
		t.Errorf("failed to create key; %s", err.Error())
		return
	}

	data := common.RandomString(128)
	nonce := common.RandomString(12)

	encryptedDataResponse, err := provide.EncryptWithNonce(*token, vault.ID.String(), key.ID.String(), data, nonce)

	if err != nil {
		t.Errorf("failed to encrypt message for vault: %s", vault.ID)
		return
	}

	decryptedDataResponse, err := provide.Decrypt(*token, vault.ID.String(), key.ID.String(), map[string]interface{}{
		"data": encryptedDataResponse.Data,
	})

	if decryptedDataResponse.Data != data {
		t.Errorf("decrypted data mismatch, expected %s, got %s", data, decryptedDataResponse.Data)
		return
	}
}

func TestAPIDecryptNoNonce(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	key, err := keyFactory(*token, vault.ID.String(), "symmetric", "encrypt/decrypt", "aes-256-GCM", "namey name", "cute description")
	if err != nil {
		t.Errorf("failed to create key; %s", err.Error())
		return
	}

	data := common.RandomString(128)

	encryptedDataResponse, err := provide.Encrypt(*token, vault.ID.String(), key.ID.String(), data)

	if err != nil {
		t.Errorf("failed to encrypt message for vault: %s", vault.ID)
		return
	}

	decryptedDataResponse, err := provide.Decrypt(*token, vault.ID.String(), key.ID.String(), map[string]interface{}{
		"data": encryptedDataResponse.Data,
	})

	if decryptedDataResponse.Data != data {
		t.Errorf("decrypted data mismatch, expected %s, got %s", data, decryptedDataResponse.Data)
		return
	}
}

func TestCreateHDWallet(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	key, err := keyFactory(*token, vault.ID.String(), "asymmetric", "sign/verify", "bip39", "hdwallet", "integration test hd wallet")
	if err != nil {
		t.Errorf("failed to create key; %s", err.Error())
		return
	}

	if key.PublicKey == nil {
		t.Errorf("failed to assign xpub key on hd wallet; %s", key.ID)
		return
	}

	opts := map[string]interface{}{}
	json.Unmarshal([]byte(`{"hdwallet":{"coin":60, "index":0}}`), &opts)

	payloadBytes, _ := common.RandomBytes(32)
	messageToSign := hex.EncodeToString(payloadBytes)
	sigresponse, err := provide.SignMessage(*token, vault.ID.String(), key.ID.String(), messageToSign, opts)
	if err != nil {
		t.Errorf("failed to sign message %s", err.Error())
		return
	}

	verifyresponse, err := provide.VerifySignature(*token, vault.ID.String(), key.ID.String(), messageToSign, *sigresponse.Signature, opts)
	if err != nil {
		t.Errorf("failed to verify signature for vault: %s", err.Error())
		return
	}

	if verifyresponse.Verified != true {
		t.Errorf("failed to verify signature for vault")
		return
	}
}

func TestCreateHDWalletFailsWithInvalidCoin(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	key, err := keyFactory(*token, vault.ID.String(), "asymmetric", "sign/verify", "bip39", "hdwallet", "integration test hd wallet")
	if err != nil {
		t.Errorf("failed to create key; %s", err.Error())
		return
	}

	if key.PublicKey == nil {
		t.Errorf("failed to assign xpub key on hd wallet; %s", key.ID)
		return
	}

	opts := map[string]interface{}{}
	json.Unmarshal([]byte(`{"hdwallet":{"coin":61, "index":0}}`), &opts) // coin: 61 <-- this is not supported

	payloadBytes, _ := common.RandomBytes(32)
	messageToSign := hex.EncodeToString(payloadBytes)
	sigresponse, err := provide.SignMessage(*token, vault.ID.String(), key.ID.String(), messageToSign, opts)
	if err != nil {
		t.Errorf("failed to sign message %s", err.Error())
		return
	}

	verifyresponse, err := provide.VerifySignature(*token, vault.ID.String(), key.ID.String(), messageToSign, *sigresponse.Signature, opts)
	if err != nil {
		t.Errorf("failed to verify signature for vault: %s", err.Error())
		return
	}

	if verifyresponse.Verified != true {
		t.Errorf("failed to verify signature for vault")
		return
	}
}

func TestCreateHDWalletCoinAbbr(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	key, err := keyFactory(*token, vault.ID.String(), "asymmetric", "sign/verify", "bip39", "hdwallet", "integration test hd wallet")
	if err != nil {
		t.Errorf("failed to create key; %s", err.Error())
		return
	}

	if key.PublicKey == nil {
		t.Errorf("failed to assign xpub key on hd wallet; %s", key.ID)
		return
	}

	opts := map[string]interface{}{}
	json.Unmarshal([]byte(`{"hdwallet":{"coin_abbr":"ETH", "index":0}}`), &opts)

	payloadBytes, _ := common.RandomBytes(32)
	messageToSign := hex.EncodeToString(payloadBytes)
	sigresponse, err := provide.SignMessage(*token, vault.ID.String(), key.ID.String(), messageToSign, opts)
	if err != nil {
		t.Errorf("failed to sign message %s", err.Error())
		return
	}

	verifyresponse, err := provide.VerifySignature(*token, vault.ID.String(), key.ID.String(), messageToSign, *sigresponse.Signature, opts)
	if err != nil {
		t.Errorf("failed to verify signature for vault: %s", err.Error())
		return
	}

	if verifyresponse.Verified != true {
		t.Errorf("failed to verify signature for vault")
		return
	}
}

func TestHDWalletAutoSign(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	key, err := keyFactory(*token, vault.ID.String(), "asymmetric", "sign/verify", "bip39", "hdwallet", "integration test hd wallet")
	if err != nil {
		t.Errorf("failed to create key; %s", err.Error())
		return
	}

	for iteration := 0; iteration < 10; iteration++ {
		payloadBytes, _ := common.RandomBytes(32)
		messageToSign := hex.EncodeToString(payloadBytes)
		sigresponse, err := provide.SignMessage(*token, vault.ID.String(), key.ID.String(), messageToSign, nil)
		if err != nil {
			t.Errorf("failed to sign message %s", err.Error())
			return
		}

		//ensure we have returned a derivation path
		if sigresponse.DerivationPath == nil {
			t.Errorf("No derivation path returned for derived key sign operation")
			return
		}

		//ensure we have returned an address
		if sigresponse.Address == nil {
			t.Errorf("no address returned for derived key sign operation")
			return
		}

		// set up the verification options
		opts := map[string]interface{}{}
		options := fmt.Sprintf(`{"hdwallet":{"coin_abbr":"ETH", "index":%d}}`, iteration)
		json.Unmarshal([]byte(options), &opts)

		verifyresponse, err := provide.VerifySignature(*token, vault.ID.String(), key.ID.String(), messageToSign, *sigresponse.Signature, opts)
		if err != nil {
			t.Errorf("failed to verify signature for vault: %s", err.Error())
			return
		}

		if verifyresponse.Verified != true {
			t.Errorf("failed to verify signature for vault!")
			return
		}
	}
}

func TestListKeys(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	// set how many keys we're going to generate
	const numberOfKeys = 24
	var inputKey [numberOfKeys + 1]*provide.Key
	inputKey[0] = nil //ignoring the vault master key

	for looper := 1; looper <= numberOfKeys; looper++ {
		keyName := fmt.Sprintf("integration test ethereum key %d", looper)
		key, err := provide.CreateKey(*token, vault.ID.String(), map[string]interface{}{
			"type":        "asymmetric",
			"usage":       "sign/verify",
			"spec":        "secp256K1",
			"name":        keyName,
			"description": "organization eth/stablecoin wallet",
		})

		if err != nil {
			t.Errorf("failed to create key. error %s", err.Error())
		}

		inputKey[looper] = key

		if len(*inputKey[looper].Address) != 42 {
			t.Errorf("invalid address length for key 01. expected 42, got %d", len(*inputKey[looper].Address))
			return
		}
	}

	listVaultKeysResponse, err := provide.ListKeys(*token, vault.ID.String(), map[string]interface{}{})
	if err != nil {
		t.Errorf("failed to list keys. error %s", err.Error())
	}

	if len(listVaultKeysResponse) != numberOfKeys+1 {
		t.Errorf("invalid number of keys returned")
		return
	}

	var outputKey [numberOfKeys + 1]*provide.Key
	for looper := 0; looper <= numberOfKeys; looper++ {
		outputKey[looper] = listVaultKeysResponse[looper]

		if looper > 0 {
			if *inputKey[looper].Address != *outputKey[looper].Address {
				t.Errorf("address mismatch. expected %s, got %s", *inputKey[looper].Address, *outputKey[looper].Address)
			}

			if *inputKey[looper].Description != *outputKey[looper].Description {
				t.Errorf("description mismatch. expected %s, got %s", *inputKey[looper].Description, *outputKey[looper].Description)
			}

			if inputKey[looper].ID != outputKey[looper].ID {
				t.Errorf("id mismatch. expected %s, got %s", inputKey[looper].ID, outputKey[looper].ID)
			}

			if *inputKey[looper].Name != *outputKey[looper].Name {
				t.Errorf("name mismatch. expected %s, got %s", *inputKey[looper].Name, *outputKey[looper].Name)
			}

			if *inputKey[looper].Spec != *outputKey[looper].Spec {
				t.Errorf("spec mismatch. expected %s, got %s", *inputKey[looper].Spec, *outputKey[looper].Spec)
			}

			if *inputKey[looper].Type != *outputKey[looper].Type {
				t.Errorf("type mismatch. expected %s, got %s", *inputKey[looper].Type, *outputKey[looper].Type)
			}

			if *inputKey[looper].Usage != *outputKey[looper].Usage {
				t.Errorf("usage mismatch. expected %s, got %s", *inputKey[looper].Usage, *outputKey[looper].Usage)
			}

			if inputKey[looper].VaultID.String() != outputKey[looper].VaultID.String() {
				t.Errorf("vault_id mismatch. expected %s, got %s", inputKey[looper].VaultID, outputKey[looper].VaultID)
			}

			if *inputKey[looper].PublicKey != *outputKey[looper].PublicKey {
				t.Errorf("public_key mismatch. expected %s, got %s", *inputKey[looper].PublicKey, *outputKey[looper].PublicKey)
			}

			t.Logf("key %d of %d validated", looper, numberOfKeys)
		}
	}
}

func TestListKeys_Filtered(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	// generate a key that will be filtered out
	_, err = provide.CreateKey(*token, vault.ID.String(), map[string]interface{}{
		"type":        "asymmetric",
		"usage":       "sign/verify",
		"spec":        "babyjubjub",
		"name":        "babyjubjub key to be filtered out",
		"description": "baseline babyjubjub key",
	})

	if err != nil {
		t.Errorf("failed to create key. error %s", err.Error())
	}

	// set how many keys we're going to generate for the filter
	const numberOfKeys = 2
	var inputKey [numberOfKeys + 1]*provide.Key
	//inputKey[0] = nil //ignoring the vault master key

	for looper := 0; looper < numberOfKeys; looper++ {
		keyName := fmt.Sprintf("integration test ethereum key %d", looper)
		key, err := provide.CreateKey(*token, vault.ID.String(), map[string]interface{}{
			"type":        "asymmetric",
			"usage":       "sign/verify",
			"spec":        "SECP256k1",
			"name":        keyName,
			"description": "organization eth/stablecoin wallet",
		})

		if err != nil {
			t.Errorf("failed to create key. error %s", err.Error())
		}

		inputKey[looper] = key

		if len(*inputKey[looper].Address) != 42 {
			t.Errorf("invalid address length for key 01. expected 42, got %d", len(*inputKey[looper].Address))
			return
		}
	}

	// first run without filter
	listVaultKeysResponse, err := provide.ListKeys(*token, vault.ID.String(), map[string]interface{}{})
	if err != nil {
		t.Errorf("failed to list keys. error %s", err.Error())
	}

	if len(listVaultKeysResponse) != numberOfKeys+2 {
		t.Errorf("invalid number of keys returned")
		return
	}

	// now filter to just secp256k1
	listVaultKeysResponse, err = provide.ListKeys(*token, vault.ID.String(), map[string]interface{}{
		"spec": "secp256k1",
	})
	if err != nil {
		t.Errorf("failed to list keys. error %s", err.Error())
	}

	if len(listVaultKeysResponse) != numberOfKeys {
		t.Errorf("invalid number of secp256k1 keys returned")
		return
	}

	// now filter to babyjubjub
	listVaultKeysResponse, err = provide.ListKeys(*token, vault.ID.String(), map[string]interface{}{
		"spec": "babyJubJub",
	})
	if err != nil {
		t.Errorf("failed to list keys. error %s", err.Error())
	}

	if len(listVaultKeysResponse) != 1 {
		t.Errorf("invalid number of baby jub jub keys returned")
		return
	}

	// now filter to symmetric (should be just the master key)
	listVaultKeysResponse, err = provide.ListKeys(*token, vault.ID.String(), map[string]interface{}{
		"type": "symmetric",
	})
	if err != nil {
		t.Errorf("failed to list keys. error %s", err.Error())
	}

	if len(listVaultKeysResponse) != 1 {
		t.Errorf("invalid number of symmetric keys returned")
		return
	}

	// now filter to asymmetric (should be babyjubjub + numberOfKeys secp256k1 keys)
	listVaultKeysResponse, err = provide.ListKeys(*token, vault.ID.String(), map[string]interface{}{
		"type": "asymmetric",
	})
	if err != nil {
		t.Errorf("failed to list keys. error %s", err.Error())
	}

	if len(listVaultKeysResponse) != (numberOfKeys + 1) {
		t.Errorf("invalid number of asymmetric keys returned")
		return
	}

	//now check the value of all the secp256k1 keys added
	// now filter to babyjubjub
	listVaultKeysResponse, err = provide.ListKeys(*token, vault.ID.String(), map[string]interface{}{
		"spec": "secp256k1",
	})
	if err != nil {
		t.Errorf("failed to list keys. error %s", err.Error())
	}

	if len(listVaultKeysResponse) != numberOfKeys {
		t.Errorf("invalid number of secp256k1 keys returned")
		return
	}

	var outputKey [numberOfKeys + 1]*provide.Key
	for looper := 0; looper < numberOfKeys; looper++ {
		outputKey[looper] = listVaultKeysResponse[looper]

		if *inputKey[looper].Address != *outputKey[looper].Address {
			t.Errorf("address mismatch. expected %s, got %s", *inputKey[looper].Address, *outputKey[looper].Address)
		}

		if *inputKey[looper].Description != *outputKey[looper].Description {
			t.Errorf("description mismatch. expected %s, got %s", *inputKey[looper].Description, *outputKey[looper].Description)
		}

		if inputKey[looper].ID != outputKey[looper].ID {
			t.Errorf("id mismatch. expected %s, got %s", inputKey[looper].ID, outputKey[looper].ID)
		}

		if *inputKey[looper].Name != *outputKey[looper].Name {
			t.Errorf("name mismatch. expected %s, got %s", *inputKey[looper].Name, *outputKey[looper].Name)
		}

		if *inputKey[looper].Spec != *outputKey[looper].Spec {
			t.Errorf("spec mismatch. expected %s, got %s", *inputKey[looper].Spec, *outputKey[looper].Spec)
		}

		if *inputKey[looper].Type != *outputKey[looper].Type {
			t.Errorf("type mismatch. expected %s, got %s", *inputKey[looper].Type, *outputKey[looper].Type)
		}

		if *inputKey[looper].Usage != *outputKey[looper].Usage {
			t.Errorf("usage mismatch. expected %s, got %s", *inputKey[looper].Usage, *outputKey[looper].Usage)
		}

		if inputKey[looper].VaultID.String() != outputKey[looper].VaultID.String() {
			t.Errorf("vault_id mismatch. expected %s, got %s", inputKey[looper].VaultID, outputKey[looper].VaultID)
		}

		if *inputKey[looper].PublicKey != *outputKey[looper].PublicKey {
			t.Errorf("public_key mismatch. expected %s, got %s", *inputKey[looper].PublicKey, *outputKey[looper].PublicKey)
		}

		if *inputKey[looper].Spec == "secp256k1" {
			if outputKey[looper].Address == nil {
				t.Errorf("output address non-nil for %s key %s", *inputKey[looper].Spec, inputKey[looper].ID.String())
			}
		} else {
			if inputKey[looper].Address != nil {
				t.Errorf("input address was non-nil for %s key %s", *inputKey[looper].Spec, inputKey[looper].ID.String())
			}
		}

		t.Logf("key %d of %d validated", looper+1, numberOfKeys)
	}
}

func TestAPIDerivedChachaDecrypt(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	key, err := keyFactory(*token, vault.ID.String(), "symmetric", "encrypt/decrypt", "chacha20", "namey name", "cute description")
	if err != nil {
		t.Errorf("failed to create key; %s", err.Error())
		return
	}

	nonce := 1
	context := common.RandomString(32)
	name := "derived key 01"
	description := "derived key 01 description"

	derivedKey, err := provide.DeriveKey(*token, vault.ID.String(), key.ID.String(), map[string]interface{}{
		"nonce":       nonce,
		"context":     context,
		"name":        name,
		"description": description,
	})

	if err != nil {
		t.Errorf("failed to derive key for vault: %s", vault.ID)
		return
	}

	if *derivedKey.Name != name {
		t.Errorf("name field incorrect. expected %s, got %s", name, *derivedKey.Name)
		return
	}

	if *derivedKey.Description != description {
		t.Errorf("description field incorrect. expected %s, got %s", description, *derivedKey.Description)
		return
	}

	data := common.RandomString(128)

	encryptedDataResponse, err := provide.Encrypt(*token, derivedKey.VaultID.String(), derivedKey.ID.String(), data)

	if err != nil {
		t.Errorf("failed to encrypt message for vault: %s", vault.ID)
		return
	}

	decryptedDataResponse, err := provide.Decrypt(*token, derivedKey.VaultID.String(), derivedKey.ID.String(), map[string]interface{}{
		"data": encryptedDataResponse.Data,
	})

	if decryptedDataResponse.Data != data {
		t.Errorf("decrypted data mismatch, expected %s, got %s", data, decryptedDataResponse.Data)
		return
	}
}

func TestAPIDerivedChachaDecryptNoNonce(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	key, err := keyFactory(*token, vault.ID.String(), "symmetric", "encrypt/decrypt", "chacha20", "namey name", "cute description")
	if err != nil {
		t.Errorf("failed to create key; %s", err.Error())
		return
	}

	context := common.RandomString(32)
	name := "derived key 01"
	description := "derived key 01 description"

	derivedKey, err := provide.DeriveKey(*token, vault.ID.String(), key.ID.String(), map[string]interface{}{
		"context":     context,
		"name":        name,
		"description": description,
	})

	if err != nil {
		t.Errorf("failed to derive key for vault: %s", vault.ID)
		return
	}

	if *derivedKey.Name != name {
		t.Errorf("name field incorrect. expected %s, got %s", name, *derivedKey.Name)
		return
	}

	if *derivedKey.Description != description {
		t.Errorf("description field incorrect. expected %s, got %s", description, *derivedKey.Description)
		return
	}

	data := common.RandomString(128)

	encryptedDataResponse, err := provide.Encrypt(*token, derivedKey.VaultID.String(), derivedKey.ID.String(), data)

	if err != nil {
		t.Errorf("failed to encrypt message for vault: %s", vault.ID)
		return
	}

	decryptedDataResponse, err := provide.Decrypt(*token, derivedKey.VaultID.String(), derivedKey.ID.String(), map[string]interface{}{
		"data": encryptedDataResponse.Data,
	})

	if decryptedDataResponse.Data != data {
		t.Errorf("decrypted data mismatch, expected %s, got %s", data, decryptedDataResponse.Data)
		return
	}
}

func TestAPIDerivedNonChachaDecryptNoNonce(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	key, err := keyFactory(*token, vault.ID.String(), "symmetric", "encrypt/decrypt", "aes-256-gcm", "namey name", "cute description")
	if err != nil {
		t.Errorf("failed to create key; %s", err.Error())
		return
	}

	context := common.RandomString(32)
	name := "derived key 01"
	description := "derived key 01 description"

	_, err = provide.DeriveKey(*token, vault.ID.String(), key.ID.String(), map[string]interface{}{
		"context":     context,
		"name":        name,
		"description": description,
	})

	if err == nil {
		t.Errorf("incorrectly derived non-chacha20 key for vault: %s.", vault.ID)
		return
	}

	if err != nil {
		t.Logf("correctly returned error deriving non-chacha20 key. Error: %s", err.Error())
	}
}

func TestAPIDeriveBIP39(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	key, err := keyFactory(*token, vault.ID.String(), "asymmetric", "sign/verify", "bip39", "namey name", "cute description")
	if err != nil {
		t.Errorf("failed to create key; %s", err.Error())
		return
	}

	derivedKey, err := provide.DeriveKey(*token, vault.ID.String(), key.ID.String(), map[string]interface{}{})

	if err != nil {
		t.Errorf("failed to derive key for vault: %s", vault.ID)
		return
	}

	if derivedKey.Address == nil {
		t.Errorf("address should be non-nil for derived secp256k1 BIP39 HD wallet key")
		return
	}

	if derivedKey.HDDerivationPath == nil {
		t.Errorf("derivation path should be non-nil for derived secp256k1 BIP39 HD wallet key")
		return
	}

	if derivedKey.PublicKey == nil {
		t.Errorf("public key should be non-nil for derived secp256k1 BIP39 HD wallet key")
		return
	}
}

func TestEphemeralCreation(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	tt := []struct {
		Name        string
		Description string
		Type        string
		Usage       string
		Spec        string
		KSpec       string
	}{
		{"ephemeral key", "ephemeral key description", "asymmetric", "sign/verify", "c25519", cryptovault.KeySpecECCC25519},
		{"ephemeral key", "ephemeral key description", "asymmetric", "sign/verify", "ed25519", cryptovault.KeySpecECCEd25519},
		{"ephemeral key", "ephemeral key description", "asymmetric", "sign/verify", "Secp256k1", cryptovault.KeySpecECCSecp256k1},
		{"ephemeral key", "ephemeral key description", "asymmetric", "sign/verify", "BabyJubJub", cryptovault.KeySpecECCBabyJubJub},
		{"ephemeral key", "ephemeral key description", "asymmetric", "sign/verify", "Bip39", cryptovault.KeySpecECCBIP39},
		{"ephemeral key", "ephemeral key description", "asymmetric", "sign/verify", "Rsa-2048", cryptovault.KeySpecRSA2048},
		{"ephemeral key", "ephemeral key description", "asymmetric", "sign/verify", "Rsa-3072", cryptovault.KeySpecRSA3072},
		{"ephemeral key", "ephemeral key description", "asymmetric", "sign/verify", "Rsa-4096", cryptovault.KeySpecRSA4096},
		{"ephemeral key", "ephemeral key description", "symmetric", "encrypt/decrypt", "Aes-256-Gcm", cryptovault.KeySpecAES256GCM},
		{"ephemeral key", "ephemeral key description", "symmetric", "encrypt/decrypt", "chacha20", cryptovault.KeySpecChaCha20},
	}

	for _, tc := range tt {
		key, err := keyFactoryEphemeral(*token, vault.ID.String(), tc.Type, tc.Usage, tc.Spec, tc.Name, tc.Description)
		if err != nil {
			t.Errorf("failed to create key; %s", err.Error())
			return
		}

		if *key.Name != tc.Name {
			t.Errorf("name mismatch. expected %s, got %s", tc.Name, *key.Name)
			return
		}

		if *key.Description != tc.Description {
			t.Errorf("description mismatch. expected %s, got %s", tc.Description, *key.Description)
			return
		}

		if *key.Type != tc.Type {
			t.Errorf("type mismatch. expected %s, got %s", tc.Type, *key.Type)
			return
		}

		if *key.Usage != tc.Usage {
			t.Errorf("usage mismatch. expected %s, got %s", tc.Usage, *key.Usage)
			return
		}
		if *key.Spec != tc.KSpec {
			t.Errorf("spec mismatch. expected %s, got %s", tc.Spec, *key.Spec)
			return
		}

		switch *key.Spec {
		case cryptovault.KeySpecECCC25519:
			if key.PrivateKey == nil {
				t.Errorf("no private key returned for ephemeral %s key", tc.Spec)
			}
			if key.PublicKey == nil {
				t.Errorf("no private key returned for ephemeral %s key", tc.Spec)
			}
		case cryptovault.KeySpecECCEd25519:
			if key.Seed == nil {
				t.Errorf("no seed returned for ephemeral %s key", tc.Spec)
			}
			if key.PublicKey == nil {
				t.Errorf("no private key returned for ephemeral %s key", tc.Spec)
			}
		case cryptovault.KeySpecECCSecp256k1:
			if key.PrivateKey == nil {
				t.Errorf("no private key returned for ephemeral %s key", tc.Spec)
			}
			if key.PublicKey == nil {
				t.Errorf("no public key returned for ephemeral %s key", tc.Spec)
			}
			if key.Address == nil {
				t.Errorf("no address returned for ephemeral %s key", tc.Spec)
			}
		case cryptovault.KeySpecECCBabyJubJub:
			if key.PrivateKey == nil {
				t.Errorf("no private key returned for ephemeral %s key", tc.Spec)
			}
			if key.PublicKey == nil {
				t.Errorf("no public key returned for ephemeral %s key", tc.Spec)
			}
		case cryptovault.KeySpecECCBIP39:
			if key.Seed == nil {
				t.Errorf("no seed returned for ephemeral %s key", tc.Spec)
			}
		case cryptovault.KeySpecRSA2048:
			if key.PrivateKey == nil {
				t.Errorf("no private key returned for ephemeral %s key", tc.Spec)
			}
			if key.PublicKey == nil {
				t.Errorf("no public key returned for ephemeral %s key", tc.Spec)
			}
		case cryptovault.KeySpecRSA3072:
			if key.PrivateKey == nil {
				t.Errorf("no private key returned for ephemeral %s key", tc.Spec)
			}
			if key.PublicKey == nil {
				t.Errorf("no public key returned for ephemeral %s key", tc.Spec)
			}
		case cryptovault.KeySpecRSA4096:
			if key.PrivateKey == nil {
				t.Errorf("no private key returned for ephemeral %s key", tc.Spec)
			}
			if key.PublicKey == nil {
				t.Errorf("no public key returned for ephemeral %s key", tc.Spec)
			}
		case cryptovault.KeySpecAES256GCM:
			if key.PrivateKey == nil {
				t.Errorf("no private key returned for ephemeral %s key", tc.Spec)
			}
		case cryptovault.KeySpecChaCha20:
			if key.Seed == nil {
				t.Errorf("no seed returned for ephemeral %s key", tc.Spec)
			}
		default:
			t.Errorf("unknown key spec generated: %s", tc.Spec)
			return
		}
	}
}

func TestNonEphemeralCreation(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	tt := []struct {
		Name        string
		Description string
		Type        string
		Usage       string
		Spec        string
		KSpec       string
	}{
		{"ephemeral key", "ephemeral key description", "asymmetric", "sign/verify", "c25519", cryptovault.KeySpecECCC25519},
		{"ephemeral key", "ephemeral key description", "asymmetric", "sign/verify", "ed25519", cryptovault.KeySpecECCEd25519},
		{"ephemeral key", "ephemeral key description", "asymmetric", "sign/verify", "Secp256k1", cryptovault.KeySpecECCSecp256k1},
		{"ephemeral key", "ephemeral key description", "asymmetric", "sign/verify", "BabyJubJub", cryptovault.KeySpecECCBabyJubJub},
		{"ephemeral key", "ephemeral key description", "asymmetric", "sign/verify", "Bip39", cryptovault.KeySpecECCBIP39},
		{"ephemeral key", "ephemeral key description", "asymmetric", "sign/verify", "Rsa-2048", cryptovault.KeySpecRSA2048},
		{"ephemeral key", "ephemeral key description", "asymmetric", "sign/verify", "Rsa-3072", cryptovault.KeySpecRSA3072},
		{"ephemeral key", "ephemeral key description", "asymmetric", "sign/verify", "Rsa-4096", cryptovault.KeySpecRSA4096},
		{"ephemeral key", "ephemeral key description", "symmetric", "encrypt/decrypt", "Aes-256-Gcm", cryptovault.KeySpecAES256GCM},
		{"ephemeral key", "ephemeral key description", "symmetric", "encrypt/decrypt", "chacha20", cryptovault.KeySpecChaCha20},
	}

	for _, tc := range tt {
		key, err := keyFactory(*token, vault.ID.String(), tc.Type, tc.Usage, tc.Spec, tc.Name, tc.Description)
		if err != nil {
			t.Errorf("failed to create key; %s", err.Error())
			return
		}

		if *key.Name != tc.Name {
			t.Errorf("name mismatch. expected %s, got %s", tc.Name, *key.Name)
			return
		}

		if *key.Description != tc.Description {
			t.Errorf("description mismatch. expected %s, got %s", tc.Description, *key.Description)
			return
		}

		if *key.Type != tc.Type {
			t.Errorf("type mismatch. expected %s, got %s", tc.Type, *key.Type)
			return
		}

		if *key.Usage != tc.Usage {
			t.Errorf("usage mismatch. expected %s, got %s", tc.Usage, *key.Usage)
			return
		}
		if *key.Spec != tc.KSpec {
			t.Errorf("spec mismatch. expected %s, got %s", tc.Spec, *key.Spec)
			return
		}

		switch *key.Spec {
		case cryptovault.KeySpecECCC25519:
			if key.PrivateKey != nil {
				t.Errorf("private key returned for regular %s key", tc.Spec)
			}
			if key.PublicKey == nil {
				t.Errorf("no private key returned for regular %s key", tc.Spec)
			}
		case cryptovault.KeySpecECCEd25519:
			if key.Seed != nil {
				t.Errorf("seed returned for regular %s key", tc.Spec)
			}
			if key.PublicKey == nil {
				t.Errorf("no private key returned for regular %s key", tc.Spec)
			}
		case cryptovault.KeySpecECCSecp256k1:
			if key.PrivateKey != nil {
				t.Errorf("private key returned for regular %s key", tc.Spec)
			}
			if key.PublicKey == nil {
				t.Errorf("no public key returned for regular %s key", tc.Spec)
			}
			if key.Address == nil {
				t.Errorf("no address returned for regular %s key", tc.Spec)
			}
		case cryptovault.KeySpecECCBabyJubJub:
			if key.PrivateKey != nil {
				t.Errorf("private key returned for regular %s key", tc.Spec)
			}
			if key.PublicKey == nil {
				t.Errorf("no public key returned for regular %s key", tc.Spec)
			}
		case cryptovault.KeySpecECCBIP39:
			if key.Seed != nil {
				t.Errorf("seed returned for regular %s key", tc.Spec)
			}
		case cryptovault.KeySpecRSA2048:
			if key.PrivateKey != nil {
				t.Errorf("private key returned for regular %s key", tc.Spec)
			}
			if key.PublicKey == nil {
				t.Errorf("no public key returned for regular %s key", tc.Spec)
			}
		case cryptovault.KeySpecRSA3072:
			if key.PrivateKey != nil {
				t.Errorf("private key returned for regular %s key", tc.Spec)
			}
			if key.PublicKey == nil {
				t.Errorf("no public key returned for regular %s key", tc.Spec)
			}
		case cryptovault.KeySpecRSA4096:
			if key.PrivateKey != nil {
				t.Errorf("private key returned for regular %s key", tc.Spec)
			}
			if key.PublicKey == nil {
				t.Errorf("no public key returned for regular %s key", tc.Spec)
			}
		case cryptovault.KeySpecAES256GCM:
			if key.PrivateKey != nil {
				t.Errorf("private key returned for regular %s key", tc.Spec)
			}
		case cryptovault.KeySpecChaCha20:
			if key.Seed != nil {
				t.Errorf("seed returned for regular %s key", tc.Spec)
			}
		default:
			t.Errorf("unknown key spec generated: %s", tc.Spec)
			return
		}
	}
}

func TestArbitrarySignatureSecp256k1(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	key, err := keyFactory(*token, vault.ID.String(), "asymmetric", "sign/verify", cryptovault.KeySpecECCSecp256k1, "namey name", "cute description")
	if err != nil {
		t.Errorf("failed to create key; %s", err.Error())
		return
	}

	payloadBytes, _ := common.RandomBytes(32)
	messageToSign := hex.EncodeToString(payloadBytes)

	opts := map[string]interface{}{}

	sigresponse, err := provide.SignMessage(*token, vault.ID.String(), key.ID.String(), messageToSign, opts)
	if err != nil {
		t.Errorf("failed to sign message %s", err.Error())
		return
	}

	verifyresponse, err := provide.VerifySignature(*token, vault.ID.String(), key.ID.String(), messageToSign, *sigresponse.Signature, opts)
	if err != nil {
		t.Errorf("failed to verify signature for vault: %s", err.Error())
		return
	}

	if verifyresponse.Verified != true {
		t.Error("failed to verify signature for vault")
		return
	}
}

func TestArbitrarySignatureSecp256k1_ShouldFail(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	key, err := keyFactory(*token, vault.ID.String(), "asymmetric", "sign/verify", cryptovault.KeySpecECCSecp256k1, "namey name", "cute description")
	if err != nil {
		t.Errorf("failed to create key; %s", err.Error())
		return
	}

	payloadBytes, _ := common.RandomBytes(33)
	messageToSign := hex.EncodeToString(payloadBytes)

	opts := map[string]interface{}{}
	//json.Unmarshal([]byte(`{"algorithm":"PS256"}`), &opts)

	_, err = provide.SignMessage(*token, vault.ID.String(), key.ID.String(), messageToSign, opts)
	if err == nil {
		t.Errorf("signed 33-byte message with secp256k1 - should only sign 32-byte messages")
		return
	}
}

func TestArbitrarySignatureEd25519(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	key, err := keyFactory(*token, vault.ID.String(), "asymmetric", "sign/verify", cryptovault.KeySpecECCEd25519, "namey name", "cute description")
	if err != nil {
		t.Errorf("failed to create key; %s", err.Error())
		return
	}

	payloadBytes, _ := common.RandomBytes(1000)
	messageToSign := hex.EncodeToString(payloadBytes)

	opts := map[string]interface{}{}

	sigresponse, err := provide.SignMessage(*token, vault.ID.String(), key.ID.String(), messageToSign, opts)
	if err != nil {
		t.Errorf("failed to sign message %s", err.Error())
		return
	}

	verifyresponse, err := provide.VerifySignature(*token, vault.ID.String(), key.ID.String(), messageToSign, *sigresponse.Signature, opts)
	if err != nil {
		t.Errorf("failed to verify signature for vault: %s", err.Error())
		return
	}

	if verifyresponse.Verified != true {
		t.Error("failed to verify signature for vault")
		return
	}
}

func TestDetachedSignatureVerification(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	tt := []struct {
		Name        string
		Description string
		Type        string
		Usage       string
		Spec        string
		Options     string
	}{
		{"regular key", "regular key description", cryptovault.KeyTypeAsymmetric, cryptovault.KeyUsageSignVerify, "ed25519", ""},
		{"regular key", "regular key description", cryptovault.KeyTypeAsymmetric, cryptovault.KeyUsageSignVerify, "babyjubjub", ""},
		{"regular key", "regular key description", cryptovault.KeyTypeAsymmetric, cryptovault.KeyUsageSignVerify, "rsa-2048", `{"algorithm":"PS256"}`},
		{"regular key", "regular key description", cryptovault.KeyTypeAsymmetric, cryptovault.KeyUsageSignVerify, "rsa-3072", `{"algorithm":"PS256"}`},
		{"regular key", "regular key description", cryptovault.KeyTypeAsymmetric, cryptovault.KeyUsageSignVerify, "rsa-4096", `{"algorithm":"PS256"}`},
		{"regular key", "regular key description", cryptovault.KeyTypeAsymmetric, cryptovault.KeyUsageSignVerify, "rsa-2048", `{"algorithm":"PS384"}`},
		{"regular key", "regular key description", cryptovault.KeyTypeAsymmetric, cryptovault.KeyUsageSignVerify, "rsa-3072", `{"algorithm":"PS384"}`},
		{"regular key", "regular key description", cryptovault.KeyTypeAsymmetric, cryptovault.KeyUsageSignVerify, "rsa-4096", `{"algorithm":"PS384"}`},
		{"regular key", "regular key description", cryptovault.KeyTypeAsymmetric, cryptovault.KeyUsageSignVerify, "rsa-2048", `{"algorithm":"PS512"}`},
		{"regular key", "regular key description", cryptovault.KeyTypeAsymmetric, cryptovault.KeyUsageSignVerify, "rsa-3072", `{"algorithm":"PS512"}`},
		{"regular key", "regular key description", cryptovault.KeyTypeAsymmetric, cryptovault.KeyUsageSignVerify, "rsa-4096", `{"algorithm":"PS512"}`},
		{"regular key", "regular key description", cryptovault.KeyTypeAsymmetric, cryptovault.KeyUsageSignVerify, "rsa-2048", `{"algorithm":"RS256"}`},
		{"regular key", "regular key description", cryptovault.KeyTypeAsymmetric, cryptovault.KeyUsageSignVerify, "rsa-3072", `{"algorithm":"RS256"}`},
		{"regular key", "regular key description", cryptovault.KeyTypeAsymmetric, cryptovault.KeyUsageSignVerify, "rsa-4096", `{"algorithm":"RS256"}`},
		{"regular key", "regular key description", cryptovault.KeyTypeAsymmetric, cryptovault.KeyUsageSignVerify, "rsa-2048", `{"algorithm":"RS384"}`},
		{"regular key", "regular key description", cryptovault.KeyTypeAsymmetric, cryptovault.KeyUsageSignVerify, "rsa-3072", `{"algorithm":"RS384"}`},
		{"regular key", "regular key description", cryptovault.KeyTypeAsymmetric, cryptovault.KeyUsageSignVerify, "rsa-4096", `{"algorithm":"RS384"}`},
		{"regular key", "regular key description", cryptovault.KeyTypeAsymmetric, cryptovault.KeyUsageSignVerify, "rsa-2048", `{"algorithm":"RS512"}`},
		{"regular key", "regular key description", cryptovault.KeyTypeAsymmetric, cryptovault.KeyUsageSignVerify, "rsa-3072", `{"algorithm":"RS512"}`},
		{"regular key", "regular key description", cryptovault.KeyTypeAsymmetric, cryptovault.KeyUsageSignVerify, "rsa-4096", `{"algorithm":"RS512"}`},
		{"regular key", "regular key description", cryptovault.KeyTypeAsymmetric, cryptovault.KeyUsageSignVerify, "secp256K1", ""},
	}

	for _, tc := range tt {

		key, err := keyFactory(*token, vault.ID.String(), tc.Type, tc.Usage, tc.Spec, "namey name", "cute description")
		if err != nil {
			t.Errorf("failed to create key; %s", err.Error())
			return
		}

		payloadBytes, _ := common.RandomBytes(32)
		messageToSign := hex.EncodeToString(payloadBytes)

		opts := map[string]interface{}{}
		if tc.Options != "" {
			json.Unmarshal([]byte(tc.Options), &opts)
		}

		sigresponse, err := provide.SignMessage(*token, vault.ID.String(), key.ID.String(), messageToSign, opts)
		if err != nil {
			t.Errorf("failed to sign message %s", err.Error())
			return
		}

		verifyresponse, err := provide.VerifySignature(*token, vault.ID.String(), key.ID.String(), messageToSign, *sigresponse.Signature, opts)
		if err != nil {
			t.Errorf("failed to verify signature for vault: %s", err.Error())
			return
		}

		if verifyresponse.Verified != true {
			t.Error("failed to verify signature for vault")
			return
		}

		detachedverifyresponse, err := provide.VerifyDetachedSignature(*token, tc.Spec, messageToSign, *sigresponse.Signature, *key.PublicKey, opts)
		if err != nil {
			t.Errorf("failed to verify detached signature: %s", err.Error())
			return
		}

		if detachedverifyresponse.Verified != true {
			t.Errorf("failed to verify detached signature for %s key type", tc.Spec)
			return
		}
	}
}

func TestDetachedSignatureVerification_ShouldFail(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	tt := []struct {
		Name        string
		Description string
		Type        string
		Usage       string
		Spec        string
		Options     string
	}{
		{"regular key", "regular key description", cryptovault.KeyTypeAsymmetric, cryptovault.KeyUsageSignVerify, cryptovault.KeySpecRSA2048, `{"algorithm":"PS512"}`},
		{"regular key", "regular key description", cryptovault.KeyTypeAsymmetric, cryptovault.KeyUsageSignVerify, cryptovault.KeySpecRSA3072, `{"algorithm":"PS512"}`},
		{"regular key", "regular key description", cryptovault.KeyTypeAsymmetric, cryptovault.KeyUsageSignVerify, cryptovault.KeySpecRSA4096, `{"algorithm":"PS512"}`},
	}

	for _, tc := range tt {

		key, err := keyFactory(*token, vault.ID.String(), tc.Type, tc.Usage, tc.Spec, "namey name", "cute description")
		if err != nil {
			t.Errorf("failed to create key; %s", err.Error())
			return
		}

		payloadBytes, _ := common.RandomBytes(32)
		messageToSign := hex.EncodeToString(payloadBytes)

		opts := map[string]interface{}{}
		if tc.Options != "" {
			json.Unmarshal([]byte(tc.Options), &opts)
		}

		sigresponse, err := provide.SignMessage(*token, vault.ID.String(), key.ID.String(), messageToSign, opts)
		if err != nil {
			t.Errorf("failed to sign message %s", err.Error())
			return
		}

		verifyresponse, err := provide.VerifySignature(*token, vault.ID.String(), key.ID.String(), messageToSign, *sigresponse.Signature, opts)
		if err != nil {
			t.Errorf("failed to verify signature for vault: %s", err.Error())
			return
		}

		if verifyresponse.Verified != true {
			t.Error("failed to verify signature for vault")
			return
		}

		// now we will run the detached verification with invalid parameters for each of the test keys
		// testing no spec, expect 422
		_, err = provide.VerifyDetachedSignature(*token, "", messageToSign, *sigresponse.Signature, *key.PublicKey, opts)
		if err == nil {
			t.Errorf("verified invalid detached signature with no spec")
		}

		// testing no message, expecting 422
		_, err = provide.VerifyDetachedSignature(*token, tc.Spec, "", *sigresponse.Signature, *key.PublicKey, opts)
		if err == nil {
			t.Errorf("verified invalid detached signature with no message")
		}

		// testing no signature, expecting 422
		_, err = provide.VerifyDetachedSignature(*token, tc.Spec, messageToSign, "", *key.PublicKey, opts)
		if err == nil {
			t.Errorf("verified invalid detached signature with no signature")
		}

		// testing no pubkey, expecting 422
		_, err = provide.VerifyDetachedSignature(*token, tc.Spec, messageToSign, *sigresponse.Signature, "", opts)
		if err == nil {
			t.Errorf("verified invalid detached signature with no public key")
		}

		// testing no algorithm (for RSA), expecting 422
		_, err = provide.VerifyDetachedSignature(*token, tc.Spec, messageToSign, *sigresponse.Signature, *key.PublicKey, map[string]interface{}{})
		if err == nil {
			t.Errorf("verified invalid detached RSA signature with no options")
			return
		}

		// testing invalid spec, will return an error with the input issue
		_, err = provide.VerifyDetachedSignature(*token, "invalid_spec", messageToSign, *sigresponse.Signature, *key.PublicKey, opts)
		if err == nil {
			t.Errorf("verified invalid spec signature")
			return
		}

		invalidPayload, _ := common.RandomBytes(32)
		invalidMessage := hex.EncodeToString(invalidPayload)

		// CHECKME testing invalid signature
		// returns a 200 with verified false (for consistency)
		// because the only thing that has gone wrong is the signature is invalid
		// as opposed to a parameter error
		verifyresponse, err = provide.VerifyDetachedSignature(*token, tc.Spec, invalidMessage, *sigresponse.Signature, *key.PublicKey, opts)
		if verifyresponse.Verified != false {
			t.Errorf("verified signature with invalid message")
			return
		}

	}
}

func TestCreateHDWalletWithSeed(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	seed := "traffic charge swing glimpse will citizen push mutual embrace volcano siege identify gossip battle casual exit enrich unlock muscle vast female initial please day"
	key, err := keyFactoryWithSeed(*token, vault.ID.String(), "asymmetric", "sign/verify", "bip39", "hdwallet", "integration test hd wallet", seed)
	if err != nil {
		t.Errorf("failed to create key; %s", err.Error())
		return
	}

	//
	if key.PublicKey == nil {
		t.Errorf("failed to assign xpub key on hd wallet; %s", key.ID)
		return
	}

	t.Logf("publickey:\n\t%s\n", *key.PublicKey)
	key2, err := keyFactoryWithSeed(*token, vault.ID.String(), "asymmetric", "sign/verify", "bip39", "hdwallet", "integration test hd wallet", seed)
	if err != nil {
		t.Errorf("failed to create key; %s", err.Error())
		return
	}

	if key2.PublicKey == nil {
		t.Errorf("failed to assign xpub key on hd wallet; %s", key2.ID)
		return
	}

	if *key.PublicKey != *key2.PublicKey {
		t.Errorf("deterministic wallet with seed did not create the same key")
		return
	}
}

func TestCreateHDWalletWithInvalidSeed(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	seed := "kraffic charge swing glimpse will citizen push mutual embrace volcano siege identify gossip battle casual exit enrich unlock muscle vast female initial please day"
	_, err = keyFactoryWithSeed(*token, vault.ID.String(), "asymmetric", "sign/verify", "bip39", "hdwallet", "integration test hd wallet", seed)
	if err == nil {
		t.Errorf("created HD wallet with invalid seed")
		return
	}

	if err != nil {
		t.Logf("error received: %s", err.Error())
	}
}

func TestHDWalletSeedAutoSign(t *testing.T) {

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	seed := "traffic charge swing glimpse will citizen push mutual embrace volcano siege identify gossip battle casual exit enrich unlock muscle vast female initial please day"
	key, err := keyFactoryWithSeed(*token, vault.ID.String(), "asymmetric", "sign/verify", "bip39", "hdwallet", "integration test hd wallet", seed)
	if err != nil {
		t.Errorf("failed to create key; %s", err.Error())
		return
	}

	key2, err := keyFactoryWithSeed(*token, vault.ID.String(), "asymmetric", "sign/verify", "bip39", "hdwallet", "integration test hd wallet", seed)
	if err != nil {
		t.Errorf("failed to create key; %s", err.Error())
		return
	}

	for iteration := 0; iteration < 10; iteration++ {
		payloadBytes, _ := common.RandomBytes(32)
		messageToSign := hex.EncodeToString(payloadBytes)

		sigresponse, err := provide.SignMessage(*token, vault.ID.String(), key.ID.String(), messageToSign, nil)
		if err != nil {
			t.Errorf("failed to sign message %s", err.Error())
			return
		}

		sigresponse2, err := provide.SignMessage(*token, vault.ID.String(), key2.ID.String(), messageToSign, nil)
		if err != nil {
			t.Errorf("failed to sign message %s", err.Error())
			return
		}

		if *sigresponse.Signature != *sigresponse2.Signature {
			t.Errorf("mismatch in signatures from key with provided seed")
			return
		}

		// set up the verification options
		opts := map[string]interface{}{}
		options := fmt.Sprintf(`{"hdwallet":{"coin_abbr":"ETH", "index":%d}}`, iteration)
		json.Unmarshal([]byte(options), &opts)

		verifyresponse, err := provide.VerifySignature(*token, vault.ID.String(), key.ID.String(), messageToSign, *sigresponse.Signature, opts)
		if err != nil {
			t.Errorf("failed to verify signature for vault: %s", err.Error())
			return
		}

		verifyresponse2, err := provide.VerifySignature(*token, vault.ID.String(), key2.ID.String(), messageToSign, *sigresponse2.Signature, opts)
		if err != nil {
			t.Errorf("failed to verify signature 2 for vault: %s", err.Error())
			return
		}

		if verifyresponse.Verified != true {
			t.Errorf("failed to verify signature for vault!")
			return
		}

		if verifyresponse2.Verified != true {
			t.Errorf("failed to verify signature 2 for vault!")
			return
		}

	}
}

func TestHDWalletSeedLedgerDerivationPath(t *testing.T) {

	// this test will generate an ethereum address using a particular derivation path
	// using the mechanism that the ledger wallet uses to create new ethereum addresses
	// we will then confirm that repeating this process
	// for two different keys using the same seed phrase
	// generates the same ETH address

	token, err := userTokenFactory()
	if err != nil {
		t.Errorf("failed to create token; %s", err.Error())
		return
	}

	vault, err := vaultFactory(*token, "vaulty vault", "just a vault with a key")
	if err != nil {
		t.Errorf("failed to create vault; %s", err.Error())
		return
	}

	seed := "traffic charge swing glimpse will citizen push mutual embrace volcano siege identify gossip battle casual exit enrich unlock muscle vast female initial please day"

	key, err := keyFactoryWithSeed(*token, vault.ID.String(), "asymmetric", "sign/verify", "bip39", "hdwallet", "integration test hd wallet", seed)
	if err != nil {
		t.Errorf("failed to create key; %s", err.Error())
		return
	}

	key2, err := keyFactoryWithSeed(*token, vault.ID.String(), "asymmetric", "sign/verify", "bip39", "hdwallet", "integration test hd wallet", seed)
	if err != nil {
		t.Errorf("failed to create key; %s", err.Error())
		return
	}

	// set up the verification options using a ledger-style Account path
	opts := map[string]interface{}{}
	path := `m/44'/60'/2'/0/0`
	options := fmt.Sprintf(`{"hdwallet":{"hd_derivation_path":"%s"}}`, path)
	json.Unmarshal([]byte(options), &opts)

	payloadBytes, _ := common.RandomBytes(32)
	messageToSign := hex.EncodeToString(payloadBytes)

	sigresponse, err := provide.SignMessage(*token, vault.ID.String(), key.ID.String(), messageToSign, opts)
	if err != nil {
		t.Errorf("failed to sign message %s", err.Error())
		return
	}

	sigresponse2, err := provide.SignMessage(*token, vault.ID.String(), key2.ID.String(), messageToSign, opts)
	if err != nil {
		t.Errorf("failed to sign message %s", err.Error())
		return
	}

	if *sigresponse.Signature != *sigresponse2.Signature {
		t.Errorf("mismatch in signatures from derived key with provided seed")
		return
	}

	if *sigresponse.Address != *sigresponse2.Address {
		t.Errorf("mismatch in generated address from derived key and seed")
		return
	}
	if *sigresponse.DerivationPath != *sigresponse2.DerivationPath {
		t.Errorf("mismatch in derivation path")
		return
	}

	if *sigresponse.DerivationPath != path {
		t.Errorf("returned derivation path does not correspond to provided derivation path. Expected: %s, received %s", path, *sigresponse.DerivationPath)
		return
	}
}
