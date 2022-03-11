DROP PROCEDURE IF EXISTS `setup_password`;

CREATE PROCEDURE `setup_password` ()
  BEGIN
    IF NOT EXISTS (
      SELECT * FROM `information_schema`.`columns` WHERE `table_name` = 'users' AND `column_name` = 'password'
    ) THEN
      ALTER TABLE `users` ADD `password` VARCHAR (255) NOT NULL DEFAULT "";
    END IF;
  END;

CALL `setup_password`;

DROP PROCEDURE IF EXISTS `setup_password`;
