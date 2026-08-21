import Image from "next/image";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "@/components/ui/carousel"

export function Hero() {
  return (
<div className="flex flex-col gap-16 items-center">
  <section id="title">
  <div className="relative">
    <Image
      src="/logos/logo.png"
      alt="Cheap Chick"
      width={100}
      height={50}
      className="absolute right-full top-1/4 -translate-y-1/2 mr-4"
    />

    <p className="text-center text-5xl leading-none font-bold">
      Cheap Chick <br />
    </p>
    <p className="text-sm leading-none text-center">
        Compare groceries from the biggest retailers
    </p>
  </div>
  </section>
  <section id="Compare Items Feature">
        <h1 className="text-center text-5xl leading-none font-bold underline decoration-[9px] decoration-light-yellow-500 underline-offset-[10px]">Compare Items! <br/></h1>
        <p className="text-center my-5">Description</p>
<div className="flex flex-row w-full gap-20 my-10"> 
<div className="flex w-1/2">
  <Carousel>
  <CarouselContent>
    <CarouselItem>BEANS1</CarouselItem>
    <CarouselItem>BEANS2</CarouselItem>
    <CarouselItem>BEANS3</CarouselItem>
  </CarouselContent>
  <CarouselPrevious />
  <CarouselNext />
</Carousel></div>
<div className ="flex w-1/2"><Carousel>
  <CarouselContent>
    <CarouselItem>TARGET BEANS</CarouselItem>
    <CarouselItem>TARGET BENAS</CarouselItem>
    <CarouselItem>TRAGET BEANS</CarouselItem>
  </CarouselContent>
  <CarouselPrevious />
  <CarouselNext />
</Carousel></div>



</div> 

  </section>
  <section id="Grocery List Feature">
  <div className="flex gap-8 justify-center items-center">
        <h1 className="text-center text-5xl leading-none font-bold underline decoration-[9px] decoration-light-yellow-500 underline-offset-[10px]">Upload Grocery Lists!</h1>
  </div>
  </section>
  <section id="Recipets Feature">
  <div className="flex gap-8 justify-center items-center">
        <h1 className="text-center text-5xl leading-none font-bold underline decoration-[9px] decoration-light-yellow-500 underline-offset-[10px]">Get Recipets!</h1>
  </div>
  </section>
  <section>
  <div className="flex gap-8 justify-center items-center">
        <h1 className="text-center text-5xl leading-none font-bold underline decoration-[9px] decoration-light-yellow-500 underline-offset-[10px]">All Local!</h1>
  </div>
  </section>
</div>
  );
}
