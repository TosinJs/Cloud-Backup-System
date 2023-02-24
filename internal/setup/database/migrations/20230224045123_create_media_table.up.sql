CREATE TABLE IF NOT EXISTS `railway`.`Media` (
  `id` INT NOT NULL,
  `username` VARCHAR(255) NOT NULL,
  `filepath` VARCHAR(255) NOT NULL,
  `flag_count` INT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  INDEX `username_idx` (`username` ASC) VISIBLE,
  CONSTRAINT `username`
    FOREIGN KEY (`username`)
    REFERENCES `railway`.`Users` (`username`)
    ON DELETE CASCADE
    ON UPDATE CASCADE);
