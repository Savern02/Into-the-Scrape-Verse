"use client"

import { useRef } from "react"
import { toPng } from "html-to-image"

interface ReceiptItem {
  name: string
  price: number
}

interface ReceiptProps {
  retailer: string
  zipcode: string
  items: ReceiptItem[]
  taxRate?: number
}

export function Receipt({
  retailer,
  zipcode,
  items,
  taxRate = 0.09,
}: ReceiptProps) {
  const receiptRef = useRef<HTMLDivElement>(null)

  const subtotal = items.reduce(
    (total, item) => total + item.price,
    0
  )

  const tax = subtotal * taxRate
  const total = subtotal + tax

  const downloadReceipt = async () => {
    if (!receiptRef.current) return

    const dataUrl = await toPng(receiptRef.current, {
      pixelRatio: 2,
    })

    const link = document.createElement("a")
    link.download = `cheap-chick-${retailer}-receipt.png`
    link.href = dataUrl
    link.click()
  }

  return (
    <div className="flex flex-col items-center gap-4">
      <div
        ref={receiptRef}
        className="w-80 bg-white p-6 font-mono text-black"
      >
        <div className="text-center">
          <p className="border-t-2 border-dashed border-black" />

          <h2 className="mt-3 text-xl font-bold">
            CHEAP CHICK
          </h2>

          <p>{retailer}</p>
          <p>ZIP: {zipcode}</p>

          <p className="mt-3 border-t-2 border-dashed border-black" />
        </div>

        <div className="my-5 space-y-2">
          {items.map((item, index) => (
            <div
              key={index}
              className="flex justify-between gap-4"
            >
              <span>{item.name}</span>
              <span>${item.price.toFixed(2)}</span>
            </div>
          ))}
        </div>

        <div className="border-t-2 border-dashed border-black pt-3">
          <div className="flex justify-between">
            <span>SUBTOTAL</span>
            <span>${subtotal.toFixed(2)}</span>
          </div>

          <div className="flex justify-between">
            <span>TAX</span>
            <span>${tax.toFixed(2)}</span>
          </div>

          <div className="my-3 border-t-2 border-dashed border-black" />

          <div className="flex justify-between text-lg font-bold">
            <span>TOTAL</span>
            <span>${total.toFixed(2)}</span>
          </div>
        </div>

        <div className="mt-6 text-center">
          <div className="mb-3 border-t-2 border-dashed border-black" />

          <p className="font-bold">CHEAP CHICK</p>
          <p>PRICE COMPARISON</p>

          <div className="mt-3 border-t-2 border-dashed border-black" />
        </div>
      </div>

      <button
        onClick={downloadReceipt}
        className="rounded-lg bg-seaweed-600 px-6 py-2 font-bold text-white transition hover:scale-105"
      >
        Download Receipt
      </button>
    </div>
  )
}
