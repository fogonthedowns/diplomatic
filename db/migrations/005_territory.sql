USE diplomacy;
DROP TABLE IF EXISTS territory;
CREATE TABLE territory (
  id int(11) NOT NULL AUTO_INCREMENT,
  game_id int(11) NOT NULL,
  owner char(15),
  country char(11),
  PRIMARY KEY (id)
);
