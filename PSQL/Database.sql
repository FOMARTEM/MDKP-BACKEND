BEGIN;


CREATE TABLE IF NOT EXISTS public."Employee"
(
    id bigserial NOT NULL,
    "LastName" character varying(64) COLLATE pg_catalog."default" NOT NULL,
    "FirstName" character varying(64) COLLATE pg_catalog."default" NOT NULL,
    "MiddleName" character varying(64) COLLATE pg_catalog."default" NOT NULL,
    "Email" character varying(255) COLLATE pg_catalog."default" NOT NULL,
    "Phone" character varying(32) COLLATE pg_catalog."default" NOT NULL,
    "BirthDate" date NOT NULL,
    "PasswordHash" character varying(511) COLLATE pg_catalog."default" NOT NULL,
    "IsActive" boolean NOT NULL,
    "Position_ID" integer NOT NULL,
    CONSTRAINT "Employee_pkey" PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public."Material"
(
    id bigserial NOT NULL,
    "FileName" character varying(64) COLLATE pg_catalog."default" NOT NULL,
    "Extension" character varying(4) COLLATE pg_catalog."default" NOT NULL,
    "Description" text COLLATE pg_catalog."default" NOT NULL,
    "EmployeeID" bigint,
    "TaskID" bigint,
    CONSTRAINT "Material_pkey" PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public."Position"
(
    id serial NOT NULL,
    "Name" character(63) COLLATE pg_catalog."default" NOT NULL,
    "Description" character varying(255) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT "Position_pkey" PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public."Revision"
(
    id bigserial NOT NULL,
    "Description" text COLLATE pg_catalog."default" NOT NULL,
    "CreationDate" date NOT NULL,
    "EmployeeID" bigint NOT NULL,
    "VersionID" bigint NOT NULL,
    "StatusID" integer NOT NULL,
    CONSTRAINT "Revision_pkey" PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public."Status"
(
    id serial NOT NULL,
    "Name" character varying(63) COLLATE pg_catalog."default" NOT NULL,
    "Description" character varying(127) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT "Status_pkey" PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public."Task"
(
    id bigserial NOT NULL,
    "Title" character varying(127) COLLATE pg_catalog."default" NOT NULL,
    "Description" text COLLATE pg_catalog."default" NOT NULL,
    "CreateDate" date NOT NULL,
    "DeadlineDate" date NOT NULL,
    "ReadyDate" date,
    "Priority" integer NOT NULL,
    "CreatorId" bigint NOT NULL,
    "EditorID" bigint NOT NULL,
    "AuthorID" bigint NOT NULL,
    "StatusID" integer NOT NULL,
    CONSTRAINT "Task_pkey" PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public."Version"
(
    id bigserial NOT NULL,
    "VersionNumber" integer NOT NULL,
    "CreationDate" date NOT NULL,
    "Title" character varying(127) COLLATE pg_catalog."default" NOT NULL,
    "Description" text COLLATE pg_catalog."default" NOT NULL,
    "EmployeeID" bigint NOT NULL,
    "MaterialID" bigint,
    "TaskID" bigint NOT NULL,
    CONSTRAINT "Version_pkey" PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.audit_log
(
    id bigserial NOT NULL,
    employee_id bigint,
    action character varying(64) COLLATE pg_catalog."default" NOT NULL,
    entity_type character varying(64) COLLATE pg_catalog."default" NOT NULL,
    entity_id bigint,
    details jsonb NOT NULL DEFAULT '{}'::jsonb,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT audit_log_pkey PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS public."Employee"
    ADD CONSTRAINT "Position" FOREIGN KEY ("Position_ID")
    REFERENCES public."Position" (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;

COMMENT ON CONSTRAINT "Position" ON public."Employee"
    IS 'Связь для определения прав (Админ, редактор, автор, руководитель и тп)';



ALTER TABLE IF EXISTS public."Material"
    ADD CONSTRAINT "Employee ID" FOREIGN KEY ("EmployeeID")
    REFERENCES public."Employee" (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;


ALTER TABLE IF EXISTS public."Material"
    ADD CONSTRAINT "Task ID" FOREIGN KEY ("TaskID")
    REFERENCES public."Task" (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;
CREATE INDEX IF NOT EXISTS ix_material_task
    ON public."Material"("TaskID");


ALTER TABLE IF EXISTS public."Revision"
    ADD CONSTRAINT "Employee ID" FOREIGN KEY ("EmployeeID")
    REFERENCES public."Employee" (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;


ALTER TABLE IF EXISTS public."Revision"
    ADD CONSTRAINT "Status ID" FOREIGN KEY ("StatusID")
    REFERENCES public."Status" (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;


ALTER TABLE IF EXISTS public."Revision"
    ADD CONSTRAINT "Version ID" FOREIGN KEY ("VersionID")
    REFERENCES public."Version" (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;


ALTER TABLE IF EXISTS public."Task"
    ADD CONSTRAINT "Author ID" FOREIGN KEY ("AuthorID")
    REFERENCES public."Employee" (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;
CREATE INDEX IF NOT EXISTS ix_task_author
    ON public."Task"("AuthorID");


ALTER TABLE IF EXISTS public."Task"
    ADD CONSTRAINT "Creator ID" FOREIGN KEY ("CreatorId")
    REFERENCES public."Employee" (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;
CREATE INDEX IF NOT EXISTS ix_task_creator
    ON public."Task"("CreatorId");


ALTER TABLE IF EXISTS public."Task"
    ADD CONSTRAINT "Editor ID" FOREIGN KEY ("EditorID")
    REFERENCES public."Employee" (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;
CREATE INDEX IF NOT EXISTS ix_task_editor
    ON public."Task"("EditorID");


ALTER TABLE IF EXISTS public."Task"
    ADD CONSTRAINT "Status ID" FOREIGN KEY ("StatusID")
    REFERENCES public."Status" (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;


ALTER TABLE IF EXISTS public."Version"
    ADD CONSTRAINT "Employee ID" FOREIGN KEY ("EmployeeID")
    REFERENCES public."Employee" (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;


ALTER TABLE IF EXISTS public."Version"
    ADD CONSTRAINT "Material ID" FOREIGN KEY ("MaterialID")
    REFERENCES public."Material" (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;


ALTER TABLE IF EXISTS public."Version"
    ADD CONSTRAINT "Task ID" FOREIGN KEY ("TaskID")
    REFERENCES public."Task" (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;


ALTER TABLE IF EXISTS public.audit_log
    ADD CONSTRAINT audit_log_employee_id_fkey FOREIGN KEY (employee_id)
    REFERENCES public."Employee" (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;

END;