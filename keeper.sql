CREATE TABLE IF NOT EXISTS Users(
    ID_User serial PRIMARY KEY,
    login VARCHAR(256) NOT NULL,
    passwd TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS Credentials (
    ID_Credential serial PRIMARY KEY,
    resource VARCHAR(256),
    login VARCHAR(256),
    passwd TEXT,
    modification_time timestamp
);

CREATE TABLE IF NOT EXISTS Users_Credentials(
    ID_User_Credentials serial PRIMARY KEY,
    ID_User bigint REFERENCES Users(ID_User),
    ID_Credential bigint REFERENCES Credentials (ID_Credential)
);

CREATE TABLE IF NOT EXISTS CreditCards (
    ID_CreditCard bigserial PRIMARY KEY,
    bank varchar(50),
    card_number varchar(16) NOT NULL,
    valid_thru timestamp,
    cvv varchar(3),
    modification_time timestamp
);

CREATE TABLE IF NOT EXISTS Users_CreditCards(
    ID_User_CreditCard SERIAL PRIMARY KEY,
    ID_User bigint REFERENCES Users (ID_User),
    ID_CreditCard bigint REFERENCES CreditCards (ID_CreditCard)
);

CREATE TABLE IF NOT EXISTS Texts (
    ID_Text bigserial PRIMARY KEY,
    description TEXT,
    content TEXT NOT NULL,
    modification_time timestamp
);

CREATE TABLE IF NOT EXISTS Users_CreditCards(
    ID_User_CreditCard SERIAL PRIMARY KEY,
    ID_User bigint REFERENCES Users (ID_User),
    ID_Text bigint REFERENCES Texts (ID_Text)
);


CREATE TABLE IF NOT EXISTS Files (
    ID_File bigserial PRIMARY KEY,
    description TEXT,
    file_name varchar(256) NOT NULL,
    file_path varchar(256) NOT NULL,
    file_size int NOT NULL,
    modification_time timestamp
);

INSERT INTO Users(login,passwd) values('admin', crypt('admin', gen_salt('bf')));
select * from users;
INSERT INTO CreditCards(bank,card_number,valid_thru,cvv,modification_time) values('SBER','1234567890123456','2022-10-10',123,'2022-10-10') RETURNING ID_CreditCard;
INSERT INTO CreditCards(bank,card_number,valid_thru,cvv,modification_time) values('ALPHA','1234567890123456','2024-10-10',123,'2023-10-10') RETURNING ID_CreditCard;

INSERT INTO Users_CreditCards (ID_User,ID_CreditCard) values (3,1);
INSERT INTO Users_CreditCards (ID_User,ID_CreditCard) values (3,2);
select * from Users_CreditCards

DELETE FROM Users;

select * from creditcards




UPDATE CreditCards SET bank = 'Lol', card_number = '1234567890123456', valid_thru = '2024-10-12',cvv = '345', modification_time = '2024-10-12' WHERE id_creditcard = 2;

SELECT creditcards.id_creditcard,creditcards.bank,CreditCards.card_number,CreditCards.cvv,CreditCards.valid_thru,CreditCards.modification_time FROM creditcards JOIN users_creditcards on creditcards.id_creditcard = users_creditcards.ID_CreditCard WHERE Users_CreditCards.id_user = 3 and creditcards.id_creditcard in  (2,3)

SELECT creditcards.id_creditcard,creditcards.bank,CreditCards.card_number,CreditCards.cvv,CreditCards.valid_thru,CreditCards.modification_time FROM creditcards JOIN users_creditcards on creditcards.id_creditcard = users_creditcards.ID_CreditCard WHERE Users_CreditCards.id_user = 3
