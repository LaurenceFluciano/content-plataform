'use client'
import { useState } from "react";
import Button from "../atom/Button";
import { Field } from "../molecule/Field";

export default function AlmostThereForm()
{  
    const labelStyle = 
        `
            text-c-medium md:text-c-large
            text-neutral-primary/80
            group-focus-within:text-brand-primary/90
        `;
    
    const [formData, setFormData] = useState({name: ''});

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const {id, value} = e.target;
        setFormData(prev => ({...prev, [id]: value}));
    };

    const isFormValid = formData.name.length > 0;

    return (
        <>
            <form className="w-[70%] md:w-[60%] xl:w-[40%] flex flex-col gap-8 lg:shadow-500 lg:px-4xl lg:py-3xl lg:rounded-xl">
                <h1 className="text-h3 text-brand-secondary mx-auto mb-2">Almost There!</h1>

                <Field id="name" className="w-full">
                    <Field.Label className={labelStyle} text="Enter your profile name: " />
                    <Field.Input onChange={handleChange} type="text" className="input-medium lg:input-large" variant="default" placeholder="E-mail" />
                </Field>


                <Button disabled={!isFormValid} className="w-full mt-5 btn-medium lg:btn-large" variant="brand" type="submit">
                    {"Finish"}
                </Button>

                <span className="mx-auto text-c-medium lg:text-c-large font-brand text-neutral-primary/70">
                    Step 3 / 3
                </span>
            </form>
        </>
    )
}