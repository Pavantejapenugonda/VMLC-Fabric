// CREDIT: https://hyperledger-fabric.readthedocs.io/en/release-2.2/chaincode4ade.html

package main

import (
  "encoding/json"
  "fmt"
  "log"
  "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
   type SmartContract struct {
      contractapi.Contract
    }

// Asset describes basic details of what makes up a simple asset
   type Asset struct {
	  ID			   string `json:"ID"`
      Make             string `json:"Make"`
      Model            string `json:"color"`
      Variant          int    `json:"size"`
      Owner            string `json:"owner"`
      AppraisedValue int    `json:"appraisedValue"`
    }

// InitLedger adds a base set of assets to the ledger
   func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
    assets := []Asset{
      {ID : 1,Make: "Hyundai", Model: "i10", Variant: "HTE", Owner: "Tomoko", AppraisedValue: 50000},
      {ID : 2,Make: "Kia", Model: "Seltos", Variant: "HTE", Owner: "Brad", AppraisedValue: 1000000},
      {ID : 3,Make: "TATA", Model: "Tiago", Variant: "SUV", Owner: "Jin Soo", AppraisedValue: 1200000},
      {ID : 4,Make: "Tyota", Model: "Innova", Variant: "HTE", Owner: "Max", AppraisedValue: 1500000},
      {ID : 5,Make: "Maruthi Suzuki", Model: "Swift", Variant: "HTE", Owner: "Adriana", AppraisedValue: 1000000},
      {ID : 6,Make: "BMW", Model: "X7", Variant: "HTE", Owner: "Michel", AppraisedValue: 10000000},
    }

    for _, asset := range assets {
      assetJSON, err := json.Marshal(asset)
      if err != nil {
        return err
      }

      err = ctx.GetStub().PutState(asset.ID, assetJSON)
      if err != nil {
        return fmt.Errorf("failed to put to world state. %v", err)
      }
    }

    return nil
  }

// CreateAsset issues a new asset to the world state with given details.
   func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size int, owner string, appraisedValue int) error {
    exists, err := s.AssetExists(ctx, id)
    if err != nil {
      return err
    }
    if exists {
      return fmt.Errorf("the asset %s already exists", id)
    }

    asset := Asset{
      ID:             id,
	  Make:			 make,
	  Model:         model,
	  Variant:       variant,
      Owner:          owner,
      AppraisedValue: appraisedValue,
    }
    assetJSON, err := json.Marshal(asset)
    if err != nil {
      return err
    }

    return ctx.GetStub().PutState(id, assetJSON)
  }

// ReadAsset returns the asset stored in the world state with given id.
   func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
    assetJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
      return nil, fmt.Errorf("failed to read from world state: %v", err)
    }
    if assetJSON == nil {
      return nil, fmt.Errorf("the asset %s does not exist", id)
    }

    var asset Asset
    err = json.Unmarshal(assetJSON, &asset)
    if err != nil {
      return nil, err
    }

    return &asset, nil
  }

// UpdateAsset updates an existing asset in the world state with provided parameters.
   func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, make string, model string, variant string, owner string, appraisedValue int) error {
    exists, err := s.AssetExists(ctx, id)
    if err != nil {
      return err
    }
    if !exists {
      return fmt.Errorf("the asset %s does not exist", id)
    }

    // overwriting original asset with new asset
    asset := Asset{
      ID:             id,
      Make:           make,
      Model:          model,
	  Variant:		  variant,
      Owner:          owner,
      AppraisedValue: appraisedValue,
    }
    assetJSON, err := json.Marshal(asset)
    if err != nil {
      return err
    }

    return ctx.GetStub().PutState(id, assetJSON)
  }

  // DeleteAsset deletes an given asset from the world state.
  func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
    exists, err := s.AssetExists(ctx, id)
    if err != nil {
      return err
    }
    if !exists {
      return fmt.Errorf("the asset %s does not exist", id)
    }

    return ctx.GetStub().DelState(id)
  }

// AssetExists returns true when asset with given ID exists in world state
   func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
    assetJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
      return false, fmt.Errorf("failed to read from world state: %v", err)
    }

    return assetJSON != nil, nil
  }

// TransferAsset updates the owner field of asset with given id in world state.
   func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {
    asset, err := s.ReadAsset(ctx, id)
    if err != nil {
      return err
    }

    asset.Owner = newOwner
    assetJSON, err := json.Marshal(asset)
    if err != nil {
      return err
    }

    return ctx.GetStub().PutState(id, assetJSON)
  }

// GetAllAssets returns all assets found in world state
   func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
// range query with empty string for startKey and endKey does an
// open-ended query of all assets in the chaincode namespace.
    resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
    if err != nil {
      return nil, err
    }
    defer resultsIterator.Close()

    var assets []*Asset
    for resultsIterator.HasNext() {
      queryResponse, err := resultsIterator.Next()
      if err != nil {
        return nil, err
      }

      var asset Asset
      err = json.Unmarshal(queryResponse.Value, &asset)
      if err != nil {
        return nil, err
      }
      assets = append(assets, &asset)
    }

    return assets, nil
  }

  func main() {
    assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
    if err != nil {
      log.Panicf("Error creating asset-transfer-basic chaincode: %v", err)
    }

    if err := assetChaincode.Start(); err != nil {
      log.Panicf("Error starting asset-transfer-basic chaincode: %v", err)
    }
  }