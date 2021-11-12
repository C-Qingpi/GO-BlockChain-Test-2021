package blockchain

import (
	"fmt"
	"github.com/dgraph-io/badger"
	// This is our database import
)

const dbPath = "./tmp/blocks"

// maybe

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
	// instead of points to blocks, this points to the database
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func InitBlockChain() *BlockChain {
	// 1. open connection to database
	// 2. check if there is a blockchain
	// 3. grab if there is; create if there is not

	var lastHash []byte
	opts := badger.DefaultOptions(dbPath)
	db, err := badger.Open(opts)
	Handle(err)
	//part 1 finished
	err = db.Update(func(txn *badger.Txn) error {
		// "lh" stand for last hash, it stores the latest block in the blockchain database. Could be used to quickly verify the state of the current blockchain
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")
			genesis := Genesis()
			fmt.Println("Genesis Allowed")
			err = txn.Set(genesis.Hash, genesis.Serialize())
			Handle(err)
			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash

			return err
			//part 2/3 finished
		} else {
			item, err := txn.Get([]byte("lh"))
			Handle(err)
			err = item.Value(func(val []byte) error {
				lastHash = val
				return nil
			})
			Handle(err)
			return err
		}
	})
	Handle(err)

	blockchain := BlockChain{lastHash, db}
	return &blockchain
	//part 3/3 finished
}

func (chain *BlockChain) AddBlock(data string) {
	var lastHash []byte
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		Handle(err)
		return err
		// return error to the viewing action unless correctly read the last hash field
	})
	Handle(err)
	// now knowing the lastHash

	newBlock := CreateBlock(data, lastHash)
	// (data, prevHash) as arg; now the computer
	// found a new block, it's trying to add it to the database
	// in the databse, blocks are stored as
	//  [Hash]-[Serialized Block Bytes] pairs

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		// update the last hash since a new block came in.
		err = txn.Set([]byte("lh"), newBlock.Hash)

		// now update the chain in the executing memory.
		chain.LastHash = newBlock.Hash
		return err
	})
	Handle(err)

}

// takes a blockchain and return to a function that looks over
// the databse one by one back from the latest block.
func (chain *BlockChain) Iterator() *BlockChainIterator {
	iterator := BlockChainIterator{chain.LastHash, chain.Database}
	// in formal application this could be a problem **,
	// as one synchronizing the database, the block could grew large.
	// at the mean time
	return &iterator
}

func (iterator *BlockChainIterator) Next() *Block {
	var block *Block

	err := iterator.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iterator.CurrentHash)
		Handle(err)

		err = item.Value(func(val []byte) error {
			block = Deserialize(val)
			return nil
		})
		Handle(err)
		return err
	})
	Handle(err)

	// after updating the database, also update the variable
	// in executing memory
	iterator.CurrentHash = block.PrevHash

	return block
}
