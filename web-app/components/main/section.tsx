"use client";

import { ReactNode, useEffect, useRef, useState } from "react";

interface FeatureSectionProps {
  id?: string;
  title: string;
  description: string;
  children?: ReactNode;
  className?: string;
}

export function Section({
  id,
  title,
  description,
  children,
  className = "",
}: FeatureSectionProps) {
  const sectionRef = useRef<HTMLElement>(null);
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        setIsVisible(entry.isIntersecting);
      },
      {
        rootMargin: "-40% 0px -40% 0px",
        threshold: 0,
      }
    );

    if (sectionRef.current) {
      observer.observe(sectionRef.current);
    }

    return () => observer.disconnect();
  }, []);

  return (
    <section ref={sectionRef} id={id} className={className}>
      <div
        className={`flex flex-col justify-center items-center bg-seaweed-500 border-4 transition-all duration-300 ${
          isVisible
            ? "border-light-yellow-500"
            : "border-transparent"
        }`}
      >
        <h1
          className={`text-center text-5xl leading-none font-bold underline decoration-[4px] underline-offset-[14px] transition-all duration-300 ${
            isVisible
              ? "decoration-light-yellow-500"
              : "decoration-transparent"
          }`}
        >
          {title}
        </h1>

        <p className="text-center my-5 rounded-full  px-6 py-2">
          {description}
        </p>

        {children}
      </div>
    </section>
  );
}
