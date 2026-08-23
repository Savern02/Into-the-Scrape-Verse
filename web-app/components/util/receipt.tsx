"use client";

import { useRef } from "react";
import { toPng } from "html-to-image";

interface ReceiptItem {
  name: string;
  price: number;
}

interface ReceiptProps {
  items: ReceiptItem[];
  taxRate?: number;
}

export function Receipt({
  items,
  taxRate = 0.09,
}: ReceiptProps) {
  const receiptRef = useRef<HTMLDivElement>(null);

  const subtotal = items.reduce(
    (total, item) => total + item.price,
    0
  );

  const tax = subtotal * taxRate;
  const total = subtotal + tax;

  const downloadReceipt = async () => {
    if (!receiptRef.current) return;

    const dataUrl = await toPng(receiptRef.current, {
      pixelRatio: 2,
    });

    const link = document.createElement("a");
    link.download = "cheap-chick-sample-receipt.png";
    link.href = dataUrl;
    link.click();
  };

  return (
    <div className="flex flex-col items-center gap-4">
      <div
        ref={receiptRef}
        className="w-80 bg-white text-black p-6 font-mono"
      >
        <div className="text-center">
          <p className="border-t-2 border-dashed border-black" />

          <h2 className="text-xl font-bold mt-3">
            CHEAP CHICK
          </h2>

          <p>SAMPLE RECEIPT</p>

          <p className="border-t-2 border-dashed border-black mt-3" />
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

          <div className="border-t-2 border-dashed border-black my-3" />

          <div className="flex justify-between text-lg font-bold">
            <span>TOTAL</span>
            <span>${total.toFixed(2)}</span>
          </div>
        </div>

        <div className="text-center mt-6">
          <div className="border-t-2 border-dashed border-black mb-3" />

          <p className="font-bold">SAMPLE RECEIPT</p>
          <p>FOR DEMONSTRATION</p>

          <div className="border-t-2 border-dashed border-black mt-3" />
        </div>
      </div>

      <button
        onClick={downloadReceipt}
        className="rounded-lg bg-seaweed-600 px-6 py-2 font-bold text-white hover:scale-105 transition"
      >
        Download Receipt
      </button>
    </div>
  );
}
