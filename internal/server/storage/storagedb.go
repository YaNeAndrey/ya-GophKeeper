package storage

import (
	"context"
	"database/sql"
	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
	log "github.com/sirupsen/logrus"
	"time"
	"ya-GophKeeper/internal/content"
	"ya-GophKeeper/internal/server/srverror"
)

type StorageDB struct {
	connectionString string
}

func InitStorageDB(connectionString string) (*StorageDB, error) {
	db, err := TryToOpenDBConnection(connectionString)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	myContext := context.TODO()

	_, err = db.ExecContext(myContext, "CREATE TABLE IF NOT EXISTS Users(ID_User serial PRIMARY KEY, login VARCHAR(256) NOT NULL, passwd TEXT NOT NULL);")
	if err != nil {
		return nil, err
	}
	_, err = db.ExecContext(myContext, "CREATE TABLE IF NOT EXISTS Credentials(ID_Credential serial PRIMARY KEY, resource VARCHAR(256), login VARCHAR(256), passwd TEXT, modification_time timestamp);")
	if err != nil {
		return nil, err
	}
	_, err = db.ExecContext(myContext, "CREATE TABLE IF NOT EXISTS Users_Credentials(ID_User_Credentials serial PRIMARY KEY, ID_User int REFERENCES Users(ID_User), ID_Credential int REFERENCES Credentials(ID_Credential));")
	if err != nil {
		return nil, err
	}
	_, err = db.ExecContext(myContext, "CREATE TABLE IF NOT EXISTS CreditCards( ID_CreditCard serial PRIMARY KEY, bank varchar(50), card_number varchar(16) NOT NULL, valid_thru timestamp, cvv varchar(3),  modification_time timestamp);")
	if err != nil {
		return nil, err
	}
	_, err = db.ExecContext(myContext, "CREATE TABLE IF NOT EXISTS Users_CreditCards(ID_User_CreditCard serial PRIMARY KEY, ID_User int REFERENCES Users(ID_User), ID_CreditCard int REFERENCES CreditCards(ID_CreditCard));")
	if err != nil {
		return nil, err
	}
	_, err = db.ExecContext(myContext, "CREATE TABLE IF NOT EXISTS Texts(ID_Text serial PRIMARY KEY, description TEXT, content TEXT NOT NULL, modification_time timestamp);")
	if err != nil {
		return nil, err
	}
	_, err = db.ExecContext(myContext, "CREATE TABLE IF NOT EXISTS Users_Texts(ID_User_Text serial PRIMARY KEY, ID_User int REFERENCES Users(ID_User), ID_Text int REFERENCES Texts(ID_Text));")
	if err != nil {
		return nil, err
	}
	_, err = db.ExecContext(myContext, "CREATE TABLE IF NOT EXISTS Files(ID_File serial PRIMARY KEY, description TEXT, file_name varchar(256) NOT NULL, file_path varchar(256) NOT NULL, file_size int NOT NULL, modification_time timestamp);")
	if err != nil {
		return nil, err
	}
	_, err = db.ExecContext(myContext, "CREATE TABLE IF NOT EXISTS Users_Files(ID_User_File serial PRIMARY KEY, ID_User int REFERENCES Users(ID_User), ID_File int REFERENCES Files(ID_File));")
	if err != nil {
		return nil, err
	}

	var resStorage StorageDB
	resStorage.connectionString = connectionString
	return &resStorage, nil
}

