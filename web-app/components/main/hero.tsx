import Image from "next/image";
import Link from "next/link";
import { Button } from "@/components/ui/button"
import { ProductCard } from "@/components/g-ui/item"
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
    <CarouselItem>
                  <ProductCard image="https://i5.walmartimages.com/seo/Red-Bull-Amber-Edition-Sugar-Free-Energy-Drink-Strawberry-Apricot-114mg-Caffeine-12-fl-oz-Pack-of-4-Cans_89dee0d8-212d-4e07-bbdd-d9d707003190.83ecf6741e5ae0f8a5affc87bb508771.jpeg"
  name="Red Bull Amber Edition Sugarfree Energy Drink, Strawberry Apricot"
  itemUrl="https://www.walmart.com/ip/Red-Bull-Amber-Edition-Sugar-Free-Energy-Drink-Strawberry-Apricot-12-fl-oz-Pack-of-4-Cans/5332753715"
  price="$10.75"
  website="Walmart"/>
                </CarouselItem>
    <CarouselItem>BEANS2</CarouselItem>
    <CarouselItem>BEANS3</CarouselItem>
  </CarouselContent>
  <CarouselPrevious />
  <CarouselNext />
</Carousel></div>
<div className ="grid col-start-4"><Carousel>
  <CarouselContent>
    <CarouselItem> <ProductCard image="https://i5.walmartimages.com/seo/Red-Bull-Sugar-Free-Energy-Drink-4-pk-Cans-12-oz_75b01883-f18b-45cc-9735-6d56c5a60717.1242b8f1b72564bf2d5aa60c42e8c90a.jpeg"
                name="Red Bull Sugar Free Energy Drink"
                itemUrl="https://www.walmart.com/ip/Red-Bull-Sugar-Free-Energy-Drink-4-pk-Cans-12-oz/8639967110?classType=REGULAR Price:19.75 Currency:USD Specs:[{Name:Brand Value:Red Bull} {Name:Flavor Value:Sugar Free} {Name:Size Value:12 oz} {Name:Container type Value:Can} {Name:Texture Value:lightly sparkling} {Name:Food form Value:Liquids} {Name:Food & drug fact label type Value:Nutrition Fact Panel}] Images:[https://i5.walmartimages.com/seo/Red-Bull-Sugar-Free-Energy-Drink-4-pk-Cans-12-oz_75b01883-f18b-45cc-9735-6d56c5a60717.1242b8f1b72564bf2d5aa60c42e8c90a.jpeg https://i5.walmartimages.com/asr/b580bf46-1883-4cb0-b7b8-03486e0d8d7e.e9660ebab2bf280042ecadf1ac5ed60b.jpeg https://i5.walmartimages.com/asr/56400165-82bf-4a70-8f31-bc58b05e2971.7d8a840d9527400b513e3d003959befa.jpeg https://i5.walmartimages.com/asr/dc76e1b4-13aa-48d2-842c-d934b497dbd7.6d278601b08567697b53df620310b862.jpeg https://i5.walmartimages.com/asr/549db7dc-54fa-42ca-b487-ce78fe8f915b.07198fb8e3e03888b7e27f622e5f03d2.jpeg"
  price="$19.75"
  website="Walmart"/>
</CarouselItem>
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
