USE diplomacy;
ALTER TABLE moves ADD second_location_submitted VARCHAR(10) CHARACTER SET utf8; 
ALTER TABLE moves ADD piece_owner varchar(15);

