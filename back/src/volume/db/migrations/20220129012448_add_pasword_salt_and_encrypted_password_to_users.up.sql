DROP PROCEDURE IF EXISTS `setup_password_encryption`;

CREATE PROCEDURE `setup_password_encryption` (
  IN `target_column` VARCHAR(20)
)
  BEGIN
    IF NOT EXISTS (
      SELECT * FROM `information_schema`.`columns` WHERE `table_name` = 'users' AND `column_name` = target_column
    ) THEN
      SET @SQL := CONCAT('ALTER TABLE `users` ADD `', target_column, '` VARCHAR (255) NOT NULL DEFAULT "";');
      PREPARE add_stmt FROM @SQL;
      EXECUTE add_stmt;
      DEALLOCATE PREPARE add_stmt;
    END IF;
  END;

CALL `setup_password_encryption`('password_salt');
CALL `setup_password_encryption`('encrypted_password');

DROP PROCEDURE IF EXISTS `setup_password_encryption`;
