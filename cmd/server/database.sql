-- Updated-at trigger helper
create or replace function set_updated_at()
    returns trigger as $$
begin
    new.updated_at = now();
    return new;
end;
$$ language plpgsql;


-- terms table (core entries)
create table if not exists terms (
                                     id serial primary key,
                                     term_en varchar(255) not null,
                                     definition_en text,
                                     context text,
                                     created_at timestamp not null default now(),
                                     updated_at timestamp not null default now()
);


create index if not exists idx_terms_term_en on terms (term_en);
create index if not exists idx_terms_created_at on terms (created_at);


create trigger trg_terms_updated_at
    before update on terms
    for each row execute procedure set_updated_at();


-- categories table (e.g. programming, security, networking)
create table if not exists categories (
                                          id serial primary key,
                                          name varchar(100) not null unique
);


-- junction table term_categories (many-to-many)
create table if not exists term_categories (
                                               term_id int not null references terms(id) on delete cascade,
                                               category_id int not null references categories(id) on delete cascade,
                                               primary key (term_id, category_id)
);


create index if not exists idx_term_categories_term on term_categories (term_id);
create index if not exists idx_term_categories_category on term_categories (category_id);


-- translations table (fa equivalents + sources)
create table if not exists translations (
                                            id serial primary key,
                                            term_id int not null references terms(id) on delete cascade,
                                            term_fa varchar(255),
                                            definition_fa text,
                                            source varchar(255)
);


create index if not exists idx_translations_term on translations (term_id);


-- related_terms table (self-referencing relations)
create table if not exists related_terms (
                                             id serial primary key,
                                             term_id int not null references terms(id) on delete cascade,
                                             related_term_id int not null references terms(id) on delete cascade,
                                             relation_type varchar(50)
);


create index if not exists idx_related_terms_term on related_terms (term_id);
create index if not exists idx_related_terms_related on related_terms (related_term_id);


-- Seed a few common categories (idempotent)
insert into categories (name)
select * from (values ('programming'),('security'),('networking'),('operating systems'),('databases')) as v(name)
where not exists (
    select 1 from categories c where c.name = v.name
);