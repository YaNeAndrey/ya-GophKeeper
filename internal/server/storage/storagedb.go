package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"time"
	"ya-GophKeeper/internal/constants/srverror"
	"ya-GophKeeper/internal/constants/urlsuff"
	"ya-GophKeeper/internal/content"
)

type StorageDB struct {
	connectionString string
}

func InitStorageDB(connectionString string) *StorageDB {
	db, err := TryToOpenDBConnection(connectionString)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer db.Close()

	myContext := context.TODO()

	_, err = db.ExecContext(myContext, "CREATE TABLE IF NOT EXISTS Users(ID_User serial PRIMARY KEY, login VARCHAR(256) NOT NULL, passwd TEXT NOT NULL);")
	if err != nil {
		log.Println(err)
		return nil
	}
	_, err = db.ExecContext(myContext, "CREATE TABLE IF NOT EXISTS Credentials(ID_Credential serial PRIMARY KEY, resource VARCHAR(256), login VARCHAR(256), passwd TEXT, modification_time timestamp);")
	if err != nil {
		log.Println(err)
		return nil
	}
	_, err = db.ExecContext(myContext, "CREATE TABLE IF NOT EXISTS Users_Credentials(ID_User_Credentials serial PRIMARY KEY, ID_User int REFERENCES Users(ID_User), ID_Credential int REFERENCES Credentials(ID_Credential));")
	if err != nil {
		log.Println(err)
		return nil
	}
	_, err = db.ExecContext(myContext, "CREATE TABLE IF NOT EXISTS CreditCards( ID_CreditCard serial PRIMARY KEY, bank varchar(50), card_number varchar(16) NOT NULL, valid_thru timestamp, cvv varchar(3),  modification_time timestamp);")
	if err != nil {
		log.Println(err)
		return nil
	}
	_, err = db.ExecContext(myContext, "CREATE TABLE IF NOT EXISTS Users_CreditCards(ID_User_CreditCard serial PRIMARY KEY, ID_User int REFERENCES Users(ID_User), ID_CreditCard int REFERENCES CreditCards(ID_CreditCard));")
	if err != nil {
		log.Println(err)
		return nil
	}
	_, err = db.ExecContext(myContext, "CREATE TABLE IF NOT EXISTS Texts(ID_Text serial PRIMARY KEY, description TEXT, content TEXT NOT NULL, modification_time timestamp);")
	if err != nil {
		log.Println(err)
		return nil
	}
	_, err = db.ExecContext(myContext, "CREATE TABLE IF NOT EXISTS Users_Texts(ID_User_Text serial PRIMARY KEY, ID_User int REFERENCES Users(ID_User), ID_Text int REFERENCES Texts(ID_Text));")
	if err != nil {
		log.Println(err)
		return nil
	}
	_, err = db.ExecContext(myContext, "CREATE TABLE IF NOT EXISTS Files(ID_File serial PRIMARY KEY, description TEXT, file_name varchar(256) NOT NULL, file_path varchar(256) NOT NULL, file_size int NOT NULL, md5 uuid NOT NULL, modification_time timestamp);")
	if err != nil {
		log.Println(err)
		return nil
	}

	_, err = db.ExecContext(myContext, "CREATE TABLE IF NOT EXISTS Users_Files(ID_User_File serial PRIMARY KEY, ID_User int REFERENCES Users(ID_User), ID_File int REFERENCES Files(ID_File));")
	if err != nil {
		log.Println(err)
		return nil
	}

	var resStorage StorageDB
	resStorage.connectionString = connectionString
	return &resStorage
}

