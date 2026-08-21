"use client";

import { useEffect, useRef } from "react";

type YellowBarProps = {
  side: "left" | "right";
};

const EDGE_MARGIN = 20;

export function YellowBar({ side }: YellowBarProps) {
  const barRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    function handleScroll() {
      if (!barRef.current) return;

      const maxScroll = document.documentElement.scrollHeight - window.innerHeight;
      const scrolled = maxScroll > 0 ? window.scrollY / maxScroll : 0;

      // right bar travels top -> bottom as you scroll, left bar goes the opposite way
      const progress = side === "right" ? scrolled : 1 - scrolled;

      // keep the bar's own height inside the viewport so it never scrolls off screen
      const travel = window.innerHeight - barRef.current.offsetHeight - EDGE_MARGIN * 2;
      barRef.current.style.top = `${EDGE_MARGIN + progress * travel}px`;
    }

    handleScroll();
    window.addEventListener("scroll", handleScroll);
    // sizes change across breakpoints, so a resize needs to reposition the bar too
    window.addEventListener("resize", handleScroll);
    return () => {
      window.removeEventListener("scroll", handleScroll);
      window.removeEventListener("resize", handleScroll);
    };
  }, [side]);

  const isRight = side === "right";

  // below lg the green section fills the screen, so a small edge margin already
  // sits on the green part; at lg (1024px) it caps at max-w-5xl and gets centered,
  // so the bar has to chase its edge instead of the viewport edge
  const sidePosition = isRight
    ? "right-2 lg:right-[calc(50vw_-_32rem_+_20px)]"
    : "left-2 lg:left-[calc(50vw_-_32rem_+_20px)]";

  return (
    <div ref={barRef} className={`fixed z-0 pointer-events-none ${sidePosition}`}>
      <div
        className="flex items-center justify-center border-[5px] border-[#FFEC51] bg-[#FFFBDB]
          w-[26px] h-[162px]
          md:w-[35px] md:h-[216px] md:border-[6px]
          lg:w-[44px] lg:h-[270px] lg:border-[8px]"
      >
        <div className="bg-[#000] w-[16px] h-[73px] md:w-[21px] md:h-[98px] lg:w-[26px] lg:h-[122px]"></div>
      </div>
    </div>
  );
}
