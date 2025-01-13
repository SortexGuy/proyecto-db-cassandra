package counters

type Counter struct {
	id_name    string
	current_id int64
}

/*
CREATE TABLE counters (
    id_name text PRIMARY KEY,
    current_id bigint
);
*/
