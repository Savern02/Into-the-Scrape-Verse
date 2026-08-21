import { createClient } from "@/lib/supabase/server";

/**
 * The read layer for the price data.
 *
 * Everything here queries VIEWS, not tables. `current_prices` and
 * `supply_list_costs` are the contract between the Go pipeline and this app:
 * the pipeline can restructure how it stores observations without breaking a
 * single component here, as long as those two views keep their shape.
 */

export type CurrentPrice = {
  offer_id: number;
  product_id: number;
  name: string;
  brand: string | null;
  category: string | null;
  unit_type: string | null;
  unit_count: number | null;
  retailer: string;
  url: string;
  edu_eligible: boolean;
  price_cents: number;
  price_per_unit_cents: number | null;
  qty_break: number;
  in_stock: boolean | null;
  observed_at: string;
};

export type CollectorHealth = {
  id: number;
  collector_id: string;
  rows_returned: number;
  rows_dropped: number;
  field_null_rates: Record<string, number>;
  verdict: "healthy" | "degraded" | "broken";
  heal_command: string | null;
  heal_status: string | null;
  healed_at: string | null;
  checked_at: string;
};

export type SupplyListCost = {
  supply_list_id: number;
  district: string;
  grade: string;
  school_year: string;
  retailer: string;
  items_matched: number;
  items_total: number;
  total_dollars: number;
  all_in_stock: boolean;
};

/** Every product with a live price, cheapest per unit first. */
export async function getCurrentPrices(limit = 100): Promise<CurrentPrice[]> {
  const supabase = await createClient();
  const { data, error } = await supabase
    .from("current_prices")
    .select("*")
    .order("price_per_unit_cents", { ascending: true, nullsFirst: false })
    .limit(limit);

  if (error) {
    // Surface it rather than rendering an empty page that looks like "no data".
    // The most common cause is a missing RLS policy — see 0003_rls.sql.
    console.error("getCurrentPrices:", error.message);
    return [];
  }
  return data ?? [];
}

/** Multi-pack products only — the ones where per-unit cost actually differs. */
export async function getComparablePrices(limit = 40): Promise<CurrentPrice[]> {
  const supabase = await createClient();
  const { data, error } = await supabase
    .from("current_prices")
    .select("*")
    .gt("unit_count", 1)
    .order("price_per_unit_cents", { ascending: true })
    .limit(limit);

  if (error) {
    console.error("getComparablePrices:", error.message);
    return [];
  }
  return data ?? [];
}

/** Most recent scrape health check per collector. */
export async function getCollectorHealth(limit = 10): Promise<CollectorHealth[]> {
  const supabase = await createClient();
  const { data, error } = await supabase
    .from("collector_health")
    .select("*")
    .order("checked_at", { ascending: false })
    .limit(limit);

  if (error) {
    console.error("getCollectorHealth:", error.message);
    return [];
  }
  return data ?? [];
}

/** What a district's supply list costs at each retailer. */
export async function getSupplyListCosts(): Promise<SupplyListCost[]> {
  const supabase = await createClient();
  const { data, error } = await supabase
    .from("supply_list_costs")
    .select("*")
    .order("total_dollars", { ascending: true });

  if (error) {
    console.error("getSupplyListCosts:", error.message);
    return [];
  }
  return data ?? [];
}

// --- formatting -------------------------------------------------

export function dollars(cents: number): string {
  return (cents / 100).toLocaleString("en-US", {
    style: "currency",
    currency: "USD",
  });
}

/**
 * Per-unit cost. Below a dollar we show cents with two decimals, because
 * the difference between 6.06¢ and 12.50¢ per crayon is the entire point
 * and rounding it to "$0.06 vs $0.13" throws away the comparison.
 */
export function perUnit(cents: number | null): string {
  if (cents === null) return "—";
  if (cents < 100) return `${cents.toFixed(2)}¢`;
  return dollars(cents);
}

export function relativeTime(iso: string): string {
  const then = new Date(iso).getTime();
  const mins = Math.round((Date.now() - then) / 60000);
  if (mins < 1) return "just now";
  if (mins < 60) return `${mins}m ago`;
  const hrs = Math.round(mins / 60);
  if (hrs < 24) return `${hrs}h ago`;
  return `${Math.round(hrs / 24)}d ago`;
}
