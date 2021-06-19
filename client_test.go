package chia_client

import (
	"encoding/hex"
	"encoding/json"
	bls "github.com/chuwt/chia-bls-go"
	"testing"
)

var testClient *ChiaClient

func init() {
	testClient = NewChiaClient(
		"https://192.168.1.58:8555",
		TlsCertOpt(
			"./ssl/full_node/private_full_node.crt",
			"./ssl/full_node/private_full_node.key",
		),
	)
}

func TestClient(t *testing.T) {
	coins, err := testClient.GetCoins(GetCoinsReq{
		Address: "xch1935w6gvt60wqzy3h5xmecfnfjp0hv0wr78k5l6ndmvtkp5fwmduqn90u5q",
		//Address: "xch14452k0srjew8f865ej3dj7wgfc2qg5t0epzjmg7pwca03z3pkl4q2ekruc",
		Start: 337682,
		End:   439020,
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, coin := range coins.CoinRecords {
		t.Log(coin.Coin.ToJson())
	}
}

func TestSendTx(t *testing.T) {
	// get this tx json string from github.com/chuwt/chia-tx
	// this tx json is signed
	reqBytes := []byte(`{"coin_solutions": [{"coin": {"parent_coin_info": "0xb5d31c65960840ea826be97ef7dae140a680a047d48434475eef8bd9062b63e8", "puzzle_hash": "0x2c68ed218bd3dc011237a1b79c2669905f763dc3f1ed4fea6ddb1760d12edb78", "amount": 80}, "puzzle_reveal": "0xff02ffff01ff02ffff01ff02ffff03ff0bffff01ff02ffff03ffff09ff05ffff1dff0bffff1effff0bff0bffff02ff06ffff04ff02ffff04ff17ff8080808080808080ffff01ff02ff17ff2f80ffff01ff088080ff0180ffff01ff04ffff04ff04ffff04ff05ffff04ffff02ff06ffff04ff02ffff04ff17ff80808080ff80808080ffff02ff17ff2f808080ff0180ffff04ffff01ff32ff02ffff03ffff07ff0580ffff01ff0bffff0102ffff02ff06ffff04ff02ffff04ff09ff80808080ffff02ff06ffff04ff02ffff04ff0dff8080808080ffff01ff0bffff0101ff058080ff0180ff018080ffff04ffff01b0a9cc8198f9453fa1945c74a45a32037aa42b406896d966118cab49786938d7082bd13a61d36fc24208f9fc491baffd01ff018080", "solution": "0xff80ffff01ffff33ffa0ad68ab3e03965c749f54cca2d979c84e1404516fc8452da3c1763af88a21b7eaff0a80ffff33ffa02c68ed218bd3dc011237a1b79c2669905f763dc3f1ed4fea6ddb1760d12edb78ff4680ffff3cffa0ac13b120c3aa90755a5cd41f3679c694ce0a2fa6d691309dbd6a30863d6ca27a8080ff8080"}], "aggregated_signature": "0xb166d38fa4b84eccd3da941264c8b3c6decd97f3ad020bc48b8e51d2013a7e200f3ac54c676e9eb6786a6cc29c3bcc070040c5c5fd608295461e9c918959b608e1a9b9a0911aac425d22d8068fbd538e5468692eccfbaaf36a5f9267d3170c22"}`)
	bundle := new(SpendBundle)
	if err := json.Unmarshal(reqBytes, bundle); err != nil {
		t.Fatal(err)
	}
	req := PushTxReq{
		SpendBundle: *bundle,
	}
	body, err := testClient.PushTx(req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
}

func TestSignTx(t *testing.T) {

}

func TestSignTxFunc(t *testing.T) {
	var (
		// 签名的walletSk
		skHexString = "58a8b3237c9981ff476a897fc0d6b377bd5b2e57cbfcdf664c76963a52041012"
		// 待签名的msg
		msgHex = "10f4962dfabb2e21217ae886084a10a8626e873d692c353b9004331d0e9966e33445218ca583311ea1490b1a8cdf2af8ad84d583adb31c2cfa141bace8cc9fc3ccd5bb71183532bff220ba46c268991a3ff07eb358e8255a65c30a2dce0e5fbb"
		// 签名的alletSk对应的pk
		pkHex = "a9cc8198f9453fa1945c74a45a32037aa42b406896d966118cab49786938d7082bd13a61d36fc24208f9fc491baffd01"
	)

	sk, _ := bls.KeyFromHexString(skHexString)

	msgBytes, _ := hex.DecodeString(msgHex)
	msgList := [][]byte{msgBytes}

	pkBytes, _ := hex.DecodeString(pkHex)
	pkList := [][]byte{pkBytes}

	signBytes, err := testClient.signTx(sk, msgList, pkList)
	if err != nil {
		t.Fatal(err)
	}
	// should be 0xb166d38fa4b84eccd3da941264c8b3c6decd97f3ad020bc48b8e51d2013a7e200f3ac54c676e9eb6786a6cc29c3bcc070040c5c5fd608295461e9c918959b608e1a9b9a0911aac425d22d8068fbd538e5468692eccfbaaf36a5f9267d3170c22
	t.Log(hex.EncodeToString(signBytes))
}
