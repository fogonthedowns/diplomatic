USE diplomacy;
DROP TABLE IF EXISTS users_games;
DROP TABLE IF EXISTS users;
CREATE TABLE users (
  id int(11) NOT NULL AUTO_INCREMENT,
  email VARCHAR(64) CHARACTER SET utf8 NOT NULL,
  user_name VARCHAR(64) CHARACTER SET utf8 NOT NULL,
  password VARCHAR(64) CHARACTER SET utf8 NOT NULL,
  PRIMARY KEY (id)
);

DROP TABLE IF EXISTS users_games;
CREATE TABLE users_games (
    user_id int(11) NOT NULL,
    game_id int(11) NOT NULL,
    country char(11) NOT NULL,
    PRIMARY KEY (user_id, game_id),
    CONSTRAINT Constr_Users_Games_User_id_fk
        FOREIGN KEY User_id_fk (user_id) REFERENCES users (id),
    CONSTRAINT Constr_Users_Games_Game_id_fk
        FOREIGN KEY Game_id_fk (game_id) REFERENCES games (id)
) CHARACTER SET utf8 COLLATE utf8_general_ci;
