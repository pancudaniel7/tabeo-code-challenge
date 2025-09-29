CREATE TABLE appointments
(
    id         BINARY(16)    NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name  VARCHAR(100) NOT NULL,
    visit_date DATE         NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uq_visit_date (visit_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;