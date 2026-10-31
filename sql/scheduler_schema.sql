CREATE TABLE go_job_trigger (
 job_id BIGINT NOT NULL PRIMARY KEY,
 next_fire BIGINT NOT NULL
) ENGINE=InnoDB COMMENT='Go 调度器的下一次触发时间';
