'use client'

import { createContext, InputHTMLAttributes, ReactNode, useContext } from "react";
import { Input, inputVariants } from "../atom/Input";
import { Label } from "../atom/Label";
import { cn } from "@/modules/lib/utils";
import { VariantProps } from "class-variance-authority";


interface FieldContextProps 
{
    id: string;
}

const FieldContext = createContext<FieldContextProps | undefined>(undefined);

interface FieldProps extends FieldContextProps {
    children: ReactNode;
    className?: string;
}

function FieldRoot({ children, id, className }: FieldProps) {
    return (
        <FieldContext.Provider value={{ id }}>
            <div className={cn("group flex flex-col gap-2", className)}>
                {children}
            </div>
        </FieldContext.Provider>
    );
}

export interface LabelProps 
    extends React.LabelHTMLAttributes<HTMLLabelElement> {
    text: string;
}

function FieldLabel({ text, ...props }: LabelProps) {
    const context = useContext(FieldContext);
    if (!context) throw new Error("Field.Label deve estar dentro de <Field />");
    
    return <Label
                htmlFor={context.id}
                text={text} 
                {...props}
            />;
};

interface InputProps
    extends InputHTMLAttributes<HTMLInputElement>,
    VariantProps<typeof inputVariants> 
{ }

function FieldInput({ ...props }: InputProps) {
    const context = useContext(FieldContext);
    if (!context) throw new Error("Field.Input deve estar dentro de <Field />");
    
    return (
        <Input
            id={context.id}
            {...props}
        />
    );
};

export const Field = Object.assign(FieldRoot, {
    Label: FieldLabel,
    Input: FieldInput,
});