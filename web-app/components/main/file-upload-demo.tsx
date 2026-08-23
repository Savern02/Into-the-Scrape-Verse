"use client"

import { useEffect, useState } from "react"
import { FileText } from "lucide-react"

export function FileUploadDemo() {
  const [dragging, setDragging] = useState(false)
  const [uploaded, setUploaded] = useState(false)

  useEffect(() => {
    let timeout: ReturnType<typeof setTimeout>

    const animate = () => {
      setDragging(false)
      setUploaded(false)

      // Wait before starting
      timeout = setTimeout(() => {
        setDragging(true)

        // File reaches upload area
        timeout = setTimeout(() => {
          setDragging(false)
          setUploaded(true)

          // Hold uploaded state
          timeout = setTimeout(() => {
            setUploaded(false)

            // Wait before restarting
            timeout = setTimeout(animate, 1000)
          }, 1800)
        }, 1400)
      }, 1200)
    }

    animate()

    return () => clearTimeout(timeout)
  }, [])

  return (
    <div className="relative flex h-64 w-full items-center justify-center">
      {/* Fake file */}
      <div
        className={`absolute left-1/2 z-10 -translate-x-1/2 transition-all duration-[1400ms] ease-in-out ${
          dragging
            ? "top-1/2 -translate-y-1/2 scale-75 opacity-0"
            : "top-4 scale-100 opacity-100"
        }`}
      >
        <div className="flex flex-col items-center">
          <FileText className="h-10 w-10" />

          <span className="mt-2 whitespace-nowrap text-xs font-medium">
            grocery-list.txt
          </span>
        </div>
      </div>

      {/* Fake upload box */}
      <div
        className={`flex h-44 w-full max-w-md flex-col items-center justify-center rounded-lg border-2 border-dashed text-center transition-all duration-300 ${
          dragging
            ? "scale-[1.03] border-primary bg-muted"
            : "border-muted-foreground/25"
        }`}
      >
        {uploaded ? (
          <>
            <FileText className="mb-3 h-8 w-8" />

            <p className="text-sm font-medium">
              File uploaded!
            </p>

            <p className="mt-1 text-xs text-muted-foreground">
              grocery-list.txt
            </p>
          </>
        ) : (
          <>
            <p className="text-sm font-medium">
              Drop your file here
            </p>

            <p className="mt-1 text-xs text-muted-foreground">
              or click to browse
            </p>
          </>
        )}
      </div>
    </div>
  )
}
