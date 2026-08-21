import Image from "next/image";
import Link from "next/link";
import { Button } from "@/components/ui/button"
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
  <div className="h-px w-full bg-white -my-10"></div>
  <section id="Compare Items Feature">
        <h1 className="text-center text-5xl leading-none font-bold underline decoration-[9px] decoration-light-yellow-500 underline-offset-[10px]">Compare Items! <br/></h1>
        <p className="text-center my-5 rounded-full bg-seaweed-600 px-6 py-2">Compare grocery prices across different stores to find the best deals and save money on your everyday groceries. <br /> Easily see where each item costs less and make smarter choices for your budget.</p>
<div className="grid grid-cols-5 items-center gap-5"> 
<div className="grid col-start-2">
  <Carousel>
  <CarouselContent>
    <CarouselItem>BEANS1</CarouselItem>
    <CarouselItem>BEANS2</CarouselItem>
    <CarouselItem>BEANS3</CarouselItem>
  </CarouselContent>
  <CarouselPrevious />
  <CarouselNext />
</Carousel></div>
<div className ="grid col-start-4"><Carousel>
  <CarouselContent className="flex justify-center">
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
  <div className="flex flex-col justify-center items-center">
        <h1 className="text-center text-5xl leading-none font-bold underline decoration-[9px] decoration-light-yellow-500 underline-offset-[10px]">Upload Grocery Lists!</h1>
   <p className="text-center my-5 rounded-full bg-seaweed-600 px-6 py-2">Upload your grocery list to quickly compare prices across different stores. Save time, find better deals, and make sure you're getting the most for your money on every shopping trip.</p>

        </div>
  </section>
  <section id="Recipets Feature">
  <div className="flex flex-col justify-center items-center">
        <h1 className="text-center text-5xl leading-none font-bold underline decoration-[9px] decoration-light-yellow-500 underline-offset-[10px]">Get Recipets!</h1>
   <p className="text-center my-5 rounded-full bg-seaweed-600 px-6 py-2">Generate a sample receipt for your grocery list and compare the total cost across different stores. See how much your entire shopping trip would cost and find where you can save the most.</p>

        </div>
  </section>
  <section>
  <div className="flex flex-col justify-center items-center">
        <h1 className="text-center text-5xl leading-none font-bold underline decoration-[9px] decoration-light-yellow-500 underline-offset-[10px]">All Local!</h1>
           <p className="text-center my-5 rounded-full bg-seaweed-600 px-6 py-2"># All Local! Enter your ZIP code to see grocery prices from stores in your area. Get local pricing for your shopping list and easily compare nearby stores to find the best deals and save money.
</p>

  </div>
  <div className="flex flex-col justify-center items-center my-20">
<Button asChild size="lg" variant={"default"} className="ease-in-out hover:scale-110">
        <Link className="font-bold" href="/auth/sign-up">Sign up</Link>
      </Button>
        </div>
  </section>
</div>
  );
}
