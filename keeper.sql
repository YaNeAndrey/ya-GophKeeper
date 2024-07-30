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