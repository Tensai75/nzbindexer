-- phpMyAdmin SQL Dump
-- version 5.1.1
-- https://www.phpmyadmin.net/
--

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
  `hash` varchar(32) NOT NULL,
  `header_hash` varchar(32) NOT NULL,
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
  `file` varchar(32) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `headers`
--

DROP TABLE IF EXISTS `headers`;
CREATE TABLE `headers` (
  `hash` varchar(32) NOT NULL,
  `files` int(11) NOT NULL,
  `total_files` int(11) NOT NULL,
  `size` bigint(20) NOT NULL,
  `date` int(11) NOT NULL,
  `poster` int(11) NOT NULL
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
  `file_hash` varchar(32) NOT NULL,
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
  ADD UNIQUE KEY `hash` (`hash`);
ALTER TABLE `files` ADD FULLTEXT KEY `subject` (`subject`);

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
  ADD UNIQUE KEY `file_to_group` (`group_id`,`file`);

--
-- Indexes for table `headers`
--
ALTER TABLE `headers`
  ADD UNIQUE KEY `hash` (`hash`);

--
-- Indexes for table `poster`
--
ALTER TABLE `poster`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `poster` (`poster`(255));

--
-- Indexes for table `segments`
--
ALTER TABLE `segments`
  ADD UNIQUE KEY `segment_id` (`segment_id`(255));

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `groups`
--
ALTER TABLE `groups`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `poster`
--
ALTER TABLE `poster`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
