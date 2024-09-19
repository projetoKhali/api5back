-- Created by Vertabelo (http://vertabelo.com)
-- Last modification date: 2024-09-18 16:16:42.857

-- tables
-- Table: dim_datetime
CREATE TABLE dim_datetime (
    dim_datetime_id int  NOT NULL,
    dim_datetime_date date  NOT NULL,
    dim_datetime_year int  NOT NULL,
    dim_datetime_month int  NOT NULL,
    dim_datetime_weekday int  NOT NULL,
    dim_datetime_day int  NOT NULL,
    dim_datetime_hour int  NOT NULL,
    dim_datetime_minute int  NOT NULL,
    dim_datetime_second int  NOT NULL,
    CONSTRAINT dim_datetime_pk PRIMARY KEY (dim_datetime_id)
);

-- Table: dim_process
CREATE TABLE dim_process (
    dim_pc_id int  NOT NULL,
    dim_pc_title varchar  NOT NULL,
    dim_pc_initial_date timestamp  NOT NULL,
    dim_pc_finish_date timestamp  NOT NULL,
    dim_pc_status int  NOT NULL DEFAULT 1,
    dim_usr_id int  NOT NULL,
    dim_pc_description varchar  NULL,
    CONSTRAINT dim_process_pk PRIMARY KEY (dim_pc_id)
);

-- Table: dim_user
CREATE TABLE dim_user (
    dim_usr_id int  NOT NULL,
    dim_usr_name varchar  NOT NULL,
    dim_usr_ocupation varchar  NOT NULL,
    CONSTRAINT dim_user_pk PRIMARY KEY (dim_usr_id)
);

-- Table: dim_vacancy
CREATE TABLE dim_vacancy (
    vc_id int  NOT NULL,
    vc_title varchar  NOT NULL,
    vc_num_positions int  NOT NULL,
    req_id int  NOT NULL,
    vc_status int  NOT NULL DEFAULT 1,
    vc_location varchar  NOT NULL,
    usr_id int  NOT NULL,
    vc_opening_date timestamp  NOT NULL,
    vc_closing_date int  NOT NULL,
    CONSTRAINT dim_vacancy_pk PRIMARY KEY (vc_id)
);

-- Table: fact_hiring_process
CREATE TABLE fact_hiring_process (
    fac_id int  NOT NULL,
    dim_process_dim_pc_id int  NOT NULL,
    dim_vacancy_vc_id int  NOT NULL,
    dim_user_dim_usr_id int  NOT NULL,
    dim_date_dim_date_id int  NOT NULL,
    met_total_candidates_applied int  NOT NULL,
    met_total_candidates_interviewed int  NOT NULL,
    met_total_candidates_hired int  NOT NULL,
    met_sum_duration_hiring_proces int  NOT NULL,
    met_sum_salary_initial int  NOT NULL,
    met_total_feedback_positive int  NOT NULL,
    met_total_neutral int  NOT NULL,
    met_total_negative int  NOT NULL,
    CONSTRAINT fact_hiring_process_pk PRIMARY KEY (fac_id)
);

-- foreign keys
-- Reference: fact_hiring_process_dim_date (table: fact_hiring_process)
ALTER TABLE fact_hiring_process ADD CONSTRAINT fact_hiring_process_dim_date
    FOREIGN KEY (dim_date_dim_date_id)
    REFERENCES dim_datetime (dim_datetime_id)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: fact_hiring_process_dim_process (table: fact_hiring_process)
ALTER TABLE fact_hiring_process ADD CONSTRAINT fact_hiring_process_dim_process
    FOREIGN KEY (dim_process_dim_pc_id)
    REFERENCES dim_process (dim_pc_id)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: fact_hiring_process_dim_user (table: fact_hiring_process)
ALTER TABLE fact_hiring_process ADD CONSTRAINT fact_hiring_process_dim_user
    FOREIGN KEY (dim_user_dim_usr_id)
    REFERENCES dim_user (dim_usr_id)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: fact_hiring_process_dim_vacancy (table: fact_hiring_process)
ALTER TABLE fact_hiring_process ADD CONSTRAINT fact_hiring_process_dim_vacancy
    FOREIGN KEY (dim_vacancy_vc_id)
    REFERENCES dim_vacancy (vc_id)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- End of file.

