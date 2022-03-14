-- phpMyAdmin SQL Dump
-- version 5.1.1
-- https://www.phpmyadmin.net/


SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `nzbindexer`
--
CREATE DATABASE IF NOT EXISTS `nzbindexer` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
USE `nzbindexer`;

-- --------------------------------------------------------

--
-- Table structure for table `files`
--

DROP TABLE IF EXISTS `files`;
CREATE TABLE `files` (
  `id` int(11) NOT NULL,
  `header_id` int(11) NOT NULL,
  `subject` text NOT NULL,
  `file_no` int(11) NOT NULL,
  `segments` int(11) NOT NULL,
  `total_segments` int(11) NOT NULL,
  `size` bigint(20) NOT NULL,
  `date` int(11) NOT NULL,
  `poster` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `file_hashes`
--

DROP TABLE IF EXISTS `file_hashes`;
CREATE TABLE `file_hashes` (
  `id` int(11) NOT NULL,
  `hash` varchar(32) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `groups`
--

DROP TABLE IF EXISTS `groups`;
CREATE TABLE `groups` (
  `id` int(11) NOT NULL,
  `group_name` text NOT NULL,
  `first_message_id` int(11) DEFAULT NULL,
  `last_message_id` bigint(11) DEFAULT NULL,
  `current_message_id` bigint(11) DEFAULT NULL,
  `headers` int(11) DEFAULT 0,
  `files` int(11) DEFAULT 0,
  `segments` bigint(20) DEFAULT 0,
  `size` bigint(20) DEFAULT 0,
  `date` int(11) DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `groups_to_files`
--

DROP TABLE IF EXISTS `groups_to_files`;
CREATE TABLE `groups_to_files` (
  `group_id` int(11) NOT NULL,
  `file_id` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `headers`
--

DROP TABLE IF EXISTS `headers`;
CREATE TABLE `headers` (
  `id` int(11) NOT NULL,
  `files` int(11) NOT NULL,
  `total_files` int(11) NOT NULL,
  `size` bigint(20) NOT NULL,
  `date` int(11) NOT NULL,
  `poster` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `header_hashes`
--

DROP TABLE IF EXISTS `header_hashes`;
CREATE TABLE `header_hashes` (
  `id` int(11) NOT NULL,
  `hash` varchar(32) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `poster`
--

DROP TABLE IF EXISTS `poster`;
CREATE TABLE `poster` (
  `id` int(11) NOT NULL,
  `poster` text DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `segments`
--

DROP TABLE IF EXISTS `segments`;
CREATE TABLE `segments` (
  `id` int(11) NOT NULL,
  `file_id` int(11) NOT NULL,
  `segment_id` text NOT NULL,
  `segment_no` int(11) NOT NULL,
  `size` int(11) NOT NULL,
  `date` int(11) NOT NULL,
  `poster` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Indexes for dumped tables
--

--
-- Indexes for table `files`
--
ALTER TABLE `files`
  ADD PRIMARY KEY (`id`),
  ADD KEY `header_id` (`header_id`),
  ADD KEY `files_to_poster` (`poster`);
ALTER TABLE `files` ADD FULLTEXT KEY `subject` (`subject`);

--
-- Indexes for table `file_hashes`
--
ALTER TABLE `file_hashes`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `file_hash` (`hash`(32));

--
-- Indexes for table `groups`
--
ALTER TABLE `groups`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `group_name` (`group_name`(255));

--
-- Indexes for table `groups_to_files`
--
ALTER TABLE `groups_to_files`
  ADD UNIQUE KEY `file_to_group` (`group_id`,`file_id`),
  ADD KEY `files` (`file_id`);

--
-- Indexes for table `headers`
--
ALTER TABLE `headers`
  ADD PRIMARY KEY (`id`),
  ADD KEY `headers_to_poster` (`poster`);

--
-- Indexes for table `header_hashes`
--
ALTER TABLE `header_hashes`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `header_hash` (`hash`(32));

--
-- Indexes for table `poster`
--
ALTER TABLE `poster`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `poster` (`poster`(255)) USING HASH;

--
-- Indexes for table `segments`
--
ALTER TABLE `segments`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `segment_id` (`segment_id`(255)),
  ADD KEY `file_id` (`file_id`),
  ADD KEY `segments_to_poster` (`poster`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `file_hashes`
--
ALTER TABLE `file_hashes`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `groups`
--
ALTER TABLE `groups`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `headers`
--
ALTER TABLE `headers`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `header_hashes`
--
ALTER TABLE `header_hashes`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `poster`
--
ALTER TABLE `poster`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `segments`
--
ALTER TABLE `segments`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `files`
--
ALTER TABLE `files`
  ADD CONSTRAINT `file_id` FOREIGN KEY (`id`) REFERENCES `file_hashes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `files_to_group` FOREIGN KEY (`id`) REFERENCES `groups_to_files` (`file_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `files_to_header` FOREIGN KEY (`header_id`) REFERENCES `header_hashes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `files_to_poster` FOREIGN KEY (`poster`) REFERENCES `poster` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `groups_to_files`
--
ALTER TABLE `groups_to_files`
  ADD CONSTRAINT `groups_to_files` FOREIGN KEY (`group_id`) REFERENCES `groups` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `headers`
--
ALTER TABLE `headers`
  ADD CONSTRAINT `header_id` FOREIGN KEY (`id`) REFERENCES `header_hashes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `headers_to_poster` FOREIGN KEY (`poster`) REFERENCES `poster` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `segments`
--
ALTER TABLE `segments`
  ADD CONSTRAINT `segments_to_file` FOREIGN KEY (`file_id`) REFERENCES `file_hashes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `segments_to_poster` FOREIGN KEY (`poster`) REFERENCES `poster` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
