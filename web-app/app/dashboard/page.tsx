import Link from "next/link"
import { createClient } from "@/lib/supabase/server"

import { Button } from "@/components/ui/button"
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "@/components/ui/carousel"

import { ProductCard } from "@/components/g-ui/item"

export default async function DashboardPage() {
  const supabase = await createClient()

  const { data: products, error } = await supabase
    .from("product_data")
    .select("*")
    .order("created_at", { ascending: false })
    .limit(10)

  if (error) {
    console.error("Failed to fetch products:", error)
  }

  return (
    <div className="space-y-8">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold">Dashboard</h1>
        <p className="text-muted-foreground">
          Compare grocery prices from different stores.
        </p>
      </div>

      {/* Products */}
      <section className="space-y-4">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-2xl font-semibold">Products</h2>
            <p className="text-sm text-muted-foreground">
              Recently added products
            </p>
          </div>

          <Button asChild variant="outline">
            <Link href="/dashboard/products">
              View All
            </Link>
          </Button>
        </div>

        <Carousel className="w-full">
          <CarouselContent>
            {products?.map((product) => (
              <CarouselItem
                key={product.id}
                className="basis-full sm:basis-1/2 lg:basis-1/3"
              >
                <ProductCard
                  image={product.image_urls[0]}
                  name={product.brand}
                  itemUrl={product.url}
                  price={`${product.currency}${product.final_price}`}
                  website={product.retailer}
                />
              </CarouselItem>
            ))}
          </CarouselContent>

          <CarouselPrevious />
          <CarouselNext />
        </Carousel>
      </section>

      {/* Comparisons */}
      <section className="space-y-4">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-2xl font-semibold">Price Comparisons</h2>
            <p className="text-sm text-muted-foreground">
              Find the best prices across stores
            </p>
          </div>

          <Button asChild variant="outline">
            <Link href="/dashboard/comparisons">
              View All
            </Link>
          </Button>
        </div>

        {/* Comparison products will go here */}
      </section>
    </div>
  )
}
