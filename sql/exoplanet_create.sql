CREATE TABLE IF NOT EXISTS exoplanet (
    id SERIAL PRIMARY KEY,
    planet_name TEXT NOT NULL,
    host_name TEXT,
    system_number INT,
    discovery_method TEXT,
    year_discovered DATE
);