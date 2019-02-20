USE diplomacy;
DROP TABLE IF EXISTS pieces_moves;
DROP TABLE IF EXISTS moves;
CREATE TABLE moves (
  id int(11) NOT NULL AUTO_INCREMENT,
  location_start VARCHAR(10) CHARACTER SET utf8 NOT NULL,
  location_submitted VARCHAR(10) CHARACTER SET utf8 NOT NULL,
  location_resolved VARCHAR(10) CHARACTER SET utf8 NOT NULL,
  type int(5) NOT NULL,
  piece_id int(11) NOT NULL,
  PRIMARY KEY (id)
);

DROP TABLE IF EXISTS pieces_moves;
CREATE TABLE pieces_moves (
    piece_id int(11) NOT NULL,
    move_id int(11) NOT NULL,
    PRIMARY KEY (piece_id, move_id),
    CONSTRAINT Constr_Pieces_Moves_Piece_id_fk
        FOREIGN KEY Piece_id_fk (piece_id) REFERENCES pieces (id),
    CONSTRAINT Constr_Pieces_Moves_Move_id_fk
        FOREIGN KEY Move_id_fk (move_id) REFERENCES moves (id)
) CHARACTER SET utf8 COLLATE utf8_general_ci;