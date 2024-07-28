--
-- Create the ZwischentonCloud Identity Database
--

-- First create the database
CREATE DATABASE IF NOT EXISTS `zwischenton_identity_database`;

-- Create the tables in the newly created database
USE zwischenton_identity_database;

/**
Create the basic entities
*/

-- Create the users table
CREATE TABLE IF NOT EXISTS `users` (

	`user_id` 			    int unsigned 	 	NOT NULL AUTO_INCREMENT 											    COMMENT 'The id of the user.',
	`user_email` 		    varchar(255)		NOT NULL													            COMMENT 'The email of the user. The email needs to be unique.',
	`user_password` 	    varchar(225) 	  	NOT NULL 												                COMMENT 'The password hash of the users password.',
	`user_createdat` 		timestamp 			NOT NULL DEFAULT current_timestamp()					      		    COMMENT 'The date and time the user was created.',
	`user_updatedat` 		timestamp 			NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()	    COMMENT 'The date and time the user data was last updated.',
    `user_role` 	  	    tinyint 		    NOT NULL DEFAULT 0											            COMMENT 'The role of the user.',

PRIMARY 	KEY (`user_id`),
UNIQUE 	    KEY (`user_email`)

) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='The user table represents a user that interacts with the Zwischenton backend.';

-- Create the service key table
CREATE TABLE IF NOT EXISTS `api_keys` (

	`api_key_id` 			int unsigned 	 	NOT NULL AUTO_INCREMENT 											COMMENT 'The id of the key.',
	`api_key` 	  	        varchar(225) 		NOT NULL 												            COMMENT 'The api key.',
    `api_key_comment` 	  	varchar(225) 		NOT NULL 												            COMMENT 'A comment about the api key.',

PRIMARY 	KEY (`api_key_id`),
UNIQUE 	  	KEY (`api_key`)

) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='This table contains all api keys.';

/**
Create the mapping tables to associate entities
*/

-- Create the table to map zwischentons to users
CREATE TABLE IF NOT EXISTS `map_zwischenton_user` (

    `map_id` 				 	int unsigned 		NOT NULL AUTO_INCREMENT		        COMMENT 'The id of the map entry.',
    `associated_zwischenton` 	int unsigned 		NOT NULL					        COMMENT 'The id of the mapped zwischenton.',
    `associated_user` 	    	int unsigned 		NOT NULL					        COMMENT 'The id of the mapped user.',

PRIMARY 	KEY (`map_id`),
UNIQUE 	  	KEY (`associated_zwischenton`),
FOREIGN 	KEY (`associated_user`)                 REFERENCES users (user_id)

) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='This table maps zwischentons to users.';

/**
Insert default admin user (default password: we4711)
*/

INSERT INTO  `users`(`user_id`, `user_email`, `user_password`, `user_role`) VALUES (0, 'admin@email.com', '$2a$12$YbAhewILx82tGkLtEZWiKOfYzBt85RSQtGXhxlQX2hV7qiP51xPES', 42);