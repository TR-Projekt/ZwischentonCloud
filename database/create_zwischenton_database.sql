--
-- Create the ZwischentonCloud API database
--

-- First create the database
CREATE DATABASE IF NOT EXISTS `zwischenton_cloud_database`;

-- Create the tables in the newly created database
USE zwischenton_cloud_database;

/**

Create the basic entities

*/

-- Create the zwischenton table
CREATE TABLE IF NOT EXISTS `zwischentons` (

	`zwischenton_id` 			int unsigned 	 	NOT NULL AUTO_INCREMENT 											COMMENT 'The id of the zwischenton.',
	`zwischenton_version` 		timestamp 			NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() 	COMMENT 'The version of the zwischenton.',
	`zwischenton_is_valid` 	  	tinyint(1) 			NOT NULL DEFAULT 0 													COMMENT 'Boolean value indicating if the zwischenton should be distributed to users.',
	`zwischenton_name` 		    varchar(255)		NOT NULL DEFAULT ''													COMMENT 'The zwischenton name. The name needs to be unique.',
	`zwischenton_description` 	text 			    NOT NULL 						      								COMMENT 'The description of the zwischenton.',

PRIMARY 	KEY (`zwischenton_id`),
UNIQUE 	  	KEY `name` (`zwischenton_name`),
		    KEY `is_valid` (`zwischenton_is_valid`)

) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='The zwischenton table represents a zwischenton and its core properties.';

-- Create the situation table
CREATE TABLE IF NOT EXISTS `situations` (

	 `situation_id` 				int unsigned		NOT NULL AUTO_INCREMENT 											COMMENT 'The id of the situation.',
	 `situation_version` 	        timestamp 			NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() 	COMMENT 'The version of the situation.',
	 `situation_lat` 				float(10,6) 		NOT NULL DEFAULT 0.000000											COMMENT 'The latitude of the situation.',
	 `situation_lon` 				float(10,6) 		NOT NULL DEFAULT 0.000000											COMMENT 'The longitude of the situation.',
     `situation_radius` 			float(10,6) 		NOT NULL DEFAULT 0.000000											COMMENT 'The radius of the situation.',
	 `situation_description` 		text				NOT NULL					          								COMMENT 'The description of the situation.',

PRIMARY 	KEY (`situation_id`),
			KEY `latitude` (`situation_lat`),
			KEY `longitude` (`situation_lon`),
            KEY `radius` (`situation_radius`)

) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='The situation table represents a situation and its core properties.';

/**

Create the mapping tables to associate entities

*/

-- Create the table to map situations to zwischentons
CREATE TABLE IF NOT EXISTS `map_zwischenton_situation` (

	`map_id` 			        int unsigned 		    NOT NULL AUTO_INCREMENT		COMMENT 'The id of the map entry.',
	`associated_zwischenton` 	int unsigned 		    NOT NULL					COMMENT 'The id of the mapped zwischenton.',
	`associated_situation` 	    int unsigned 		    NOT NULL					COMMENT 'The id of the mapped situation.',

PRIMARY 	KEY (`map_id`),
FOREIGN 	KEY (`associated_zwischenton`)		REFERENCES zwischentons (zwischenton_id)	ON DELETE CASCADE 	ON UPDATE CASCADE,
FOREIGN 	KEY (`associated_situation`) 		REFERENCES situations (situation_id) 		ON DELETE CASCADE 	ON UPDATE CASCADE

) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='The table maps situations to zwischentons.';
