USE diplomacy;
DROP TABLE IF EXISTS game_pieces;
DROP TABLE IF EXISTS pieces;
CREATE TABLE pieces (
  id int(11) NOT NULL AUTO_INCREMENT,
  game_id int(11) NOT NULL,
  owner_id int(11) NOT NULL,
  type char(5) NOT NULL,
  location VARCHAR(10) CHARACTER SET utf8 NOT NULL,
  PRIMARY KEY (id)
);

DROP TABLE IF EXISTS game_pieces;
CREATE TABLE game_pieces (
    game_id int(11) NOT NULL,
    piece_id int(11) NOT NULL,
    PRIMARY KEY (game_id, piece_id),
    CONSTRAINT Constr_Game_Pieces_Game_id_fk
        FOREIGN KEY Game_id_fk (game_id) REFERENCES games (id),
    CONSTRAINT Constr_Game_Pieces_piece_id_fk
        FOREIGN KEY Piece_id_fk (piece_id) REFERENCES pieces (id)
) CHARACTER SET utf8 COLLATE utf8_general_ci;