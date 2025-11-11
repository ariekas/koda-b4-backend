ALTER TABLE users
  ADD CONSTRAINT fk_users_profile
  FOREIGN KEY (profile_id) REFERENCES profile (id);