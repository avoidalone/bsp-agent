// Package proof contains all the functions to make a proof on the proofchain about a block replica
package proof

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"time"

	"github.com/covalenthq/bsp-agent/internal/config"
	ty "github.com/covalenthq/bsp-agent/internal/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
)

const (
	proofTxTimeout uint64 = 60
)

// SendBlockReplicaProofTx calls the proof-chain contract to make a transaction for the block-replica that it is processing
func SendBlockReplicaProofTx(ctx context.Context, config *config.EthConfig, proofChain string, ethClient *ethclient.Client, chainHeight uint64, chainLen uint64, resultSegment []byte, replicaURL string, blockReplica *ty.BlockReplica, txHash chan string) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(proofTxTimeout))
	defer cancel()

	_, opts, _, err := getTransactionOpts(ctx, config, ethClient)
	if err != nil {
		log.Error("error getting transaction ops: ", err.Error())
		txHash <- ""

		return
	}

	contractAddress := common.HexToAddress(proofChain)
	contract, err := NewProofChain(contractAddress, ethClient)
	if err != nil {
		log.Error("error binding to deployed contract: ", err.Error())
		txHash <- ""

		return
	}

	jsonResult, err := json.Marshal(resultSegment)
	if err != nil {
		log.Error("error in JSON marshaling result segment: ", err.Error())
		txHash <- ""

		return
	}
	sha256Result := sha256.Sum256(jsonResult)

	transaction, err := contract.SubmitBlockSpecimenProof(opts, blockReplica.NetworkId, chainHeight, blockReplica.Hash, sha256Result, replicaURL)

	if err != nil {
		log.Error("error calling deployed contract: ", err)
		txHash <- ""

		return
	}
	receipt, err := bind.WaitMined(ctx, ethClient, transaction)
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Error("block-result proof tx call: ", transaction.Hash(), " to proof contract failed: ", err.Error())
		txHash <- ""

		return
	}
	if err != nil {
		log.Error("error in waiting for tx to be mined on the blockchain: ", err.Error())
		txHash <- ""

		return
	}

	txHash <- receipt.TxHash.String()
}

func getTransactionOpts(ctx context.Context, config *config.EthConfig, ethClient *ethclient.Client) (common.Address, *bind.TransactOpts, uint64, error) {
	sKey := config.PrivateKey
	chainID, err := ethClient.ChainID(ctx)
	if err != nil {
		log.Error("error in getting transaction options: ", err.Error())
	}
	secretKey := crypto.ToECDSAUnsafe(common.FromHex(sKey))
	addr := crypto.PubkeyToAddress(secretKey.PublicKey)
	opts, err := bind.NewKeyedTransactorWithChainID(secretKey, chainID)
	if err != nil {
		log.Fatalf("error getting new keyed transactor with chain id: %v", err)
	}

	return addr, opts, chainID.Uint64(), err
}
