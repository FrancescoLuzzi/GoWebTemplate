-- +goose Up

-- Password123!
INSERT INTO users(id,email,firstName,lastName,password) VALUES
('be1d078b-827b-438c-bd95-fbb5627115c3','root@test.com','admin','root','$argon2id$v=19$m=15000,t=2,p=2$ZW9QUTMwSlc5UmVrVTJtQQ$kRLJ4I+I0a/qLndBwL11ZKL+Vfw0nOpeoN09h0EbrS8');


-- +goose Down

SELECT 'down SQL query';
