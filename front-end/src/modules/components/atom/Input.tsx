'use client'
import { cn } from '@/modules/lib/utils';
import { cva, type VariantProps } from 'class-variance-authority';
import { InputHTMLAttributes } from 'react';

export const inputVariants = cva(
    "input-base input-large",
    {
        variants: {
            variant: {
                default: "input-primary-default" 
            },
        },
        defaultVariants: {
            variant: "default",
        },
    }
);


interface InputProps
    extends InputHTMLAttributes<HTMLInputElement>,
    VariantProps<typeof inputVariants> 
{ }


export function Input({className="", variant="default", ...props}: InputProps)
{
    return (
        <input 
            {...props}
            className={cn(inputVariants({ variant }), className)}
        >
        </input>
    )
}