func (st *StorageDB) AddNewUser(ctx context.Context, login string, password string) error {
	db, err := TryToOpenDBConnection(st.connectionString)
	if err != nil {
		return err
	}
	var id int
	id, err = GetUserID(ctx, login, db)

	if err != nil {
		return err
	}
	if id != 0 {
		return srverror.ErrLoginAlreadyTaken
	}

	_, err = db.ExecContext(ctx, "INSERT INTO Users(login,passwd) values($1, crypt($2, gen_salt('bf')));", login, password)
	if err != nil {
		return err
	}

	return nil
}
func (st *StorageDB) CheckUserPassword(ctx context.Context, login string, password string) (bool, error) {
	db, err := TryToOpenDBConnection(st.connectionString)
	if err != nil {
		return false, err
	}
	var passwdOK bool
	err = db.QueryRowContext(ctx, "SELECT (case when (passwd = crypt($2, passwd)) then 'True' else 'False' end) as ok FROM Users WHERE login = $1", login, password).Scan(&passwdOK)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return false, srverror.ErrLoginNotFound
		default:
			return false, err
		}
	}

	return passwdOK, nil
}
func (st *StorageDB) ChangeUserPassword(ctx context.Context, login string, password string) error {
	db, err := TryToOpenDBConnection(st.connectionString)
	if err != nil {
		return err
	}

	res, err := db.ExecContext(ctx, "UPDATE Users SET passwd = $2 WHERE login = $1", login, password)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return srverror.ErrLoginNotFound
	}

	return nil
}

func (st *StorageDB) AddNewCreditCards(ctx context.Context, login string, creditCards []content.CreditCardInfo) error {
	db, err := TryToOpenDBConnection(st.connectionString)
	if err != nil {
		return err
	}
	var userID int
	userID, err = GetUserID(ctx, login, db)
	if err != nil {
		return err
	}
	if userID == 0 {
		return srverror.ErrLoginNotFound
	}
	for _, card := range creditCards {
		var cardID int
		err = db.QueryRowContext(ctx, "INSERT INTO CreditCards(bank,card_number,valid_thru,cvv,modification_time) values($1,$2,$3,$4,$5) RETURNING ID_CreditCard;", card.Bank, card.CardNumber, card.ValidThru, card.CVV, card.ModificationTime).Scan(&cardID)
		if err != nil {
			log.Println(err)
			continue
		}
		_, err = db.ExecContext(ctx, "INSERT INTO Users_CreditCards (ID_User,ID_CreditCard) values ($1,$2)", userID, cardID)
		if err != nil {
			log.Println(err)
			//mb remove creditcard?
		}
	}
	return nil
}

func (st *StorageDB) AddNewCredentials(ctx context.Context, login string, credentials []content.CredentialInfo) error {
	return nil
}
func (st *StorageDB) AddNewFiles(ctx context.Context, login string, files []content.BinaryFileInfo) error {
	return nil
}
func (st *StorageDB) AddNewTexts(ctx context.Context, login string, texts []content.TextInfo) error {
	return nil
}

func (st *StorageDB) RemoveCreditCards(ctx context.Context, login string, creditCardIDs []int) error {
	db, err := TryToOpenDBConnection(st.connectionString)
	if err != nil {
		return err
	}
	var userID int
	userID, err = GetUserID(ctx, login, db)
	if err != nil {
		return err
	}
	if userID == 0 {
		return srverror.ErrLoginNotFound
	}

	for _, cardID := range creditCardIDs {
		_, err = db.ExecContext(ctx, "DELETE FROM CreditCards where ID_CreditCard == $1", cardID)
		if err != nil {
			log.Println(err)
			continue
		}
		_, err = db.ExecContext(ctx, "DELETE FROM Users_CreditCards where ID_CreditCard == $1 and ID_User == $2", cardID, userID)
		if err != nil {
			log.Println(err)
			continue
		}
	}
	return nil
}
func (st *StorageDB) RemoveCredentials(ctx context.Context, login string, credentialIDs []int) error {
	db, err := TryToOpenDBConnection(st.connectionString)
	if err != nil {
		return err
	}
	var userID int
	userID, err = GetUserID(ctx, login, db)
	if err != nil {
		return err
	}
	if userID == 0 {
		return srverror.ErrLoginNotFound
	}

	for _, cardID := range credentialIDs {
		_, err = db.ExecContext(ctx, "DELETE FROM Credentials where ID_Credential == $1", cardID)
		if err != nil {
			log.Println(err)
			continue
		}
		_, err = db.ExecContext(ctx, "DELETE FROM Users_Credentials where ID_Credential == $1 and ID_User == $2", cardID, userID)
		if err != nil {
			log.Println(err)
			continue
		}
	}
	return nil
}
func (st *StorageDB) RemoveFiles(ctx context.Context, login string, fileIDs []int) error {
	return nil
}
func (st *StorageDB) RemoveTexts(ctx context.Context, login string, textIDs []int) error { return nil }

