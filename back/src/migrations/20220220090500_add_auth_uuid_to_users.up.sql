DROP PROCEDURE IF EXISTS `setup_auth_uuid`;

CREATE PROCEDURE `setup_auth_uuid` ()
  BEGIN
    IF NOT EXISTS (
      SELECT * FROM `information_schema`.`columns` WHERE `table_name` = 'users' AND `column_name` = 'auth_uuid'
    ) THEN
      ALTER TABLE `users` ADD `auth_uuid` VARCHAR (255) NOT NULL DEFAULT "";
    END IF;
  END;

CALL `setup_auth_uuid`;

DROP PROCEDURE IF EXISTS `setup_auth_uuid`;
