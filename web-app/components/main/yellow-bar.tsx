"use client";

import { useEffect, useRef } from "react";

type YellowBarProps = {
  side: "left" | "right";
};

const EDGE_MARGIN = 20;
const NAV_HEIGHT = 64;

export function YellowBar({ side }: YellowBarProps) {
  const barRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    function handleScroll() {
      if (!barRef.current) return;

      const maxScroll =
        document.documentElement.scrollHeight - window.innerHeight;

      const scrolled =
        maxScroll > 0 ? window.scrollY / maxScroll : 0;

      const progress =
        side === "right" ? scrolled : 1 - scrolled;

      const barHeight = barRef.current.offsetHeight;

      const minTop = NAV_HEIGHT + EDGE_MARGIN;
      const maxTop =
        window.innerHeight - barHeight - EDGE_MARGIN;

      const top =
        minTop + progress * (maxTop - minTop);

      barRef.current.style.top = `${top}px`;
    }

    handleScroll();

    window.addEventListener("scroll", handleScroll);
    window.addEventListener("resize", handleScroll);

    return () => {
      window.removeEventListener("scroll", handleScroll);
      window.removeEventListener("resize", handleScroll);
    };
  }, [side]);

  const sidePosition =
    side === "right"
      ? "right-2 lg:right-0"
      : "left-2 lg:left-0";

  return (
    <div
      ref={barRef}
      className={`fixed z-10 pointer-events-none ${sidePosition}`}
    >
      <div
        className="
          flex items-center justify-center
          border-[5px] border-[#FFEC51]
          bg-[#FFFBDB]
          w-[26px] h-[162px]

          md:w-[35px] md:h-[216px]
          md:border-[6px]

          lg:w-[44px] lg:h-[270px]
          lg:border-[8px]
        "
      >
        <div
          className="
            bg-black
            w-[16px] h-[73px]
            md:w-[21px] md:h-[98px]
            lg:w-[26px] lg:h-[122px]
          "
        />
      </div>
    </div>
  );
}