func (st *StorageDB) UpdateFiles(ctx context.Context, login string, files []content.BinaryFileInfo) error {
	return nil
}
func (st *StorageDB) UpdateTexts(ctx context.Context, login string, texts []content.TextInfo) error {
	db, err := TryToOpenDBConnection(st.connectionString)
	if err != nil {
		return err
	}

	for _, text := range texts {
		_, err = db.ExecContext(ctx, "UPDATE Texts SET description = $2, content = $3, modification_time = $4 where ID_Text == $1", text.ID, text.Description, text.Content, text.ModificationTime)
		if err != nil {
			log.Println(err)
			continue
		}
	}
	return nil
}
func (st *StorageDB) UpdateCreditCards(ctx context.Context, login string, creditCards []content.CreditCardInfo) error {
	db, err := TryToOpenDBConnection(st.connectionString)
	if err != nil {
		return err
	}

	for _, card := range creditCards {
		_, err = db.ExecContext(ctx, "UPDATE CreditCards SET bank = $2, card_number = $3, valid_thru = $4,cvv = $5, modification_time = $6 WHERE id_creditcard == $1", card.ID, card.Bank, card.CardNumber, card.ValidThru, card.CVV, card.ModificationTime)
		if err != nil {
			log.Println(err)
			continue
		}
	}
	return nil
}
func (st *StorageDB) UpdateCredentials(ctx context.Context, login string, credentials []content.CredentialInfo) error {
	db, err := TryToOpenDBConnection(st.connectionString)
	if err != nil {
		return err
	}

	for _, cred := range credentials {
		_, err = db.ExecContext(ctx, "UPDATE Credentials SET resource = $2, login = $3, passwd = $4,modification_time = $5 where id_credential == $1", cred.ID, cred.Resource, cred.Login, cred.Password, cred.ModificationTime)
		if err != nil {
			log.Println(err)
			continue
		}
	}
	return nil
}

func (st *StorageDB) GetCreditCards(ctx context.Context, login string, cardIDs []int) ([]content.CreditCardInfo, error) {
	return nil, nil
}
func (st *StorageDB) GetCredentials(ctx context.Context, login string, credIDs []int) ([]content.CredentialInfo, error) {
	return nil, nil
}
func (st *StorageDB) GetFiles(ctx context.Context, login string, fileIDs []int) ([]content.BinaryFileInfo, error) {
	return nil, nil
}
func (st *StorageDB) GetTexts(ctx context.Context, login string, textIDs []int) ([]content.TextInfo, error) {
	return nil, nil
}

func GetUserID(ctx context.Context, login string, db *sql.DB) (int, error) {
	var id int
	err := db.QueryRowContext(ctx, "SELECT Users.ID_User FROM Users WHERE login = $1", login).Scan(&id)
	if err != nil {
		return 0, err
	} else {
		return id, nil
	}
}

func TryToOpenDBConnection(dbConnectionString string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dbConnectionString)
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//var bufError error
	err = retry.Retry(

		func(attempt uint) error {
			if err = db.PingContext(ctx); err != nil {
				return err
			}

			return nil
		},
		strategy.Limit(4),
		strategy.Backoff(backoff.Incremental(-1*time.Second, 2*time.Second)),
	)

	if err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
