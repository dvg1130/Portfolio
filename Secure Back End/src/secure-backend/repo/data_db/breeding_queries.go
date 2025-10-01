package datadb

var GetBreedingMales = "SELECT * FROM breeding WHERE male1_suid = ? UNION SELECT * FROM breeding WHERE male2_suid = ? UNION SELECT * FROM breeding WHERE male3_suid = ? UNION SELECT * FROM breeding WHERE male4_suid = ?"

var GetBreedingUuid = "SELECT breeding_uuid FROM breeding WHERE uuid = ? AND event_id = ? "

// var GetAllBreeding = "SELECT breeding_year FROM breeding WHERE uuid = ? AND suid = ? ORDER BY breeding_year DESC"

var AddBreedingEvent = "INSERT INTO breeding (uuid, event_id, female_suid, male1_suid, male2_suid, male3_suid, male4_suid, female_sid, male1_sid, male2_sid, male3_sid, male4_sid, breeding_year, breeding_season, female_weight, male1_weight, male2_weight, male3_weight, male4_weight, cooling_start, cooling_end, warming_start, warming_end, pairing1_date, pairing2_date, pairing3_date, pairing4_date, gravid_date,lay_date, clutch_size, clutch_survive, outcome, notes) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?, ?, ?, ?, ?, ?, ?, ?, ?)"

var GetAllBreedingBySnake = "SELECT breeding_year FROM breeding WHERE uuid = ? AND (female_suid = ? OR male1_suid = ? OR male2_suid = ? OR male3_suid = ? OR male4_suid = ? ) ORDER BY breeding_year DESC;"

var GetAllBreedingsByUser = "SELECT breeding_uuid, event_id, breeding_year, breeding_season FROM breeding WHERE uuid = ? ORDER BY breeding_year DESC;"

var UpdateBreed = "UPDATE breeding SET %s WHERE uuid = ? AND breeding_uuid = ?"

var GetBreedingEvent = `
SELECT event_id, female_sid, male1_sid, male2_sid, male3_sid, male4_sid,
       breeding_year, breeding_season, female_weight, male1_weight, male2_weight,
       male3_weight, male4_weight, cooling_start, cooling_end, warming_start,
       warming_end, pairing1_date, pairing2_date, pairing3_date, pairing4_date,
       gravid_date, lay_date, clutch_size, clutch_survive, outcome, notes
FROM breeding
WHERE uuid = ? AND breeding_uuid = ?;
`
var DeleteBreedingEvent = "DELETE FROM breeding WHERE uuid = ? AND breeding_uuid = ?"
