--cantidad de veces que pasaron por la caja
SELECT
	box_id,
	box_name,
	COUNT(*)
FROM
	tracker
GROUP BY
	box_id,
	box_name;
-- cantidad de veces que paso una session por cada caja
SELECT
	log_id,
	box_id,
	box_name,
	COUNT(*)
FROM
	tracker
GROUP BY
	log_id,
	box_id,
	box_name;

-- son la cantidad de veces que paso una session distinta por la caja
SELECT
t.box_id, t.box_name, count(*)
FROM
	(
		SELECT
			log_id,
			box_id,
			box_name
		FROM
			tracker 
		GROUP BY
			log_id,
			box_id,
			box_name
	) t
GROUP BY t.box_id, t.box_name;




