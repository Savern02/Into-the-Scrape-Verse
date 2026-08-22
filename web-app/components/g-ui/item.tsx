import Image from "next/image"
import { Card, CardContent } from "@/components/ui/card"
import Link from "next/link"

interface ProductCardProps {
  image: string
  name: string
  price: string
  website: string
}

export function ProductCard({
  image,
  name,
  itemUrl,
  price,
  website,
}: ProductCardProps) {
  return (
    <Card className="w-full max-w-xs overflow-hidden h-[450px]">
      <div className="aspect-square relative">
        <Link href={itemUrl}>
        <Image
          src={image}
          alt={name}
          width={300}
          height={300}
          className="object-contain"
        />
        </Link>
      </div>

      <CardContent className="space-y-1 p-4">
        <p className="text-lg font-semibold">{price}</p>

        <p className="font-medium">
          {name}
        </p>

        <p className="text-sm text-muted-foreground">
          {website}
        </p>
      </CardContent>
    </Card>
  )
}
