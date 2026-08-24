"use client"

import { useRef, useState } from "react"
import { Upload } from "lucide-react"
import { Input } from "@/components/ui/input"

interface FileUploadProps {
  onFileSelect: (file: File) => void
}

export function FileUpload({ onFileSelect }: FileUploadProps) {
  const inputRef = useRef<HTMLInputElement>(null)
  const [dragging, setDragging] = useState(false)

  const handleFile = (file: File | undefined) => {
    if (!file) return
    onFileSelect(file)
  }

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
        handleFile(e.dataTransfer.files[0])
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
        Drop your grocery.txt here
      </p>

      <p className="mt-1 text-xs text-muted-foreground">
        or click to browse
      </p>

      <Input
        ref={inputRef}
        type="file"
        accept=".txt"
        className="hidden"
        onChange={(e) => {
          handleFile(e.target.files?.[0])
        }}
      />
    </div>
  )
}
