"use client"

import { useState } from "react"

import { FileUpload } from "@/components/g-ui/file-upload"
import { Receipt } from "@/components/util/receipt"

import { createClient } from "@/lib/supabase/client"

interface Product {
  id: string
  final_price: number
  currency: string
  specifications: unknown[]
  image_urls: string[]
  brand: string
  retailer: string
  is_available: boolean
  zipcode: string
}

interface ReceiptItem {
  name: string
  price: number
}

interface RetailerReceipt {
  retailer: string
  items: ReceiptItem[]
}

export default function GroceryListPage() {
  const [receipts, setReceipts] = useState<RetailerReceipt[]>([])
  const [loading, setLoading] = useState(false)
  const [message, setMessage] = useState("")

  async function handleFileSelect(file: File) {
    setLoading(true)
    setMessage("")
    setReceipts([])

    try {
      // Read grocery.txt
      const text = await file.text()

      const groceryItems = text
        .split(/\r?\n/)
        .map((item) => item.trim().toLowerCase())
        .filter(Boolean)

      if (groceryItems.length === 0) {
        setMessage("Your grocery list is empty.")
        return
      }

      const supabase = createClient()

      // Get the signed-in user
      const {
        data: { user },
        error: userError,
      } = await supabase.auth.getUser()

      if (userError || !user) {
        setMessage("You must be signed in.")
        return
      }

      // Get the user's ZIP code
      const { data: profile, error: profileError } = await supabase
        .from("profiles")
        .select("zip_code")
        .eq("id", user.id)
        .single()

      if (profileError || !profile?.zip_code) {
        setMessage("Please add your ZIP code in Settings first.")
        return
      }

      // Get products available in the user's ZIP code
      const { data: products, error: productsError } = await supabase
        .from("product_data")
        .select("*")
        .eq("zipcode", profile.zip_code.toString())
        .eq("is_available", true)

      if (productsError) {
        console.error(productsError)
        setMessage("Failed to load products.")
        return
      }

      if (!products || products.length === 0) {
        setMessage("No products were found for your ZIP code.")
        return
      }

      /*
       * Group products by retailer.
       *
       * For each grocery item:
       *   1. Find matching products
       *   2. Group them by retailer
       *   3. Keep the cheapest product for each retailer
       */
      const retailerMap = new Map<string, ReceiptItem[]>()

      for (const groceryItem of groceryItems) {
        const matches = products.filter((product) =>
          productMatches(product, groceryItem)
        )

        // Find the cheapest matching product for EACH retailer
        const cheapestByRetailer = new Map<string, Product>()

        for (const product of matches) {
          const existing = cheapestByRetailer.get(product.retailer)

          if (
            !existing ||
            Number(product.final_price) <
              Number(existing.final_price)
          ) {
            cheapestByRetailer.set(product.retailer, product)
          }
        }

        // Add the cheapest product to that retailer's receipt
        for (const [retailer, product] of cheapestByRetailer) {
          if (!retailerMap.has(retailer)) {
            retailerMap.set(retailer, [])
          }

          retailerMap.get(retailer)!.push({
            name: getProductName(product),
            price: Number(product.final_price),
          })
        }
      }

      const generatedReceipts: RetailerReceipt[] = Array.from(
        retailerMap.entries()
      ).map(([retailer, items]) => ({
        retailer,
        items,
      }))

      if (generatedReceipts.length === 0) {
        setMessage(
          "None of the grocery items could be found for your ZIP code."
        )
        return
      }

      setReceipts(generatedReceipts)
    } catch (error) {
      console.error(error)
      setMessage("Something went wrong processing your grocery list.")
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-3xl font-bold">
          Grocery List
        </h1>

        <p className="text-muted-foreground">
          Upload your grocery list to compare prices.
        </p>
      </div>

      <FileUpload onFileSelect={handleFileSelect} />

      {loading && (
        <p className="text-sm text-muted-foreground">
          Finding the best prices...
        </p>
      )}

      {message && (
        <p className="text-sm text-muted-foreground">
          {message}
        </p>
      )}

      {receipts.length > 0 && (
        <div className="flex flex-wrap gap-8">
          {receipts.map((receipt) => (
            <Receipt
              key={receipt.retailer}
              retailer={receipt.retailer}
              zipcode=""
              items={receipt.items}
            />
          ))}
        </div>
      )}
    </div>
  )
}
