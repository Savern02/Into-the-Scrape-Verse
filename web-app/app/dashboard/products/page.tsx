"use client"

import { useEffect, useMemo, useState } from "react"
import { Search } from "lucide-react"

import { createClient } from "@/lib/supabase/client"

import { Input } from "@/components/ui/input"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

import { ProductCard } from "@/components/g-ui/item"

interface Product {
  id: string
  url: string
  final_price: number
  currency: string
  specifications: unknown[]
  image_urls: string[]
  brand: string
  retailer: string
  is_available: boolean
  zipcode: string
}

export default function ProductsPage() {
  const [products, setProducts] = useState<Product[]>([])
  const [search, setSearch] = useState("")
  const [retailer, setRetailer] = useState("all")
  const [sort, setSort] = useState("lowest")
  const [loading, setLoading] = useState(true)
  const [message, setMessage] = useState("")

  useEffect(() => {
    async function loadProducts() {
      const supabase = createClient()

      try {
        const {
          data: { user },
        } = await supabase.auth.getUser()

        if (!user) {
          setMessage("You must be signed in.")
          return
        }

        const { data: profile, error: profileError } =
          await supabase
            .from("profiles")
            .select("zip_code")
            .eq("id", user.id)
            .single()

        if (profileError || !profile?.zip_code) {
          setMessage(
            "Please add your ZIP code in Settings first."
          )
          return
        }

        const { data, error } = await supabase
          .from("product_data")
          .select("*")
          .eq("zipcode", profile.zip_code.toString())
          .eq("is_available", true)

        if (error) {
          console.error(error)
          setMessage("Failed to load products.")
          return
        }

        setProducts(data ?? [])
      } catch (error) {
        console.error(error)
        setMessage("Something went wrong.")
      } finally {
        setLoading(false)
      }
    }

    loadProducts()
  }, [])

  const retailers = useMemo(() => {
    return Array.from(
      new Set(products.map((product) => product.retailer))
    ).sort()
  }, [products])

  const filteredProducts = useMemo(() => {
    const searchTerm = search.toLowerCase().trim()

    const filtered = products.filter((product) => {
      const matchesRetailer =
        retailer === "all" ||
        product.retailer === retailer

      if (!matchesRetailer) {
        return false
      }

      if (!searchTerm) {
        return true
      }

      const specifications = JSON.stringify(
        product.specifications
      ).toLowerCase()

      return (
        product.brand.toLowerCase().includes(searchTerm) ||
        specifications.includes(searchTerm)
      )
    })

    return [...filtered].sort((a, b) => {
      if (sort === "highest") {
        return Number(b.final_price) - Number(a.final_price)
      }

      return Number(a.final_price) - Number(b.final_price)
    })
  }, [products, search, retailer, sort])

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-3xl font-bold">Products</h1>

        <p className="text-muted-foreground">
          Search and compare grocery products near you.
        </p>
      </div>

      {/* Filters */}
      <div className="flex flex-col gap-3 sm:flex-row">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />

          <Input
            value={search}
            onChange={(event) =>
              setSearch(event.target.value)
            }
            placeholder="Search products..."
            className="pl-9"
          />
        </div>

        <Select
          value={retailer}
          onValueChange={setRetailer}
        >
          <SelectTrigger className="w-full sm:w-48">
            <SelectValue placeholder="Retailer" />
          </SelectTrigger>

          <SelectContent>
            <SelectItem value="all">
              All retailers
            </SelectItem>

            {retailers.map((retailer) => (
              <SelectItem
                key={retailer}
                value={retailer}
              >
                {retailer}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>

        <Select value={sort} onValueChange={setSort}>
          <SelectTrigger className="w-full sm:w-48">
            <SelectValue placeholder="Sort by" />
          </SelectTrigger>

          <SelectContent>
            <SelectItem value="lowest">
              Lowest price
            </SelectItem>

            <SelectItem value="highest">
              Highest price
            </SelectItem>
          </SelectContent>
        </Select>
      </div>

      {/* Results */}
      {loading ? (
        <p className="text-muted-foreground">
          Loading products...
        </p>
      ) : message ? (
        <p className="text-muted-foreground">
          {message}
        </p>
      ) : filteredProducts.length === 0 ? (
        <div className="py-12 text-center">
          <p className="font-medium">
            No products found.
          </p>

          <p className="mt-1 text-sm text-muted-foreground">
            Try changing your search or filters.
          </p>
        </div>
      ) : (
        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {filteredProducts.map((product) => (
            <ProductCard
              key={product.id}
              image={product.image_urls[0]}
              name={product.brand}
              itemUrl={product.url}
              price={`${product.currency}${Number(
                product.final_price
              ).toFixed(2)}`}
              website={product.retailer}
            />
          ))}
        </div>
      )}
    </div>
  )
}
