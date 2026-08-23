"use client"

import { useRef, useState } from "react"
import { Upload } from "lucide-react"
import { Input } from "@/components/ui/input"

export function FileUpload() {
  const inputRef = useRef<HTMLInputElement>(null)
  const [dragging, setDragging] = useState(false)

  return (
    <div
      onDragOver={(e) => {
        e.preventDefault()
        setDragging(true)
      }}
      onDragLeave={() => setDragging(false)}
      onDrop={(e) => {
        e.preventDefault()
        setDragging(false)

        const files = e.dataTransfer.files
        console.log(files)
      }}
      onClick={() => inputRef.current?.click()}
      className={`flex cursor-pointer flex-col items-center justify-center rounded-lg border-2 border-dashed p-10 text-center transition-colors ${
        dragging
          ? "border-primary bg-muted"
          : "border-muted-foreground/25 hover:bg-muted/50"
      }`}
    >
      <Upload className="mb-3 h-8 w-8" />

      <p className="text-sm font-medium">
        Drop your file here
      </p>

      <p className="mt-1 text-xs text-muted-foreground">
        or click to browse
      </p>

      <Input
        ref={inputRef}
        type="file"
        className="hidden"
        onChange={(e) => {
          console.log(e.target.files)
        }}
      />
    </div>
  )
}
