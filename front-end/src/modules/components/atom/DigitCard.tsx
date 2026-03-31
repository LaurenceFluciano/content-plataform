'use client'
import { cn } from "@/modules/lib/utils";
import { Digit } from "@/types/digit";
import { cva, VariantProps } from "class-variance-authority";

export const digitCardVariants = cva(
    "card-digit-base card-digit-mobile",
    {
        variants: {
            variant: {
                default: "card-digit-default" 
            },
        },
        defaultVariants: {
            variant: "default",
        },
    }
);

export interface DigitCardProps
    extends VariantProps<typeof digitCardVariants>      
{
    digit: Digit,
    className?: string,
    onClick?: () => void
}

export default function DigitCard(
    { 
        variant, 
        className, 
        digit, 
        onClick
    }: DigitCardProps & { ref?: React.Ref<HTMLSpanElement> }
){
    return (
        <span 
            onClick={onClick}
            className={cn(
                digitCardVariants({ variant }), 
                !digit && "card-digit-placeholder",
                className
            )}
        >
            {digit || "0"}
        </span>
    );
}