func (st *StorageDB) AddNewUser(ctx context.Context, login string, password string) error {
	db, err := TryToOpenDBConnection(st.connectionString)
	if err != nil {
		return err
	}
	_, err = GetUserID(ctx, login, db)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = db.ExecContext(ctx, "INSERT INTO Users(login,passwd) values($1, crypt($2, gen_salt('bf')));", login, password)
		if err != nil {
			return err
		}
		return nil
	} else if err != nil {
		return err
	}
	return srverror.ErrLoginAlreadyTaken
}
func (st *StorageDB) CheckUserPassword(ctx context.Context, login string, password string) (bool, error) {
	db, err := TryToOpenDBConnection(st.connectionString)
	if err != nil {
		return false, err
	}
	var passwdOK bool
	err = db.QueryRowContext(ctx, "SELECT (case when (passwd = crypt($2, passwd)) then 'True' else 'False' end) as ok FROM Users WHERE login = $1", login, password).Scan(&passwdOK)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
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

	res, err := db.ExecContext(ctx, "UPDATE Users SET passwd = crypt($2, gen_salt('bf')) WHERE login = $1", login, password)
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

func (st *StorageDB) AddNewCreditCards(ctx context.Context, login string, creditCards []content.CreditCardInfo) ([]content.CreditCardInfo, error) {
	db, err := TryToOpenDBConnection(st.connectionString)
	if err != nil {
		return nil, err
	}
	var userID int
	userID, err = GetUserID(ctx, login, db)
	if err != nil {
		return nil, err
	}
	if userID == 0 {
		return nil, srverror.ErrLoginNotFound
	}
	for index, card := range creditCards {
		var cardID int
		err = db.QueryRowContext(ctx, "INSERT INTO CreditCards(bank,card_number,valid_thru,cvv,modification_time) values($1,$2,$3,$4,$5) RETURNING ID_CreditCard;", card.Bank, card.CardNumber, card.ValidThru, card.CVV, card.ModificationTime).Scan(&cardID)

		if err != nil {
			log.Println(err)
			continue
		}
		creditCards[index].ID = cardID

		_, err = db.ExecContext(ctx, "INSERT INTO Users_CreditCards (ID_User,ID_CreditCard) values ($1,$2)", userID, cardID)
		if err != nil {
			log.Println(err)
			//mb remove creditcard?
		}
	}
	return creditCards, nil
}
func (st *StorageDB) AddNewCredentials(ctx context.Context, login string, credentials []content.CredentialInfo) ([]content.CredentialInfo, error) {
	db, err := TryToOpenDBConnection(st.connectionString)
	if err != nil {
		return nil, err
	}
	var userID int
	userID, err = GetUserID(ctx, login, db)
	if err != nil {
		return nil, err
	}
	if userID == 0 {
		return nil, srverror.ErrLoginNotFound
	}
	for index, cred := range credentials {
		var credID int
		err = db.QueryRowContext(ctx, "INSERT INTO credentials(resource,login,passwd,modification_time) values($1,$2,$3,$4) RETURNING id_credential;", cred.Resource, cred.Login, cred.Password, cred.ModificationTime).Scan(&credID)

		if err != nil {
			log.Println(err)
			continue
		}
		credentials[index].ID = credID

		_, err = db.ExecContext(ctx, "INSERT INTO users_credentials (ID_User,id_credential) values ($1,$2)", userID, credID)
		if err != nil {
			log.Println(err)
			//mb remove credential?
		}
	}
	return credentials, nil
}
func (st *StorageDB) AddNewFiles(ctx context.Context, login string, files []content.BinaryFileInfo) ([]content.BinaryFileInfo, error) {
	db, err := TryToOpenDBConnection(st.connectionString)
	if err != nil {
		return nil, err
	}
	var userID int
	userID, err = GetUserID(ctx, login, db)
	if err != nil {
		return nil, err
	}
	if userID == 0 {
		return nil, srverror.ErrLoginNotFound
	}
	for index, file := range files {
		var fileID int
		err = db.QueryRowContext(ctx, "INSERT INTO files(description,file_name,file_path,file_size,modification_time) values($1,$2,$3,$4,$5) RETURNING id_file;", file.Description, file.FileName, file.FilePath, file.FileSize, file.ModificationTime).Scan(&fileID)

		if err != nil {
			log.Println(err)
			continue
		}
		files[index].ID = fileID

		_, err = db.ExecContext(ctx, "INSERT INTO users_files (ID_User,id_file) values ($1,$2)", userID, fileID)
		if err != nil {
			log.Println(err)
			//mb remove credential?
		}
	}
	return files, nil
}
func (st *StorageDB) AddNewTexts(ctx context.Context, login string, texts []content.TextInfo) ([]content.TextInfo, error) {
	db, err := TryToOpenDBConnection(st.connectionString)
	if err != nil {
		return nil, err
	}
	var userID int
	userID, err = GetUserID(ctx, login, db)
	if err != nil {
		return nil, err
	}
	if userID == 0 {
		return nil, srverror.ErrLoginNotFound
	}
	for index, text := range texts {
		var textID int
		err = db.QueryRowContext(ctx, "INSERT INTO texts(description,content,modification_time) values($1,$2,$3) RETURNING id_text;", text.Description, text.Content, text.ModificationTime).Scan(&textID)

		if err != nil {
			log.Println(err)
			continue
		}
		texts[index].ID = textID

		_, err = db.ExecContext(ctx, "INSERT INTO users_texts (ID_User,id_text) values ($1,$2)", userID, textID)
		if err != nil {
			log.Println(err)
			//mb remove credential?
		}
	}
	return texts, nil
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

	_, err = db.ExecContext(ctx, "DELETE FROM Users_CreditCards where ID_CreditCard = any($1) and ID_User = $2", pq.Array(creditCardIDs), userID)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = db.ExecContext(ctx, "DELETE FROM CreditCards where ID_CreditCard = any($1)", pq.Array(creditCardIDs))
	if err != nil {
		return err
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

	_, err = db.ExecContext(ctx, "DELETE FROM Users_Credentials where id_credential = any($1) and ID_User = $2", pq.Array(credentialIDs), userID)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = db.ExecContext(ctx, "DELETE FROM Credentials where id_credential = any($1)", pq.Array(credentialIDs))
	if err != nil {
		return err
	}
	return nil
}
func (st *StorageDB) RemoveFiles(ctx context.Context, login string, fileIDs []int) error {
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

	_, err = db.ExecContext(ctx, "DELETE FROM users_files where id_file = any($1) and ID_User = $2", pq.Array(fileIDs), userID)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = db.ExecContext(ctx, "DELETE FROM files where id_file = any($1)", pq.Array(fileIDs))
	if err != nil {
		return err
	}
	return nil
}
func (st *StorageDB) RemoveTexts(ctx context.Context, login string, textIDs []int) error {
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

	_, err = db.ExecContext(ctx, "DELETE FROM users_texts where id_text = any($1) and ID_User = $2", pq.Array(textIDs), userID)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = db.ExecContext(ctx, "DELETE FROM texts where id_text = any($1)", pq.Array(textIDs))
	if err != nil {
		return err
	}
	return nil
}

func (st *StorageDB) UpdateFiles(ctx context.Context, login string, files []content.BinaryFileInfo) error {
	return nil
}
func (st *StorageDB) UpdateTexts(ctx context.Context, login string, texts []content.TextInfo) error {
	db, err := TryToOpenDBConnection(st.connectionString)
	if err != nil {
		return err
	}

	for _, text := range texts {
		_, err = db.ExecContext(ctx, "UPDATE Texts SET description = $2, content = $3, modification_time = $4 where ID_Text = $1", text.ID, text.Description, text.Content, text.ModificationTime)
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
		_, err = db.ExecContext(ctx, "UPDATE CreditCards SET bank = $2, card_number = $3, valid_thru = $4,cvv = $5, modification_time = $6 WHERE id_creditcard = $1", card.ID, card.Bank, card.CardNumber, card.ValidThru, card.CVV, card.ModificationTime)
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
		_, err = db.ExecContext(ctx, "UPDATE Credentials SET resource = $2, login = $3, passwd = $4,modification_time = $5 where id_credential = $1", cred.ID, cred.Resource, cred.Login, cred.Password, cred.ModificationTime)
		if err != nil {
			log.Println(err)
			continue
		}
	}
	return nil
}

func (st *StorageDB) GetCreditCards(ctx context.Context, login string, cardIDs []int) ([]content.CreditCardInfo, error) {
	if cardIDs == nil {
		return nil, nil
	}
	db, err := TryToOpenDBConnection(st.connectionString)
	if err != nil {
		return nil, err
	}
	var userID int
	userID, err = GetUserID(ctx, login, db)
	if err != nil {
		return nil, err
	}
	if userID == 0 {
		return nil, srverror.ErrLoginNotFound
	}

	var rows *sql.Rows
	rows, err = db.QueryContext(ctx, "SELECT creditcards.id_creditcard,creditcards.bank,CreditCards.card_number,CreditCards.cvv,CreditCards.valid_thru,CreditCards.modification_time FROM creditcards JOIN users_creditcards on creditcards.id_creditcard = users_creditcards.ID_CreditCard WHERE Users_CreditCards.id_user = $1 and creditcards.id_creditcard = any($2)", userID, pq.Array(cardIDs))

	var res []content.CreditCardInfo
	for rows.Next() {
		var r content.CreditCardInfo
		err = rows.Scan(&r.ID, &r.Bank, &r.CardNumber, &r.CVV, &r.ValidThru, &r.ModificationTime)
		if err != nil {
			break
		}
		res = append(res, r)
	}
	return res, nil
}
func (st *StorageDB) GetCredentials(ctx context.Context, login string, credIDs []int) ([]content.CredentialInfo, error) {
	if credIDs == nil {
		return nil, nil
	}
	db, err := TryToOpenDBConnection(st.connectionString)
	if err != nil {
		return nil, err
	}
	var userID int
	userID, err = GetUserID(ctx, login, db)
	if err != nil {
		return nil, err
	}
	if userID == 0 {
		return nil, srverror.ErrLoginNotFound
	}

	var rows *sql.Rows
	rows, err = db.QueryContext(ctx, "SELECT credentials.id_credential,credentials.resource,credentials.login,credentials.passwd,credentials.modification_time FROM credentials JOIN users_credentials on credentials.id_credential = users_credentials.id_credential WHERE users_credentials.id_user = $1 and credentials.id_credential = any($2)", userID, pq.Array(credIDs))

	var res []content.CredentialInfo
	for rows.Next() {
		var r content.CredentialInfo
		err = rows.Scan(&r.ID, &r.Resource, &r.Login, &r.Password, &r.ModificationTime)
		if err != nil {
			break
		}
		res = append(res, r)
	}
	return res, nil
}
func (st *StorageDB) GetFiles(ctx context.Context, login string, fileIDs []int) ([]content.BinaryFileInfo, error) {
	return nil, nil
}
func (st *StorageDB) GetTexts(ctx context.Context, login string, textIDs []int) ([]content.TextInfo, error) {
	if textIDs == nil {
		return nil, nil
	}
	db, err := TryToOpenDBConnection(st.connectionString)
	if err != nil {
		return nil, err
	}
	var userID int
	userID, err = GetUserID(ctx, login, db)
	if err != nil {
		return nil, err
	}
	if userID == 0 {
		return nil, srverror.ErrLoginNotFound
	}

	var rows *sql.Rows
	rows, err = db.QueryContext(ctx, "SELECT texts.id_text,texts.description,texts.content,texts.modification_time FROM texts JOIN users_texts on texts.id_text = users_texts.id_text WHERE users_texts.id_user = $1 and texts.id_text = any($2)", userID, pq.Array(textIDs))

	var res []content.TextInfo
	for rows.Next() {
		var r content.TextInfo
		err = rows.Scan(&r.ID, &r.Description, &r.Content, &r.ModificationTime)
		if err != nil {
			break
		}
		res = append(res, r)
	}
	return res, nil
}

func (st *StorageDB) GetModtimeWithIDs(ctx context.Context, login string, dataType string) (map[int]time.Time, error) {
	db, err := TryToOpenDBConnection(st.connectionString)
	if err != nil {
		return nil, err
	}
	var userID int
	err = db.QueryRowContext(ctx, "SELECT id_user FROM Users WHERE login = $1", login).Scan(&userID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, srverror.ErrLoginNotFound
		default:
			return nil, err
		}
	}

	var rows *sql.Rows

	switch dataType {
	case urlsuff.DatatypeCredential:
		rows, err = db.QueryContext(ctx, "SELECT credentials.id_credential,modification_time FROM credentials JOIN users_credentials on credentials.id_credential = users_credentials.id_credential WHERE users_credentials.id_user = $1", userID)
	case urlsuff.DatatypeCreditCard:
		rows, err = db.QueryContext(ctx, "SELECT creditcards.id_creditcard,modification_time FROM creditcards JOIN users_creditcards on creditcards.id_creditcard = users_creditcards.ID_CreditCard WHERE Users_CreditCards.id_user = $1", userID)
	case urlsuff.DatatypeText:
		rows, err = db.QueryContext(ctx, "SELECT texts.id_text,modification_time FROM texts JOIN users_texts on texts.id_text = users_texts.id_text WHERE users_texts.id_user = $1", userID)
	case urlsuff.DatatypeFile:
		rows, err = db.QueryContext(ctx, "SELECT files.id_file,modification_time FROM files JOIN users_files on files.id_file = users_files.id_file WHERE users_files.id_user = $1", userID)
	default:
		return nil, srverror.ErrIncorrectDataTpe
	}

	if err != nil {
		return nil, err
	}

	type row struct {
		ID      int
		modtime time.Time
	}

	res := make(map[int]time.Time)
	for rows.Next() {
		var r row
		err = rows.Scan(&r.ID, &r.modtime)
		if err != nil {
			break
		}
		res[r.ID] = r.modtime
	}
	if err != nil {
		return nil, err
	}
	return res, nil
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
