import Link from "next/link";
import { getComparablePrices, perUnit } from "@/lib/queries";
import { Button } from "@/components/ui/button";

/**
 * The hero states the thesis and then proves it with live numbers from the
 * database. A claim followed immediately by its own evidence is more
 * persuasive than either alone — and it means the landing page breaks
 * visibly if the pipeline stops running, which is the honest behaviour.
 */
export async function Hero() {
  const rows = await getComparablePrices(60);

  const priced = rows.filter((r) => r.price_per_unit_cents !== null);
  const cheapest = priced.at(0);
  const dearest = priced.at(-1);

  const multiple =
    cheapest?.price_per_unit_cents && dearest?.price_per_unit_cents
      ? dearest.price_per_unit_cents / cheapest.price_per_unit_cents
      : null;

  return (
    <section className="flex flex-col gap-8 py-12">
      <div className="flex flex-col gap-4 max-w-2xl">
        <p className="text-xs uppercase tracking-widest text-muted-foreground">
          Classroom supply pricing
        </p>
        <h1 className="text-4xl sm:text-5xl font-semibold tracking-tight text-balance leading-[1.1]">
          The cheapest box is rarely the cheapest crayon.
        </h1>
        <p className="text-muted-foreground text-pretty">
          Supply catalogues price by the package. Classrooms consume by the
          item. We scrape education suppliers, normalize every pack size, and
          show what one unit actually costs — so a set of 800 and a set of 24
          can be compared honestly.
        </p>
      </div>

      {multiple && cheapest && dearest ? (
        <div className="flex flex-col gap-4">
          <div className="flex flex-wrap items-end gap-x-10 gap-y-6">
            <Stat
              label="Cheapest per unit"
              value={perUnit(cheapest.price_per_unit_cents)}
              detail={cheapest.name}
            />
            <Stat
              label="Priciest per unit"
              value={perUnit(dearest.price_per_unit_cents)}
              detail={dearest.name}
            />
            <Stat
              label="Spread"
              value={`${multiple.toFixed(1)}×`}
              detail={`across ${priced.length} tracked products`}
            />
          </div>

          <p className="text-sm text-muted-foreground max-w-xl text-pretty">
            Same product category, same retailer, {multiple.toFixed(1)} times
            the cost per item. None of that is visible from a product grid.
          </p>
        </div>
      ) : (
        <p className="text-sm text-muted-foreground">
          No priced products yet — run the ingest pipeline to populate this.
        </p>
      )}

      <div>
        <Button asChild>
          <Link href="/supplies">See the full comparison</Link>
        </Button>
      </div>
    </section>
  );
}

function Stat({
  label,
  value,
  detail,
}: {
  label: string;
  value: string;
  detail: string;
}) {
  return (
    <div className="flex flex-col gap-1 min-w-0">
      <span className="text-xs uppercase tracking-wide text-muted-foreground">
        {label}
      </span>
      <span className="text-3xl font-semibold tabular-nums">{value}</span>
      <span className="text-xs text-muted-foreground truncate max-w-[16rem]">
        {detail}
      </span>
    </div>
  );
}
