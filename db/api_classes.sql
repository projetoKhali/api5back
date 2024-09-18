-- Created by Vertabelo (http://vertabelo.com)
-- Last modification date: 2024-09-18 16:12:11.942

-- tables
-- Table: candidate
CREATE TABLE candidate (
    cd_id int  NOT NULL,
    cd_name varchar  NOT NULL,
    cd_email varchar  NOT NULL,
    cd_phone varchar  NOT NULL,
    cd_status int  NOT NULL DEFAULT 1,
    cd_score int  NOT NULL,
    cd_last_update timestamp  NOT NULL,
    CONSTRAINT candidate_pk PRIMARY KEY (cd_id)
);

-- Table: candidate_evaluation
CREATE TABLE candidate_evaluation (
    ce_id int  NOT NULL,
    cd_id int  NOT NULL,
    pc_id int  NOT NULL,
    ce_evaluation_criteria varchar  NOT NULL,
    ce_score int  NOT NULL,
    ce_evaluation_date timestamp  NOT NULL,
    usr_id int  NOT NULL,
    CONSTRAINT candidate_evaluation_pk PRIMARY KEY (ce_id)
);

-- Table: department
CREATE TABLE department (
    dp_id int  NOT NULL,
    dp_name varchar  NOT NULL,
    dp_description varchar  NOT NULL,
    CONSTRAINT department_pk PRIMARY KEY (dp_id)
);

-- Table: feedback
CREATE TABLE feedback (
    fd_id int  NOT NULL,
    cd_id int  NULL,
    vc_id int  NOT NULL,
    fd_date timestamp  NOT NULL DEFAULT now(),
    fd_type int  NOT NULL,
    fd_content varchar  NOT NULL,
    usr_id int  NOT NULL,
    CONSTRAINT feedback_pk PRIMARY KEY (fd_id)
);

-- Table: hiring
CREATE TABLE hiring (
    hr_id int  NOT NULL,
    cd_id int  NOT NULL,
    vc_id int  NOT NULL,
    hr_contract_start_date date  NOT NULL,
    hr_initial_salary decimal(8,2)  NOT NULL,
    hr_employment_classifications varchar  NOT NULL,
    hr_offer_acceptance_date timestamp  NOT NULL,
    CONSTRAINT hiring_pk PRIMARY KEY (hr_id)
);

-- Table: interview
CREATE TABLE interview (
    iw_id int  NOT NULL,
    cd_id int  NULL,
    vc_id int  NOT NULL,
    iw_date timestamp  NOT NULL,
    usr_id int  NOT NULL,
    iw_result varchar  NOT NULL,
    iw_observation varchar  NULL,
    iw_schreduling_date timestamp  NOT NULL DEFAULT now(),
    iw_conclusion_date timestamp  NOT NULL,
    CONSTRAINT interview_pk PRIMARY KEY (iw_id)
);

-- Table: process
CREATE TABLE process (
    pc_id int  NOT NULL,
    pc_title varchar  NOT NULL,
    pc_initial_date timestamp  NOT NULL,
    pc_expected_finish_date timestamp  NOT NULL,
    pc_finish_date timestamp  NOT NULL,
    pc_status int  NOT NULL DEFAULT 1,
    usr_id int  NOT NULL,
    pc_description varchar  NULL,
    dp_id int  NOT NULL,
    CONSTRAINT process_pk PRIMARY KEY (pc_id)
);

-- Table: process_history
CREATE TABLE process_history (
    ph_id int  NOT NULL,
    ph_type_action varchar  NOT NULL,
    pc_id int  NOT NULL,
    usr_id int  NOT NULL,
    ph_insert_date timestamp  NOT NULL DEFAULT now(),
    CONSTRAINT process_history_pk PRIMARY KEY (ph_id)
);

-- Table: requirement
CREATE TABLE requirement (
    req_id int  NOT NULL,
    req_name varchar  NOT NULL,
    req_proficiency_level varchar  NULL,
    CONSTRAINT requirement_pk PRIMARY KEY (req_id)
);

-- Table: user
CREATE TABLE "user" (
    usr_id int  NOT NULL,
    usr_name varchar  NOT NULL,
    usr_email varchar  NOT NULL,
    usr_password varchar  NOT NULL,
    usr_ocupation varchar  NOT NULL,
    usr_qtd_feedback int  NOT NULL,
    dp_id int  NULL,
    CONSTRAINT user_pk PRIMARY KEY (usr_id)
);

-- Table: vacancy
CREATE TABLE vacancy (
    vc_id int  NOT NULL,
    pc_id int  NOT NULL,
    vc_title varchar  NOT NULL,
    vc_num_positions int  NOT NULL,
    vc_status int  NOT NULL DEFAULT 1,
    vc_location varchar  NOT NULL,
    usr_id int  NOT NULL,
    vc_opening_date timestamp  NOT NULL,
    vc_closing_date int  NOT NULL,
    CONSTRAINT vacancy_pk PRIMARY KEY (vc_id)
);

-- Table: vacancy_candidate
CREATE TABLE vacancy_candidate (
    cd_id int  NOT NULL,
    vc_id int  NOT NULL,
    vc_cd_insert_date timestamp  NOT NULL DEFAULT now(),
    CONSTRAINT vacancy_candidate_pk PRIMARY KEY (cd_id,vc_id)
);

-- Table: vacancy_requirement
CREATE TABLE vacancy_requirement (
    vc_rec_vc_id int  NOT NULL,
    vc_req_req_id int  NOT NULL,
    vc_req_insertdate timestamp  NOT NULL DEFAULT now(),
    CONSTRAINT vacancy_requirement_pk PRIMARY KEY (vc_rec_vc_id,vc_req_req_id)
);

-- foreign keys
-- Reference: HIring_Candidate (table: hiring)
ALTER TABLE hiring ADD CONSTRAINT HIring_Candidate
    FOREIGN KEY (cd_id)
    REFERENCES candidate (cd_id)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: HIring_Vacancy (table: hiring)
ALTER TABLE hiring ADD CONSTRAINT HIring_Vacancy
    FOREIGN KEY (vc_id)
    REFERENCES vacancy (vc_id)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: Interviel_Candidate (table: interview)
ALTER TABLE interview ADD CONSTRAINT Interviel_Candidate
    FOREIGN KEY (cd_id)
    REFERENCES candidate (cd_id)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: Interviel_User (table: interview)
ALTER TABLE interview ADD CONSTRAINT Interviel_User
    FOREIGN KEY (usr_id)
    REFERENCES "user" (usr_id)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: Vacancy_Process (table: vacancy)
ALTER TABLE vacancy ADD CONSTRAINT Vacancy_Process
    FOREIGN KEY (pc_id)
    REFERENCES process (pc_id)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: Vacancy_User (table: vacancy)
ALTER TABLE vacancy ADD CONSTRAINT Vacancy_User
    FOREIGN KEY (usr_id)
    REFERENCES "user" (usr_id)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: candidate_evaluation_candidate (table: candidate_evaluation)
ALTER TABLE candidate_evaluation ADD CONSTRAINT candidate_evaluation_candidate
    FOREIGN KEY (cd_id)
    REFERENCES candidate (cd_id)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: candidate_evaluation_process (table: candidate_evaluation)
ALTER TABLE candidate_evaluation ADD CONSTRAINT candidate_evaluation_process
    FOREIGN KEY (pc_id)
    REFERENCES process (pc_id)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: candidate_evaluation_user (table: candidate_evaluation)
ALTER TABLE candidate_evaluation ADD CONSTRAINT candidate_evaluation_user
    FOREIGN KEY (usr_id)
    REFERENCES "user" (usr_id)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: feedback_candidate (table: feedback)
ALTER TABLE feedback ADD CONSTRAINT feedback_candidate
    FOREIGN KEY (cd_id)
    REFERENCES candidate (cd_id)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: feedback_user (table: feedback)
ALTER TABLE feedback ADD CONSTRAINT feedback_user
    FOREIGN KEY (usr_id)
    REFERENCES "user" (usr_id)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: feedback_vacancy (table: feedback)
ALTER TABLE feedback ADD CONSTRAINT feedback_vacancy
    FOREIGN KEY (vc_id)
    REFERENCES vacancy (vc_id)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: process_department (table: process)
ALTER TABLE process ADD CONSTRAINT process_department
    FOREIGN KEY (dp_id)
    REFERENCES department (dp_id)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: process_history_process (table: process_history)
ALTER TABLE process_history ADD CONSTRAINT process_history_process
    FOREIGN KEY (pc_id)
    REFERENCES process (pc_id)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: process_history_user (table: process_history)
ALTER TABLE process_history ADD CONSTRAINT process_history_user
    FOREIGN KEY (usr_id)
    REFERENCES "user" (usr_id)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: process_user (table: process)
ALTER TABLE process ADD CONSTRAINT process_user
    FOREIGN KEY (usr_id)
    REFERENCES "user" (usr_id)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: user_department (table: user)
ALTER TABLE "user" ADD CONSTRAINT user_department
    FOREIGN KEY (dp_id)
    REFERENCES department (dp_id)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: vacancy_candidate_candidate (table: vacancy_candidate)
ALTER TABLE vacancy_candidate ADD CONSTRAINT vacancy_candidate_candidate
    FOREIGN KEY (cd_id)
    REFERENCES candidate (cd_id)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: vacancy_candidate_vacancy (table: vacancy_candidate)
ALTER TABLE vacancy_candidate ADD CONSTRAINT vacancy_candidate_vacancy
    FOREIGN KEY (vc_id)
    REFERENCES vacancy (vc_id)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: vacancy_requirement_requirement (table: vacancy_requirement)
ALTER TABLE vacancy_requirement ADD CONSTRAINT vacancy_requirement_requirement
    FOREIGN KEY (vc_req_req_id)
    REFERENCES requirement (req_id)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- Reference: vacancy_requirement_vacancy (table: vacancy_requirement)
ALTER TABLE vacancy_requirement ADD CONSTRAINT vacancy_requirement_vacancy
    FOREIGN KEY (vc_rec_vc_id)
    REFERENCES vacancy (vc_id)  
    NOT DEFERRABLE 
    INITIALLY IMMEDIATE
;

-- End of file.