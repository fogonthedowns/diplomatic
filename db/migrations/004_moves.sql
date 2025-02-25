USE diplomacy;
DROP TABLE IF EXISTS moves;
CREATE TABLE moves (
  id int(11) NOT NULL AUTO_INCREMENT,
  location_start VARCHAR(10) CHARACTER SET utf8 NOT NULL,
  location_submitted VARCHAR(10) CHARACTER SET utf8 NOT NULL,
  location_resolved VARCHAR(10) CHARACTER SET utf8,
  phase int NOT NULL,
  game_year VARCHAR(5),
  type VARCHAR(10) NOT NULL,
  piece_id int(11) NOT NULL,
  game_id int(11) NOT NULL,
  PRIMARY KEY (id)
);

-- DROP TABLE IF EXISTS pieces_moves;
-- CREATE TABLE pieces_moves (
--     piece_id int(11) NOT NULL,
--     move_id int(11) NOT NULL,
--     phase int NOT NULL,
--     game_year DATE,
--     PRIMARY KEY (piece_id, move_id),
--     CONSTRAINT Constr_Pieces_Moves_Piece_id_fk
--         FOREIGN KEY Piece_id_fk (piece_id) REFERENCES pieces (id),
--     CONSTRAINT Constr_Pieces_Moves_Move_id_fk
--         FOREIGN KEY Move_id_fk (move_id) REFERENCES moves (id)
-- ) CHARACTER SET utf8 COLLATE utf8_general_ci;