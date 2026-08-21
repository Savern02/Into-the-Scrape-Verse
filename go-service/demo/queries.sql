-- 1. Scraper health over time: detection, verdict, generated heal command
select collector_id, verdict, rows_returned, rows_dropped,
       field_null_rates, heal_command, checked_at
from collector_health order by checked_at;

-- 2. The calculator: same product category, wildly different unit economics
select name, brand, price_cents/100.0 as price,
       unit_count, round(price_per_unit_cents::numeric,2) as cents_per_unit
from current_prices
where unit_count > 1
order by cents_per_unit;

-- 3. Raw preserved for reprocessing without re-scraping
select collector_id, snapshot_id, object_key, row_count, status
from raw_ingests order by received_at desc;
