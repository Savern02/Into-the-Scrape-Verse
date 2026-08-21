-- Into the Scrape-Verse :: seed + supply list layer
-- Run AFTER 0001_init.sql. Paste the CONTENTS of this file into the
-- Supabase SQL editor, not the filename.

-- ===============================================================
-- Part 1: the supply-list layer -- this is the calculator
-- ===============================================================
--
-- A district publishes "3rd grade needs 24 pencils, 4 glue sticks, 1 ream
-- of copy paper." That list is the unit a parent or a purchasing officer
-- actually cares about. Everything else in this database exists to answer
-- "what does that list cost, and where."

create table if not exists supply_lists (
    id          bigint generated always as identity primary key,
    district    text not null,          -- 'Broken Arrow Public Schools'
    school      text,                   -- null = district-wide list
    grade       text not null,          -- 'K', '3', '9-12', 'Art'
    school_year text not null,          -- '2026-2027'
    source_url  text,                   -- where you got it; cite it in the UI
    created_at  timestamptz not null default now()
);

-- Postgres does not allow expressions inside a table-level UNIQUE constraint,
-- so the coalesce goes in a unique INDEX instead. Same effect: two lists for
-- the same district+grade+year collide even when school is null.
create unique index if not exists supply_lists_ident_idx
    on supply_lists (district, coalesce(school, ''), grade, school_year);

create table if not exists supply_list_items (
    id             bigint generated always as identity primary key,
    supply_list_id bigint  not null references supply_lists (id) on delete cascade,
    -- verbatim from the district's list, e.g. '24 #2 pencils, sharpened'
    raw_text       text    not null,
    -- parsed
    description    text    not null,    -- '#2 pencils'
    quantity       numeric not null default 1,
    unit_type      text,                -- 'each' | 'ream' | 'box'
    category       text,
    -- filled in once you match the line to a real product
    product_id     bigint  references products (id),
    match_notes    text,
    position       int     not null default 0
);

create index if not exists sli_list_idx on supply_list_items (supply_list_id, position);

-- What does this list cost at each retailer that stocks every item?
-- Retailers missing an item are still shown, with a coverage figure, because
-- "cheapest but you have to shop somewhere else for 3 things" matters.
create or replace view supply_list_costs as
select
    sl.id                       as supply_list_id,
    sl.district,
    sl.grade,
    sl.school_year,
    cp.retailer,
    count(*)                    as items_matched,
    (select count(*) from supply_list_items x
      where x.supply_list_id = sl.id) as items_total,
    round(
        sum(cp.price_cents * sli.quantity) / 100.0, 2
    )                           as total_dollars,
    bool_and(coalesce(cp.in_stock, true)) as all_in_stock
from supply_lists sl
    join supply_list_items sli on sli.supply_list_id = sl.id
    join current_prices cp     on cp.product_id = sli.product_id
group by sl.id, sl.district, sl.grade, sl.school_year, cp.retailer;

-- ===============================================================
-- Part 2: seed retailers
-- ===============================================================
--
-- Education-focused suppliers, chosen because they sit in the long tail that
-- Scraper Studio is for. Before you build a scraper for any of these, confirm
-- it isn't already covered:
--     brightdata pipelines list | grep -i <name>
-- If it shows up there, drop it and pick another.

insert into retailers (slug, name, homepage, edu_program) values
    ('discount-school-supply', 'Discount School Supply', 'https://www.discountschoolsupply.com', true),
    ('really-good-stuff',      'Really Good Stuff',      'https://www.reallygoodstuff.com',     true),
    ('kaplan-early-learning',  'Kaplan Early Learning',  'https://www.kaplanco.com',            true),
    ('nasco-education',        'Nasco Education',        'https://www.enasco.com',              true),
    ('demco',                  'Demco',                  'https://www.demco.com',               true)
on conflict (slug) do nothing;

-- ===============================================================
-- Part 3: register your collector
-- ===============================================================
--
-- Replace c_REPLACE_ME with the real ID from:
--     brightdata scraper create <url> "<description>" -o create.json
--     jq -r '.collector_id' create.json
--
-- field_spec is NOT documentation. `brightdata scraper heal` rewrites the
-- extraction from plain-language field descriptions, and Bright Data's docs
-- say it plainly: vague prompts produce vague heals. Write each one the way
-- you'd explain the field to someone who can't see the page.

insert into collectors (id, retailer_id, scraper_type, target_url, field_spec)
select
    'c_REPLACE_ME',
    r.id,
    'pdp',
    'https://www.discountschoolsupply.com/REPLACE-WITH-A-REAL-CATEGORY-PAGE',
    '{
       "title":        "the full product name as displayed in the page heading",
       "brand":        "the manufacturer or brand name, not the retailer",
       "price":        "the current selling price in US dollars; if a sale price and a struck-through price both appear, this is the sale price",
       "list_price":   "the struck-through original price, only when one is shown",
       "sku":          "the retailer item number or product code",
       "availability": "whether the item is in stock, out of stock, or backordered",
       "url":          "the canonical product page URL"
     }'::jsonb
from retailers r
where r.slug = 'discount-school-supply'
on conflict (id) do nothing;

-- ===============================================================
-- Part 4: a supply list to develop against
-- ===============================================================
-- Replace with a real published list and put its URL in source_url so the
-- app can cite it. Judges notice when data has a provenance trail.

insert into supply_lists (district, school, grade, school_year, source_url)
values ('Broken Arrow Public Schools', null, '3', '2026-2027', null)
on conflict do nothing;

-- Re-runnable: the delete clears any prior seed for this exact list first,
-- so pasting this file twice does not double every line item.
delete from supply_list_items sli
 using supply_lists sl
 where sli.supply_list_id = sl.id
   and sl.district = 'Broken Arrow Public Schools'
   and sl.grade = '3'
   and sl.school_year = '2026-2027';

insert into supply_list_items
    (supply_list_id, raw_text, description, quantity, unit_type, category, position)
select sl.id, v.raw_text, v.description, v.quantity, v.unit_type, v.category, v.position
from supply_lists sl,
     (values
        ('24 #2 pencils, sharpened',      '#2 pencils',            24, 'each',  'writing',    1),
        ('4 large glue sticks',           'glue sticks, large',     4, 'each',  'art',        2),
        ('1 pack wide-ruled paper',       'wide-ruled filler paper',1, 'pack',  'paper',      3),
        ('2 boxes 24-count crayons',      'crayons, 24 count',      2, 'box',   'art',        4),
        ('1 pair blunt-tip scissors',     'blunt-tip scissors',     1, 'each',  'art',        5),
        ('3 two-pocket folders',          'two-pocket folders',     3, 'each',  'paper',      6),
        ('1 ream copy paper, 500 sheets', 'copy paper, 500 sheets', 1, 'ream',  'paper',      7),
        ('2 boxes tissues',               'facial tissues',         2, 'box',   'janitorial', 8)
     ) as v(raw_text, description, quantity, unit_type, category, position)
where sl.district = 'Broken Arrow Public Schools'
  and sl.grade = '3'
  and sl.school_year = '2026-2027';

-- ===============================================================
-- Verify
-- ===============================================================
-- select count(*) from retailers;          -- 5
-- select count(*) from supply_list_items;  -- 8
-- select * from supply_list_costs;         -- empty until products are matched
