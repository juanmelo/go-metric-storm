CREATE TABLE `go_storm_log` (
  `go_storm_log_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `topology_id` varchar(128) DEFAULT NULL,
  `workers` int(11) DEFAULT NULL,
  `bolt_parse` double(8,3) DEFAULT NULL,
  `bolt_load_config` double(8,3) DEFAULT NULL,
  `bolt_group` double(8,3) DEFAULT NULL,
  `bolt_complete` double(8,3) DEFAULT NULL,
  `bolt_persist` double(8,3) DEFAULT NULL,
  `bolt_entity_sync` double(8,3) DEFAULT NULL,
  `created_date` datetime DEFAULT NULL,
  PRIMARY KEY (`go_storm_log_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5363 DEFAULT CHARSET=utf8;

CREATE TABLE `go_storm_log_detail` (
  `go_storm_log_detail` bigint(20) NOT NULL AUTO_INCREMENT,
  `go_storm_log_id` bigint(20) DEFAULT NULL,
  `bolt_id` varchar(64) DEFAULT NULL,
  `executors` int(11) DEFAULT NULL,
  `tasks` int(11) DEFAULT NULL,
  `emitted` bigint(20) DEFAULT NULL,
  `transferred` bigint(20) DEFAULT NULL,
  `acked` bigint(20) DEFAULT NULL,
  `failed` bigint(20) DEFAULT NULL,
  `capacity` double(8,3) DEFAULT NULL,
  `execute_latency` double(8,3) DEFAULT NULL,
  PRIMARY KEY (`go_storm_log_detail`)
) ENGINE=InnoDB AUTO_INCREMENT=32137 DEFAULT CHARSET=utf8;

CREATE TABLE `go_storm_log_lag` (
  `go_storm_log_lag_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `go_storm_log_id` bigint(20) DEFAULT NULL,
  `spout_id` varchar(128) DEFAULT NULL,
  `topic` varchar(128) DEFAULT NULL,
  `topic_partition` varchar(16) DEFAULT NULL,
  `consumer_committed_offset` bigint(20) DEFAULT NULL,
  `log_head_offset` bigint(20) DEFAULT NULL,
  `lag` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`go_storm_log_lag_id`)
) ENGINE=InnoDB AUTO_INCREMENT=134490 DEFAULT CHARSET=utf8;
