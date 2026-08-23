import Link from "next/link"

import { Button } from "@/components/ui/button"
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "@/components/ui/carousel"

import { ProductCard } from "@/components/g-ui/item"

export default function DashboardPage() {
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
            <CarouselItem className="basis-full sm:basis-1/2 lg:basis-1/3">
              <ProductCard
                image="/products/redbull.jpg"
                name="Red Bull Sugar Free Energy Drink"
                itemUrl="/dashboard/products/redbull"
                price="$7.98"
                website="Walmart"
              />
            </CarouselItem>

            <CarouselItem className="basis-full sm:basis-1/2 lg:basis-1/3">
              <ProductCard
                image="/products/beans.jpg"
                name="Great Value Black Beans"
                itemUrl="/dashboard/products/beans"
                price="$1.24"
                website="Walmart"
              />
            </CarouselItem>

            <CarouselItem className="basis-full sm:basis-1/2 lg:basis-1/3">
              <ProductCard
                image="/products/milk.jpg"
                name="Great Value Vitamin D Whole Milk"
                itemUrl="/dashboard/products/milk"
                price="$3.48"
                website="Walmart"
              />
            </CarouselItem>
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

        <Carousel className="w-full">
          <CarouselContent>
            <CarouselItem className="basis-full sm:basis-1/2 lg:basis-1/3">
              <ProductCard
                image="/products/redbull.jpg"
                name="Red Bull Sugar Free Energy Drink"
                itemUrl="/dashboard/products/redbull"
                price="$7.98"
                website="Walmart"
              />
            </CarouselItem>

            <CarouselItem className="basis-full sm:basis-1/2 lg:basis-1/3">
              <ProductCard
                image="/products/redbull.jpg"
                name="Red Bull Sugar Free Energy Drink"
                itemUrl="/dashboard/products/redbull"
                price="$8.98"
                website="Target"
              />
            </CarouselItem>

            <CarouselItem className="basis-full sm:basis-1/2 lg:basis-1/3">
              <ProductCard
                image="/products/redbull.jpg"
                name="Red Bull Sugar Free Energy Drink"
                itemUrl="/dashboard/products/redbull"
                price="$9.48"
                website="Kroger"
              />
            </CarouselItem>
          </CarouselContent>

          <CarouselPrevious />
          <CarouselNext />
        </Carousel>
      </section>
    </div>
  )
}
