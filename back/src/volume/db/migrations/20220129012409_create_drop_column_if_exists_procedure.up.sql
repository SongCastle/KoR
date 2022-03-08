DROP PROCEDURE IF EXISTS `drop_column_if_exists`;

CREATE PROCEDURE `drop_column_if_exists` (
  IN `target_table_name` VARCHAR(63),
  IN `target_column` VARCHAR(63)
)
  BEGIN
    IF EXISTS (
      SELECT * FROM `information_schema`.`columns` WHERE `table_name` = target_table_name AND `column_name` = target_column
    ) THEN
      SET @SQL := CONCAT('ALTER TABLE `', target_table_name, '` DROP `', target_column, '`');
      PREPARE drop_stmt FROM @SQL;
      EXECUTE drop_stmt;
      DEALLOCATE PREPARE drop_stmt;
    END IF;
  END;
