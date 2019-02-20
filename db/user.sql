USE diplomacy;
DROP TABLE IF EXISTS users;
CREATE TABLE users (
  id int(11) NOT NULL AUTO_INCREMENT,
  email VARCHAR(64) CHARACTER SET utf8 NOT NULL,
  user_name VARCHAR(64) CHARACTER SET utf8 NOT NULL,
  password VARCHAR(64) CHARACTER SET utf8 NOT NULL,
  PRIMARY KEY (id)
);


CREATE TABLE usersgames (
    user int(11) NOT NULL,
    game int(11) NOT NULL,
    PRIMARY KEY (user, game),
    CONSTRAINT Constr_UsersGames_User_fk
        FOREIGN KEY User_fk (user) REFERENCES users (id)
        ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT Constr_UserGames_Game_fk
        FOREIGN KEY Game_fk (game) REFERENCES games (id)
        ON DELETE CASCADE ON UPDATE CASCADE
) CHARACTER SET utf8 COLLATE utf8_general_ci;
