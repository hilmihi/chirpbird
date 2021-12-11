-- phpMyAdmin SQL Dump
-- version 5.1.1
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1
-- Waktu pembuatan: 10 Des 2021 pada 05.09
-- Versi server: 10.4.21-MariaDB
-- Versi PHP: 7.3.31

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `chirpbird`
--
CREATE DATABASE IF NOT EXISTS `chirpbird` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE `chirpbird`;

-- --------------------------------------------------------

--
-- Struktur dari tabel `messages`
--

DROP TABLE IF EXISTS `messages`;
CREATE TABLE `messages` (
  `id` int(11) NOT NULL,
  `createdate` datetime(3) NOT NULL,
  `updatedate` datetime(3) NOT NULL,
  `deletedate` datetime(3) DEFAULT NULL,
  `seqid` int(11) NOT NULL,
  `room` varchar(36) NOT NULL,
  `from` bigint(20) NOT NULL,
  `content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (json_valid(`content`))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data untuk tabel `messages`
--

INSERT INTO `messages` (`id`, `createdate`, `updatedate`, `deletedate`, `seqid`, `room`, `from`, `content`) VALUES
(1, '2021-12-05 17:04:26.000', '2021-12-05 17:04:26.000', NULL, 1, 'grpAsdsWfGdgr', 2, '\"Core Error - Bus Dumped\"'),
(2, '2021-12-05 13:50:37.732', '2021-12-05 13:50:37.732', NULL, 2, 'grpAsdsWfGdgr', 1, '\"aaaaa\"'),
(3, '2021-12-05 13:52:40.633', '2021-12-05 13:52:40.633', NULL, 3, 'grpAsdsWfGdgr', 1, '\"bbbbb\"'),
(4, '2021-12-05 13:53:22.907', '2021-12-05 13:53:22.907', NULL, 4, 'grpAsdsWfGdgr', 1, '\"ccccc\"'),
(5, '2021-12-05 13:53:56.885', '2021-12-05 13:53:56.885', NULL, 5, 'grpAsdsWfGdgr', 1, '\"ddddd\"'),
(6, '2021-12-07 09:25:29.138', '2021-12-07 09:25:29.138', NULL, 1, 'caaf8725-0adc-4ede-8732-d', 1, '\"hellooo\"'),
(7, '2021-12-08 10:05:46.161', '2021-12-08 10:05:46.161', NULL, 1, '47b9a9ac-4c76-44c2-a434-5', 1, '\"halooo\"'),
(8, '2021-12-08 10:06:04.877', '2021-12-08 10:06:04.877', NULL, 2, '47b9a9ac-4c76-44c2-a434-5', 4, '\"haii\"'),
(9, '2021-12-08 10:06:43.402', '2021-12-08 10:06:43.402', NULL, 3, '47b9a9ac-4c76-44c2-a434-5', 1, '\"aaaa\"'),
(10, '2021-12-08 10:07:07.613', '2021-12-08 10:07:07.613', NULL, 4, '47b9a9ac-4c76-44c2-a434-5', 4, '\"bbbb\"'),
(11, '2021-12-08 10:26:28.925', '2021-12-08 10:26:28.925', NULL, 5, '47b9a9ac-4c76-44c2-a434-5', 1, '\"ccccc\"'),
(12, '2021-12-08 10:26:37.656', '2021-12-08 10:26:37.656', NULL, 6, '47b9a9ac-4c76-44c2-a434-5', 4, '\"dddd\"'),
(13, '2021-12-09 02:12:38.770', '2021-12-09 02:12:38.770', NULL, 7, '47b9a9ac-4c76-44c2-a434-5', 1, '\"pagi bos\"'),
(14, '2021-12-09 02:37:23.963', '2021-12-09 02:37:23.963', NULL, 8, '47b9a9ac-4c76-44c2-a434-5', 4, '\"ada apa?\"'),
(15, '2021-12-09 02:37:46.778', '2021-12-09 02:37:46.778', NULL, 9, '47b9a9ac-4c76-44c2-a434-5', 1, '\"gada apa2\"'),
(16, '2021-12-09 02:41:01.097', '2021-12-09 02:41:01.097', NULL, 6, 'grpAsdsWfGdgr', 1, '\"aaaa\"'),
(17, '2021-12-09 02:43:28.317', '2021-12-09 02:43:28.317', NULL, 7, 'grpAsdsWfGdgr', 1, '\"aaaaa\"'),
(18, '2021-12-09 03:08:21.813', '2021-12-09 03:08:21.813', NULL, 8, 'grpAsdsWfGdgr', 1, '\"ccccc\"'),
(19, '2021-12-09 03:09:01.420', '2021-12-09 03:09:01.420', NULL, 9, 'grpAsdsWfGdgr', 2, '\"ccccc\"'),
(20, '2021-12-09 03:12:23.150', '2021-12-09 03:12:23.150', NULL, 10, 'grpAsdsWfGdgr', 1, '\"aaaa\"'),
(21, '2021-12-09 03:12:31.977', '2021-12-09 03:12:31.977', NULL, 11, 'grpAsdsWfGdgr', 2, '\"cccc\"'),
(22, '2021-12-09 04:34:32.793', '2021-12-09 04:34:32.793', NULL, 12, 'grpAsdsWfGdgr', 2, '\"dddd\"');

-- --------------------------------------------------------

--
-- Struktur dari tabel `rooms`
--

DROP TABLE IF EXISTS `rooms`;
CREATE TABLE `rooms` (
  `id` int(11) NOT NULL,
  `createdate` datetime(3) NOT NULL,
  `updatedate` datetime(3) NOT NULL,
  `name` varchar(36) NOT NULL,
  `owner` int(11) NOT NULL DEFAULT 0,
  `seqid` int(11) NOT NULL DEFAULT 0,
  `public` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (json_valid(`public`)),
  `tags` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (json_valid(`tags`))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data untuk tabel `rooms`
--

INSERT INTO `rooms` (`id`, `createdate`, `updatedate`, `name`, `owner`, `seqid`, `public`, `tags`) VALUES
(1, '2021-12-05 16:26:34.000', '2021-12-05 16:26:34.000', 'grpAsdsWfGdgr', 1, 0, '{\"rn\":\"Global Chat\"}', '[\"All\"]'),
(21, '2021-12-07 08:39:09.974', '2021-12-07 08:39:09.974', '1edbfda7-e4ee-44a6-86be-7', 0, 0, '{\"rn\":\"Global Chat2\"}', NULL),
(22, '2021-12-07 09:01:30.939', '2021-12-07 09:01:30.939', '1b2f96c7-78df-41f6-bcd5-1', 0, 0, '{\"rn\":\"cobaa\"}', NULL),
(23, '2021-12-07 09:20:12.763', '2021-12-07 09:20:12.763', '60b10a2e-d2a5-480f-94e7-9', 0, 0, '{\"rn\":\"lagii\"}', NULL),
(24, '2021-12-07 09:25:01.186', '2021-12-07 09:25:01.186', 'caaf8725-0adc-4ede-8732-d', 0, 0, '{\"rn\":\"cobsss\"}', NULL),
(25, '2021-12-08 10:05:39.392', '2021-12-08 10:05:39.392', '47b9a9ac-4c76-44c2-a434-5', 0, 0, '{\"rn\":\"hilmi-hilmihi\"}', NULL),
(26, '2021-12-09 05:37:16.048', '2021-12-09 05:37:16.048', '8cd329a8-69b0-4734-bffe-e', 0, 0, '{\"rn\":\"hilmihi-robirit\"}', NULL);

-- --------------------------------------------------------

--
-- Struktur dari tabel `subscriptions`
--

DROP TABLE IF EXISTS `subscriptions`;
CREATE TABLE `subscriptions` (
  `id` int(11) NOT NULL,
  `createdate` datetime(3) NOT NULL,
  `updatedate` datetime(3) NOT NULL,
  `deletedate` datetime(3) DEFAULT NULL,
  `userid` int(11) NOT NULL,
  `room` varchar(36) NOT NULL,
  `seqid` int(11) DEFAULT 0,
  `comment` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (json_valid(`comment`))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data untuk tabel `subscriptions`
--

INSERT INTO `subscriptions` (`id`, `createdate`, `updatedate`, `deletedate`, `userid`, `room`, `seqid`, `comment`) VALUES
(1, '2021-12-05 16:26:34.000', '2021-12-05 16:26:34.000', NULL, 1, 'grpAsdsWfGdgr', 0, NULL),
(2, '2021-12-05 17:03:55.000', '2021-12-05 17:03:55.000', NULL, 2, 'grpAsdsWfGdgr', 0, NULL),
(3, '2021-12-07 08:39:09.974', '2021-12-07 08:39:09.974', NULL, 1, '1edbfda7-e4ee-44a6-86be-7', 0, '\"\"'),
(4, '2021-12-07 08:39:09.980', '2021-12-07 08:39:09.980', NULL, 3, '1edbfda7-e4ee-44a6-86be-7', 0, '\"\"'),
(5, '2021-12-07 09:01:30.939', '2021-12-07 09:01:30.939', NULL, 1, '1b2f96c7-78df-41f6-bcd5-1', 0, '\"\"'),
(6, '2021-12-07 09:01:30.943', '2021-12-07 09:01:30.943', NULL, 2, '1b2f96c7-78df-41f6-bcd5-1', 0, '\"\"'),
(7, '2021-12-07 09:20:12.763', '2021-12-07 09:20:12.763', NULL, 1, '60b10a2e-d2a5-480f-94e7-9', 0, '\"\"'),
(8, '2021-12-07 09:20:12.767', '2021-12-07 09:20:12.767', NULL, 3, '60b10a2e-d2a5-480f-94e7-9', 0, '\"\"'),
(9, '2021-12-07 09:20:12.768', '2021-12-07 09:20:12.768', NULL, 2, '60b10a2e-d2a5-480f-94e7-9', 0, '\"\"'),
(10, '2021-12-07 09:25:01.186', '2021-12-07 09:25:01.186', NULL, 1, 'caaf8725-0adc-4ede-8732-d', 0, '\"\"'),
(11, '2021-12-07 09:25:01.190', '2021-12-07 09:25:01.190', NULL, 3, 'caaf8725-0adc-4ede-8732-d', 0, '\"\"'),
(12, '2021-12-08 10:05:39.392', '2021-12-08 10:05:39.392', NULL, 1, '47b9a9ac-4c76-44c2-a434-5', 0, '\"\"'),
(13, '2021-12-08 10:05:39.397', '2021-12-08 10:05:39.397', NULL, 4, '47b9a9ac-4c76-44c2-a434-5', 0, '\"\"'),
(14, '2021-12-09 05:37:16.048', '2021-12-09 05:37:16.048', NULL, 4, '8cd329a8-69b0-4734-bffe-e', 0, '\"\"'),
(15, '2021-12-09 05:37:16.052', '2021-12-09 05:37:16.052', NULL, 2, '8cd329a8-69b0-4734-bffe-e', 0, '\"\"');

-- --------------------------------------------------------

--
-- Struktur dari tabel `users`
--

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` int(11) NOT NULL,
  `username` varchar(32) NOT NULL,
  `createdate` datetime(3) NOT NULL,
  `updatedate` datetime(3) NOT NULL,
  `state` smallint(6) NOT NULL DEFAULT 0,
  `lastseen` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data untuk tabel `users`
--

INSERT INTO `users` (`id`, `username`, `createdate`, `updatedate`, `state`, `lastseen`) VALUES
(1, 'hilmi', '2021-12-02 17:37:45.000', '2021-12-02 17:37:45.000', 1, '2021-12-02 11:37:45'),
(2, 'robirit', '2021-12-02 17:37:45.000', '2021-12-02 17:37:45.000', 1, '2021-12-02 11:37:45'),
(3, 'ardhi', '2021-12-07 11:22:24.000', '2021-12-07 11:22:24.000', 1, '2021-12-07 11:37:45'),
(4, 'hilmihi', '2021-12-08 09:48:04.684', '2021-12-08 09:48:04.684', 1, '2021-12-08 09:48:04'),
(5, 'robit', '2021-12-09 02:11:51.410', '2021-12-09 02:11:51.410', 1, '2021-12-09 02:11:51');

--
-- Indexes for dumped tables
--

--
-- Indeks untuk tabel `messages`
--
ALTER TABLE `messages`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `messages_room_seqid` (`room`,`seqid`);

--
-- Indeks untuk tabel `rooms`
--
ALTER TABLE `rooms`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `room_name` (`name`),
  ADD KEY `room_owner` (`owner`);

--
-- Indeks untuk tabel `subscriptions`
--
ALTER TABLE `subscriptions`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `subscriptions_room_userid` (`room`,`userid`),
  ADD KEY `userid` (`userid`),
  ADD KEY `subscriptions_room` (`room`);

--
-- Indeks untuk tabel `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `auth_uname` (`username`);

--
-- AUTO_INCREMENT untuk tabel yang dibuang
--

--
-- AUTO_INCREMENT untuk tabel `messages`
--
ALTER TABLE `messages`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=23;

--
-- AUTO_INCREMENT untuk tabel `rooms`
--
ALTER TABLE `rooms`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=27;

--
-- AUTO_INCREMENT untuk tabel `subscriptions`
--
ALTER TABLE `subscriptions`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=16;

--
-- AUTO_INCREMENT untuk tabel `users`
--
ALTER TABLE `users`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

--
-- Ketidakleluasaan untuk tabel pelimpahan (Dumped Tables)
--

--
-- Ketidakleluasaan untuk tabel `messages`
--
ALTER TABLE `messages`
  ADD CONSTRAINT `messages_ibfk_1` FOREIGN KEY (`room`) REFERENCES `rooms` (`name`);

--
-- Ketidakleluasaan untuk tabel `subscriptions`
--
ALTER TABLE `subscriptions`
  ADD CONSTRAINT `subscriptions_ibfk_1` FOREIGN KEY (`userid`) REFERENCES `users` (`id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
