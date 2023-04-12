package main

import (
	"fmt"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/coinbase/waas-client-library-go/auth"
	"github.com/coinbase/waas-client-library-go/clients"
	v1clients "github.com/coinbase/waas-client-library-go/clients/v1"
	v1types "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/types/v1"
	blockchain "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/blockchain/v1"
	protocols "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/protocols/v1"
	mpcKeys "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_keys/v1"
	mpcTransactions "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_transactions/v1"
	mpcWallet "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/mpc_wallets/v1"
	pools "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/pools/v1"
	"google.golang.org/api/iterator"
)

const (
	version = "v1"

	// apiKeyName is the name of the API Key to use. Fill this out before running the main function.
	apiKeyName = "<YOUR_API_KEY_NAME>"

	// apiKeyPrivateKey is the private key of the API Key to use. Fill this out before running the main function.
	apiKeyPrivateKey = "<YOUR_PRIVATE_KEY>"
)

func parseInt32(str string) (int32, error) {
	i, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("cannot parse %q as int32: %v", str, err)
	}
	return int32(i), nil
}

// An example function to demonstrate how to use the WaaS client libraries.
func main() {

	ctx := context.Background()

	authOpt := clients.WithAPIKey(&auth.APIKey{
		Name:       apiKeyName,
		PrivateKey: apiKeyPrivateKey,
	})

	// Create BlockchainServiceClient
	blockchainClient, err := v1clients.NewBlockchainServiceClient(ctx, authOpt)
	if err != nil {
		log.Fatalf("Error instantiating BlockchainServiceClient: %v", err)
	}

	// Create MPCKeyServiceClient
	mpcKeyClient, err := v1clients.NewMPCKeyServiceClient(ctx, authOpt)
	if err != nil {
		log.Fatalf("Error instantiating MPCKeyServiceClient: %v", err)
	}

	// Create MPCTransactionServiceClient
	mpcTransactionClient, err := v1clients.NewMPCTransactionServiceClient(ctx, authOpt)
	if err != nil {
		log.Fatalf("Error instantiating MPCTransactionServiceClient: %v", err)
	}

	// Create MPCWalletServiceClient
	mpcWalletClient, err := v1clients.NewMPCWalletServiceClient(ctx, authOpt)
	if err != nil {
		log.Fatalf("Error instantiating MPCWalletServiceClient: %v", err)
	}

	// Create PoolServiceClient
	poolClient, err := v1clients.NewPoolServiceClient(ctx, authOpt)
	if err != nil {
		log.Fatalf("Error instantiating PoolServiceClient: %v", err)
	}

	// Create ProtocolServiceClient
	protocolClient, err := v1clients.NewProtocolServiceClient(ctx, authOpt)
	if err != nil {
		log.Fatalf("Error instantiating ProtocolServiceClient: %v", err)
	}

	// Create a Gin router
	router := gin.Default()

	// Blockchain API - ListNetworks (GET)
	router.GET("/blockchain/v1/networks", func(c *gin.Context) {
	
		pageSize, err := parseInt32(c.DefaultQuery("pageSize", "50"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		pageToken := c.Query("pageToken")

		networksIter := blockchainClient.ListNetworks(context.Background(), &blockchain.ListNetworksRequest{PageSize: pageSize, PageToken: pageToken})

		var networks []*blockchain.Network
		for {
			network, err := networksIter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			networks = append(networks, network)
		}

		networksJSON, err := json.Marshal(networks)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "application/json")
		c.Writer.Write(networksJSON)
	})

	// Blockchain API - GetNetwork (GET)
	router.GET("/blockchain/v1/networks/:networkId", func(c *gin.Context) {
		networkId := c.Param("networkId")
		networkName := "networks/" + networkId

		network, err := blockchainClient.GetNetwork(context.Background(), &blockchain.GetNetworkRequest{Name: networkName})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		networksJSON, err := json.Marshal(network)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "application/json")
		c.Writer.Write(networksJSON)
	})

	// Blockchain API - ListAssets (GET)
	router.GET("/blockchain/v1/networks/:networkId/assets", func(c *gin.Context) {
		networkId := c.Param("networkId")
		networkName := "networks/" + networkId

		pageSize, err := parseInt32(c.DefaultQuery("pageSize", "50"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		pageToken := c.Query("pageToken")
		filter := c.Query("filter")

		assetsIter := blockchainClient.ListAssets(context.Background(), &blockchain.ListAssetsRequest{Parent: networkName, PageSize: pageSize, PageToken: pageToken, Filter: filter})

		var assets []*blockchain.Asset
		for {
			asset, err := assetsIter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			assets = append(assets, asset)
		}

		assetsJSON, err := json.Marshal(assets)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "application/json")
		c.Writer.Write(assetsJSON)
	})

	// Blockchain API - GetAsset (GET)
	router.GET("/blockchain/v1/networks/:networkId/assets/:assetId", func(c *gin.Context) {
		networkId := c.Param("networkId")
		assetId := c.Param("assetId")
		var assetName = "networks/" + networkId + "/assets/" + assetId

		asset, err := blockchainClient.GetAsset(context.Background(), &blockchain.GetAssetRequest{Name: assetName})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		assetJSON, err := json.Marshal(asset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "application/json")
		c.Writer.Write(assetJSON)
	})

	// MPC Keys API - GetMPCKey (GET)
	router.GET("/mpc_keys/v1/pools/:poolId/deviceGroups/:deviceGroupId/mpcKeys/:mpcKeyId", func(c *gin.Context) {
		poolId := c.Param("poolId")
		deviceGroupId := c.Param("deviceGroupId")
		mpcKeyId := c.Param("mpcKeyId")
		mpcKeyName := "pools/" + poolId + "/deviceGroups/" + deviceGroupId + "/mpcKeys/" + mpcKeyId

		mpcKey, err := mpcKeyClient.GetMPCKey(context.Background(), &mpcKeys.GetMPCKeyRequest{Name: mpcKeyName})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		mpcKeyJSON, err := json.Marshal(mpcKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "application/json")
		c.Writer.Write(mpcKeyJSON)
	})

	// MPC Keys API - GetDevice (GET)
	router.GET("/mpc_keys/v1/devices/:deviceId", func(c *gin.Context) {
		deviceId := c.Param("deviceId")
		deviceName := "devices/" + deviceId

		device, err := mpcKeyClient.GetDevice(context.Background(), &mpcKeys.GetDeviceRequest{Name: deviceName})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		deviceJSON, err := json.Marshal(device)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "application/json")
		c.Writer.Write(deviceJSON)
	})

	// MPC Keys API - GetDeviceGroup (GET)
	router.GET("/mpc_keys/v1/pools/:poolId/deviceGroups/:deviceGroupId", func(c *gin.Context) {
		poolId := c.Param("poolId")
		deviceGroupId := c.Param("deviceGroupId")
		deviceGroupName := "pools/" + poolId + "/deviceGroups/" + deviceGroupId

		deviceGroup, err := mpcKeyClient.GetDeviceGroup(context.Background(), &mpcKeys.GetDeviceGroupRequest{Name: deviceGroupName})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		deviceGroupJSON, err := json.Marshal(deviceGroup)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "application/json")
		c.Writer.Write(deviceGroupJSON)
	})

	// MPC Keys API - ListMPCOperations (GET)
	router.GET("/mpc_keys/v1/pools/:poolId/deviceGroups/:deviceGroupId/mpcOperations", func(c *gin.Context) {
		poolId := c.Param("poolId")
		deviceGroupId := c.Param("deviceGroupId")
		deviceGroupName := "pools/" + poolId + "/deviceGroups/" + deviceGroupId

		mpcOperations, err := mpcKeyClient.ListMPCOperations(context.Background(), &mpcKeys.ListMPCOperationsRequest{Parent: deviceGroupName})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		mpcOperationsJSON, err := json.Marshal(mpcOperations)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "application/json")
		c.Writer.Write(mpcOperationsJSON)
	})

	// MPC Keys API - RegisterDevice (POST)
	router.POST("/mpc_keys/v1/device/register", func(c *gin.Context) {
		var registerDeviceReq *mpcKeys.RegisterDeviceRequest
		if err := c.BindJSON(&registerDeviceReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response, err := mpcKeyClient.RegisterDevice(ctx, registerDeviceReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		registerDeviceJSON, err := json.Marshal(response)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "application/json")
		c.Writer.Write(registerDeviceJSON)
	})

	// MPC Keys API - CreateMPCKey (POST)
	router.POST("/mpc_keys/v1/pools/:poolId/deviceGroups/:deviceGroupId/mpcKeys", func(c *gin.Context) {
		poolId := c.Param("poolId")
		deviceGroupId := c.Param("deviceGroupId")
		deviceGroupName := "pools/" + poolId + "/deviceGroups/" + deviceGroupId

		requestId := c.Query("requestId")

		var mpcKey *mpcKeys.MPCKey
		if err := c.BindJSON(&mpcKey); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		createMpcKeyReq := &mpcKeys.CreateMPCKeyRequest{Parent: deviceGroupName, MpcKey: mpcKey, RequestId: requestId}
		response, err := mpcKeyClient.CreateMPCKey(ctx, createMpcKeyReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		mpcKeyJSON, err := json.Marshal(response)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "application/json")
		c.Writer.Write(mpcKeyJSON)
	})

	// MPC Keys API - CreateSignature (POST)
	router.POST("/mpc_keys/v1/pools/:poolId/deviceGroups/:deviceGroupId/mpcKeys/:mpcKeyId/signatures", func(c *gin.Context) {
		poolId := c.Param("poolId")
		deviceGroupId := c.Param("deviceGroupId")
		mpcKeyId := c.Param("mpcKeyId")
		mpcKeyName := "pools/" + poolId + "/deviceGroups/" + deviceGroupId + "/mpcKeys/" + mpcKeyId

		requestId := c.Query("requestId")

		var signature *mpcKeys.Signature
		if err := c.BindJSON(&signature); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		createSignatureReq := &mpcKeys.CreateSignatureRequest{Parent: mpcKeyName, Signature: signature, RequestId: requestId}
		response, err := mpcKeyClient.CreateSignature(ctx, createSignatureReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		signatureJSON, err := json.Marshal(response)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "application/json")
		c.Writer.Write(signatureJSON)
	})

	// MPC Keys API - CreateDeviceGroup (POST)
	router.POST("/mpc_keys/v1/pools/:poolId/deviceGroups", func(c *gin.Context) {
		poolId := c.Param("poolId")
		mpcKeyName := "pools/" + poolId

		deviceGroupId := c.Query("deviceGroupId")
		requestId := c.Query("requestId")

		var deviceGroup *mpcKeys.DeviceGroup
		if err := c.BindJSON(&deviceGroup); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		createDeviceGroupReq := &mpcKeys.CreateDeviceGroupRequest{Parent: mpcKeyName, DeviceGroup: deviceGroup, DeviceGroupId: deviceGroupId, RequestId: requestId}
		response, err := mpcKeyClient.CreateDeviceGroup(ctx, createDeviceGroupReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		deviceGroupJSON, err := json.Marshal(response)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "application/json")
		c.Writer.Write(deviceGroupJSON)
	})

	// MPC Transactions API - GetMPCTransaction (GET)
	router.GET("/mpc_transactions/v1/pools/:poolId/mpcWallets/:mpcWalletId/mpcTransactions/:mpcTransactionId", func(c *gin.Context) {
		poolId := c.Param("poolId")
		mpcWalletId := c.Param("mpcWalletId")
		mpcTransactionId := c.Param("mpcTransactionId")

		mpcTransactionName := "pools/" + poolId + "/mpcWallets/" + mpcWalletId + "/mpcTransactions/" + mpcTransactionId

		mpcTx, err := mpcTransactionClient.GetMPCTransaction(context.Background(), &mpcTransactions.GetMPCTransactionRequest{Name: mpcTransactionName})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		mpcTxJSON, err := json.Marshal(mpcTx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "application/json")
		c.Writer.Write(mpcTxJSON)
	})

	// MPC Transactions API - ListMPCTransactions (GET)
	router.GET("/mpc_transactions/v1/pools/:poolId/mpcWallets/:mpcWalletId/mpcTransactions", func(c *gin.Context) {
		poolId := c.Param("poolId")
		mpcWalletId := c.Param("mpcWalletId")
		mpcWalletName := "pools/" + poolId + "/mpcWallets/" + mpcWalletId

		pageSize, err := parseInt32(c.DefaultQuery("pageSize", "50"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		pageToken := c.Query("pageToken")

		mpxTxsIter := mpcTransactionClient.ListMPCTransactions(context.Background(), &mpcTransactions.ListMPCTransactionsRequest{Parent: mpcWalletName, PageSize: pageSize, PageToken: pageToken})

		var mpxTxs []*mpcTransactions.MPCTransaction
		for {
			mpxTx, err := mpxTxsIter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			mpxTxs = append(mpxTxs, mpxTx)
		}

		mpxTxsJSON, err := json.Marshal(mpxTxs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "application/json")
		c.Writer.Write(mpxTxsJSON)
	})

	// MPC Transactions API - CreateMPCTransaction (POST)
	router.POST("/mpc_transactions/v1/pools/:poolId/mpcWallets/:mpcWalletId/mpcTransactions", func(c *gin.Context) {
		poolId := c.Param("poolId")
		mpcWalletId := c.Param("mpcWalletId")
		mpcWalletName := "pools/" + poolId + "/mpcWallets/" + mpcWalletId

		var requestBody mpcTransactions.CreateMPCTransactionRequest
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		createMpcTxReq := &mpcTransactions.CreateMPCTransactionRequest{Parent: mpcWalletName, MpcTransaction: requestBody.MpcTransaction, Input: requestBody.Input, OverrideNonce: requestBody.OverrideNonce, RequestId: requestBody.RequestId}
		response, err := mpcTransactionClient.CreateMPCTransaction(ctx, createMpcTxReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		mpcTxJSON, err := json.Marshal(response)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "application/json")
		c.Writer.Write(mpcTxJSON)
	})

	// MPC Wallets API - GetMPCWallet (GET)
	router.GET("/mpc_wallets/v1/pools/:poolId/mpcWallets/:mpcWalletId", func(c *gin.Context) {
		poolId := c.Param("poolId")
		mpcWalletId := c.Param("mpcWalletId")

		mpcWalletName := "pools/" + poolId + "/mpcWallets/" + mpcWalletId
	
		wallet, err := mpcWalletClient.GetMPCWallet(context.Background(), &mpcWallet.GetMPCWalletRequest{Name: mpcWalletName})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		walletJSON, err := json.Marshal(wallet)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "application/json")
		c.Writer.Write(walletJSON)
	})

	// MPC Wallets API - ListMPCWallets (GET)
	router.GET("/mpc_wallets/v1/pools/:poolId/mpcWallets", func(c *gin.Context) {
		poolId := c.Param("poolId")
		poolName := "pools/" + poolId

		pageSize, err := parseInt32(c.DefaultQuery("pageSize", "50"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		pageToken := c.Query("pageToken")

		walletsIter := mpcWalletClient.ListMPCWallets(context.Background(), &mpcWallet.ListMPCWalletsRequest{Parent: poolName, PageSize: pageSize, PageToken: pageToken})

		var wallets []*mpcWallet.MPCWallet
		for {
			wallet, err := walletsIter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			wallets = append(wallets, wallet)
		}

		walletsJSON, err := json.Marshal(wallets)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "application/json")
		c.Writer.Write(walletsJSON)
	})

	// MPC Wallets API - GetAddress (GET)
	router.GET("/mpc_wallets/v1/networks/:networkId/addresses/:addressId", func(c *gin.Context) {
		networkId := c.Param("networkId")
		addressId := c.Param("addressId")
		networkName := "networks/" + networkId + "/addresses/" + addressId
	
		address, err := mpcWalletClient.GetAddress(context.Background(), &mpcWallet.GetAddressRequest{Name: networkName})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		addressJSON, err := json.Marshal(address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "application/json")
		c.Writer.Write(addressJSON)
	})

	// MPC Wallets API - ListAddresses (GET)
	router.GET("/mpc_wallets/v1/networks/:networkId/addresses", func(c *gin.Context) {
		networkId := c.Param("networkId")
		networkName := "networks/" + networkId

		pageSize, err := parseInt32(c.DefaultQuery("pageSize", "50"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		pageToken := c.Query("pageToken")
		wallet := c.Query("mpcWallet")

		addressesIter := mpcWalletClient.ListAddresses(context.Background(), &mpcWallet.ListAddressesRequest{Parent: networkName, MpcWallet: wallet, PageSize: pageSize, PageToken: pageToken})

		var addresses []*mpcWallet.Address
		for {
			address, err := addressesIter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			addresses = append(addresses, address)
		}

		addressesJSON, err := json.Marshal(addresses)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "application/json")
		c.Writer.Write(addressesJSON)
	})

	// MPC Wallets API - ListBalances (GET)
	router.GET("/mpc_wallets/v1/networks/:networkId/addresses/:addressId/balances", func(c *gin.Context) {
		networkId := c.Param("networkId")
		addressId := c.Param("addressId")
		addressName := "networks/" + networkId + "/addresses/" + addressId

		pageSize, err := parseInt32(c.DefaultQuery("pageSize", "50"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		pageToken := c.Query("pageToken")

		balancesIter := mpcWalletClient.ListBalances(context.Background(), &mpcWallet.ListBalancesRequest{Parent: addressName, PageSize: pageSize, PageToken: pageToken})

		var balances []*mpcWallet.Balance
		for {
			balance, err := balancesIter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			balances = append(balances, balance)
		}

		balancesJSON, err := json.Marshal(balances)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "application/json")
		c.Writer.Write(balancesJSON)
	})

	// MPC Wallets API - CreateMPCWallet (POST)
	router.POST("/mpc_wallets/v1/pools/:poolId/mpcWallets", func(c *gin.Context) {
		poolId := c.Param("poolId")
		poolName := "pools/" + poolId

		device := c.Query("device")
		requestId := c.Query("requestId")

		var wallet *mpcWallet.MPCWallet
		if err := c.BindJSON(&wallet); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		createMpcWalletReq := &mpcWallet.CreateMPCWalletRequest{Parent: poolName, MpcWallet: wallet, Device: device, RequestId: requestId}
		response, err := mpcWalletClient.CreateMPCWallet(ctx, createMpcWalletReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		mpcWalletJSON, err := json.Marshal(response)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "application/json")
		c.Writer.Write(mpcWalletJSON)
	})

	// MPC Wallets API - GenerateAddress (POST)
	router.POST("/mpc_wallets/v1/pools/:poolId/mpcWallets/:mpcWalletId/generateAddress", func(c *gin.Context) {
		poolId := c.Param("poolId")
		mpcWalletId := c.Param("mpcWalletId")
		mpcWalletName := "pools/" + poolId + "/mpcWallets/" + mpcWalletId

		var requestBody mpcWallet.GenerateAddressRequest
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		generateAddressReq := &mpcWallet.GenerateAddressRequest{MpcWallet: mpcWalletName, Network: requestBody.Network, RequestId: requestBody.RequestId}
		response, err := mpcWalletClient.GenerateAddress(ctx, generateAddressReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		addressJSON, err := json.Marshal(response)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "application/json")
		c.Writer.Write(addressJSON)
	})

	// Pools API - GetPool (GET)
	router.GET("/pools/v1/pools/:poolId", func(c *gin.Context) {
		poolId := c.Param("poolId")
		poolName := "pools/" + poolId

		pool, err := poolClient.GetPool(context.Background(), &pools.GetPoolRequest{Name: poolName})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		poolJSON, err := json.Marshal(pool)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "application/json")
		c.Writer.Write(poolJSON)
	})

	// Pools API - ListPools (GET)
	router.GET("/pools/v1/pools", func(c *gin.Context) {

		pageSize, err := parseInt32(c.DefaultQuery("pageSize", "50"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		pageToken := c.Query("pageToken")
	
		poolsIter := poolClient.ListPools(context.Background(), &pools.ListPoolsRequest{PageSize: pageSize, PageToken: pageToken})

		var pools []*pools.Pool
		for {
			pool, err := poolsIter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			pools = append(pools, pool)
		}

		poolsJSON, err := json.Marshal(pools)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "application/json")
		c.Writer.Write(poolsJSON)
	})

	// Pools API - CreatePool (POST)
	router.POST("/pools/v1/pools", func(c *gin.Context) {
		poolId := c.Query("poolId")

		var pool *pools.Pool
		if err := c.BindJSON(&pool); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		createPoolReq := &pools.CreatePoolRequest{PoolId: poolId, Pool: pool}
		response, err := poolClient.CreatePool(ctx, createPoolReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		poolJSON, err := json.Marshal(response)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "application/json")
		c.Writer.Write(poolJSON)
	})

	// Protocols API - BroadcastTransaction (POST)
	router.POST("/protocols/v1/networks/:networkId/broadcastTransaction", func(c *gin.Context) {
		networkId := c.Param("networkId")
		networkName := "networks/" + networkId
	
		var transaction *v1types.Transaction
		if err := c.BindJSON(&transaction); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	
		broadcastTxReq := &protocols.BroadcastTransactionRequest{Network: networkName, Transaction: transaction}
		response, err := protocolClient.BroadcastTransaction(ctx, broadcastTxReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		txJSON, err := json.Marshal(response)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "application/json")
		c.Writer.Write(txJSON)
	})

	// Protocols API - ConstructTransaction (POST)
	router.POST("/protocols/v1/networks/:networkId/constructTransaction", func(c *gin.Context) {
		networkId := c.Param("networkId")
		networkName := "networks/" + networkId

		var input *v1types.TransactionInput
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		constructTxReq := &protocols.ConstructTransactionRequest{Network: networkName, Input: input}
		response, err := protocolClient.ConstructTransaction(ctx, constructTxReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		txJSON, err := json.Marshal(response)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "application/json")
		c.Writer.Write(txJSON)
	})

	// Protocols API - ConstructTransferTransaction (POST)
	router.POST("/protocols/v1/networks/:networkId/constructTransferTransaction", func(c *gin.Context) {
		networkId := c.Param("networkId")
		networkName := "networks/" + networkId

		var requestBody protocols.ConstructTransferTransactionRequest
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		constructTransferTxReq := &protocols.ConstructTransferTransactionRequest{Network: networkName, Asset: requestBody.Asset, Sender: requestBody.Sender, Recipient: requestBody.Recipient, Amount: requestBody.Amount, Nonce: requestBody.Nonce, Fee: requestBody.Fee}
		response, err := protocolClient.ConstructTransferTransaction(ctx, constructTransferTxReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		transferTxJSON, err := json.Marshal(response)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "application/json")
		c.Writer.Write(transferTxJSON)
	})

	server := http.Server{
		Addr: ":8080",
		Handler: router,
	}

	log.Println("Proxy server listening on port 8080...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
