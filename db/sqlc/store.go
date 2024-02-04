package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v5/pgxpool"
)

// store provides all functions to execute db queries amd trasactions

type SQLStore struct {
    *Queries
    connPool *pgxpool.Pool
}




type Store struct {
	*Queries 
	db *sql.DB
	

}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:db,
		Queries:New(db),
	
	}
}

// execTx executes a function within a db trasaction 
func (store * Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return nil
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rberr :=	tx.Rollback(); rberr != nil{
			return fmt.Errorf("tx err: %v, rb err: %v", err, rberr)
		}
		return err
	}

	return tx.Commit()

}



// TransferTxParams contains the input parameters of the transfer transaction 
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID int64 `json:"to_account_id"`
	Amount int64 `json:"amount"`

}

// transferTxResult in the result of the transfer transcation
type TransferTxResult struct {
	Transfer Transfer `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount 	Account `json:"to_account"`
	FromEntry   Entry  `json:"from_entry"`
	ToEntry   Entry    `json:"to_entry"`

}


var txKey = struct {}{}

// Transfeer performs a money transfrer from one account to the other
// it creates a transfer record, add account entires, and update accounts'balance
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, err){
	var result TransferTxResult 

	err := store.execTx(ctx, func(q *Queries) error {

		var err error
		txName := ctx.Value(txKey)

		fmt.Println(txName,"Create transfer")

		// create a transfer record 
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,

		})

		if err != nil {
			return err
		}


		fmt.Println(txName,"Create entry1")
		// add account entires fromEntry and toEntry
		result.FromEntry, err= q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,

		})
		if err != nil {
			return err
		}	

		fmt.Println(txName,"Create entry2")
		result.ToEntry, err= q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,

		})
		if err != nil {
			return err
		}	


		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:  arg.FromAccountID,
				Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		
		// fmt.Println(txName,"update account 2 ")
		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID: arg.ToAccountID,
				Amount: arg.Amount,
		})

		if err != nil {
			return err
		}

		return nil 
	})

	
	
	return result, err
}