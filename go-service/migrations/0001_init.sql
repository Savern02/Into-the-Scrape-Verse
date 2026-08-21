-- Into the Scrape-Verse :: canonical schema
-- Run against Supabase: SQL Editor -> paste -> Run
-- Principle: raw observations are preserved; everything else is derived.

-- ---------------------------------------------------------------
-- Sources
-- ---------------------------------------------------------------

create table if not exists retailers (
    id          bigint generated always as identity primary key,
    slug        text        not null unique,
    name        text        not null,
    homepage    text,
    -- true when this retailer publishes education / institutional pricing
    edu_program boolean     not null default false,
    created_at  timestamptz not null default now()
);

-- One row per Bright Data Collector ID (c_*). This is the join point
-- between the scraper side of the team and the pipeline side.
create table if not exists collectors (
    id           text primary key,               -- e.g. c_8f2a91
    retailer_id  bigint      not null references retailers (id),
    scraper_type text        not null,           -- pdp | discovery | sitemap | search | browser
    target_url   text,
    -- the plain-language field descriptions the scraper was built from.
    -- `bdata scraper heal` re-derives selectors from these, so keep them here.
    field_spec   jsonb       not null default '{}'::jsonb,
    active       boolean     not null default true,
    created_at   timestamptz not null default now()
);

-- ---------------------------------------------------------------
-- Raw layer
-- ---------------------------------------------------------------

create table if not exists raw_ingests (
    id           bigint generated always as identity primary key,
    collector_id text        not null references collectors (id),
    snapshot_id  text,                            -- Bright Data snapshot / run id
    object_key   text        not null,            -- key in object storage
    row_count    int         not null default 0,
    status       text        not null default 'received', -- received|processing|done|failed
    error        text,
    received_at  timestamptz not null default now(),
    processed_at timestamptz,
    unique (collector_id, snapshot_id)
);

-- ---------------------------------------------------------------
-- Canonical layer
-- ---------------------------------------------------------------

create table if not exists products (
    id            bigint generated always as identity primary key,
    canonical_key text        not null unique,    -- normalized brand|name|pack, used for matching
    name          text        not null,
    brand         text,
    category      text,                            -- paper | writing | art | janitorial | tech | furniture
    unit_type     text,                            -- each | ream | case | pack | box
    unit_count    numeric,                         -- 500 (sheets), 12 (pens), 24 (crayons)
    created_at    timestamptz not null default now()
);

create index if not exists products_category_idx on products (category);

create table if not exists offers (
    id           bigint generated always as identity primary key,
    product_id   bigint  not null references products (id) on delete cascade,
    retailer_id  bigint  not null references retailers (id),
    sku          text,
    url          text    not null,
    edu_eligible boolean not null default false,
    min_bulk_qty int,
    unique (retailer_id, url)
);

-- Append-only. Never update a row here; a new scrape is a new observation.
create table if not exists price_observations (
    id                   bigint generated always as identity primary key,
    offer_id             bigint      not null references offers (id) on delete cascade,
    raw_ingest_id        bigint      references raw_ingests (id),
    price_cents          bigint      not null,
    list_price_cents     bigint,
    currency             text        not null default 'USD',
    in_stock             boolean,
    qty_break            int         not null default 1,   -- bulk tier: 1, 12, 50, 100...
    price_per_unit_cents numeric,                            -- derived, for cross-pack comparison
    observed_at          timestamptz not null,
    unique (offer_id, qty_break, observed_at)
);

create index if not exists price_obs_offer_time_idx
    on price_observations (offer_id, observed_at desc);

-- ---------------------------------------------------------------
-- Self-healing evidence  <-- judging criterion 05 lives here
-- ---------------------------------------------------------------

create table if not exists collector_health (
    id               bigint generated always as identity primary key,
    collector_id     text        not null references collectors (id),
    raw_ingest_id    bigint      references raw_ingests (id),
    rows_returned    int         not null,
    rows_dropped     int         not null,
    -- {"price": 0.98, "title": 0.0} -- fraction of rows where the field came back empty
    field_null_rates jsonb       not null default '{}'::jsonb,
    verdict          text        not null,        -- healthy | degraded | broken
    heal_command     text,                         -- the exact bdata command we generated
    heal_output      text,
    -- awaiting_approval | done | failed | rejected.
    -- `bdata scraper heal` stops at a human approval gate; a heal is NOT live
    -- until `bdata scraper approve` moves it to done.
    heal_status      text,
    healed_at        timestamptz,
    checked_at       timestamptz not null default now()
);

create index if not exists collector_health_time_idx
    on collector_health (collector_id, checked_at desc);

-- ---------------------------------------------------------------
-- Convenience view for the Next.js side
-- ---------------------------------------------------------------

create or replace view current_prices as
select distinct on (o.id)
    o.id            as offer_id,
    p.id            as product_id,
    p.name,
    p.brand,
    p.category,
    p.unit_type,
    p.unit_count,
    r.slug          as retailer,
    o.url,
    o.edu_eligible,
    po.price_cents,
    po.price_per_unit_cents,
    po.qty_break,
    po.in_stock,
    po.observed_at
from offers o
    join products p           on p.id = o.product_id
    join retailers r          on r.id = o.retailer_id
    join price_observations po on po.offer_id = o.id
where po.qty_break = 1
order by o.id, po.observed_at desc;
