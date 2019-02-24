USE diplomacy;
DROP TABLE IF EXISTS pieces;
CREATE TABLE pieces (
  id int(11) NOT NULL AUTO_INCREMENT,
  game_id int(11) NOT NULL, 
  owner char(15),
  type char(5) NOT NULL,
  is_active boolean NOT NULL DEFAULT true,
  location char(11) CHARACTER SET utf8 NOT NULL,
  PRIMARY KEY (id)
);
