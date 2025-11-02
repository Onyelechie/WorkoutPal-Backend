CREATE TABLE IF NOT EXISTS schedule (
                                        id SERIAL PRIMARY KEY,
                                        name TEXT NOT NULL,
                                        user_id INT NOT NULL,
                                        day_of_week INT NOT NULL,
                                        time_slot TIME NOT NULL,
                                        routine_length_minutes INT NOT NULL
);

CREATE TABLE IF NOT EXISTS schedule_routine (
    schedule_id INT NOT NULL REFERENCES schedule(id) ON DELETE CASCADE,
    routine_id INT NOT NULL REFERENCES workout_routine(id),
    position INT NOT NULL,
    PRIMARY KEY (schedule_id, routine_id)
);