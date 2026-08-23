"use client";

import { useEffect, useState } from "react";
import { Input } from "@/components/ui/input";

export function ZipDemo() {
  const zipCodes = ["74103", "74135", "74105", "74133"];

  const [value, setValue] = useState("");
  const [zipIndex, setZipIndex] = useState(0);

  useEffect(() => {
    const zip = zipCodes[zipIndex];
    let currentIndex = 0;

    const typeInterval = setInterval(() => {
      if (currentIndex < zip.length) {
        setValue(zip.slice(0, currentIndex + 1));
        currentIndex++;
      } else {
        clearInterval(typeInterval);

        setTimeout(() => {
          setValue("");

          setTimeout(() => {
            setZipIndex((prev) => (prev + 1) % zipCodes.length);
          }, 500);
        }, 1500);
      }
    }, 250);

    return () => clearInterval(typeInterval);
  }, [zipIndex]);

  return (
    <div className="flex flex-col items-center gap-4 text-white">
      <Input
        value={value}
        readOnly
        placeholder="Enter your ZIP code"
        className="w-48 text-center text-lg"
      />
    </div>
  );
}
