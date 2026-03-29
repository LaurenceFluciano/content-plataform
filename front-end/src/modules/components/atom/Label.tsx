'use client'
import { LabelHTMLAttributes } from 'react';

interface LabelProps
    extends LabelHTMLAttributes<HTMLLabelElement>
{
    text: string,
}


export function Label({text, ...props}: LabelProps)
{
    return (
        <label className={`transition-all duration-200 font-brand ` + props.className} {...props}>
            {text}
        </label>
    )
}