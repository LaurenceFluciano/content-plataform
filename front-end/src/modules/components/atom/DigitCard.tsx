'use client'
import { cn } from "@/modules/lib/utils";
import { cva, VariantProps } from "class-variance-authority";
import React, { forwardRef, useRef, useState } from "react";

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
    className?: string,
    onDigitChange?: (digit: string) => void
}

const DigitCard = forwardRef<HTMLSpanElement, DigitCardProps>(
    ({ variant, className, onDigitChange }, ref) => {
        const [value, setValue] = useState("");

        const handleKeyDown = (e: React.KeyboardEvent<HTMLSpanElement>) => {
            if (/^\d$/.test(e.key)) {
                e.preventDefault();
                setValue(e.key);
                if (onDigitChange) onDigitChange(e.key);
                return;
            }

            if (e.key === "Backspace" || e.key === "Delete") {
                e.preventDefault();
                const wasEmpty = value === "";
                setValue("");
                if (onDigitChange) onDigitChange(wasEmpty ? "backspace-empty" : "");
                return;
            }

            if (e.key.length === 1) e.preventDefault();
        };

        return (
            <span 
                ref={ref}
                contentEditable={true}
                onKeyDown={handleKeyDown}
                suppressContentEditableWarning={true}
                className={cn(
                    digitCardVariants({ variant }), 
                    !value && "text-neutral-primary/40",
                    className
                )}
            >
                {value || "0"}
            </span>
        );
    }
);

export default DigitCard;