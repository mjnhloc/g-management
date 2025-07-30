CREATE TABLE waitlists (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    class_id BIGINT NOT NULL,
    member_id BIGINT NOT NULL,
    joined_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    notified BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT fk_waitlist_class FOREIGN KEY (class_id) REFERENCES classes(id),
    CONSTRAINT fk_waitlist_member FOREIGN KEY (member_id) REFERENCES members(id)
);