USE diplomacy;
DROP TABLE IF EXISTS game_territories;
DROP TABLE IF EXISTS territory;
CREATE TABLE territory (
  id int(11) NOT NULL AUTO_INCREMENT,
  game_id int(11) NOT NULL,
  owner char(15),
  country char(11),
  PRIMARY KEY (id)
);

-- DROP TABLE IF EXISTS game_territories;
-- CREATE TABLE game_territories (
--     game_id int(11) NOT NULL,
--     territory_id int(11) NOT NULL,
--     PRIMARY KEY (game_id, territory_id),
--     CONSTRAINT Constr_Game_Territories_Game_id_fk
--         FOREIGN KEY Game_id_fk (game_id) REFERENCES games (id),
--     CONSTRAINT Constr_Game_Territories_Territory_id_fk
--         FOREIGN KEY Piece_id_fk (territory_id) REFERENCES territory (id)
-- ) CHARACTER SET utf8 COLLATE utf8_general_ci;