'use client'
import { cn } from '@/modules/lib/utils';
import { cva, type VariantProps } from 'class-variance-authority';
import React from 'react';

const buttonVariants = cva(
  "button-base",
  {
    variants: {
      variant: {
        brand: "btn-primary-brand",
        neutral: "btn-primary-neutral",
        danger: "btn-primary-danger",
      }
    },
    defaultVariants: {
      variant: "brand",
    },
  }
);

interface ButtonProps 
  extends 
    React.ButtonHTMLAttributes<HTMLButtonElement>, 
    VariantProps<typeof buttonVariants> 
{
  children?: React.ReactNode;
}

export default function Button({
    children = "Button", 
    variant, 
    type = "button", 
    className,
    ...props
}: ButtonProps) {
    return (
        <button 
            className={cn(buttonVariants({ variant }), className)} 
            type={type}
            {...props} 
        >
            {children}
        </button>
    )
}