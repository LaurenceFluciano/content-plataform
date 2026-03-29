'use client'
import { useState } from "react";
import Button from "../atom/Button";
import { Field } from "../molecule/Field";

export default function RegisterForm()
{  
    const labelStyle = 
        `
            text-c-medium md:text-c-large
            text-neutral-primary/80
            group-focus-within:text-brand-primary/90
        `;
    
    const [formData, setFormData] = useState({email: '', password: ''});

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const {id, value} = e.target;
        setFormData(prev => ({...prev, [id]: value}));
    };

    const isFormValid = formData.email.length > 0 && formData.password.length > 0;

    return (
        <>
            <form className="w-[70%] md:w-[60%] xl:w-[40%] flex flex-col gap-4 lg:shadow-500 lg:px-4xl lg:py-3xl lg:rounded-xl">
                <h1 className="text-h3 text-brand-secondary mx-auto mb-2">Sign Up</h1>

                <Field id="email" className="w-full">
                    <Field.Label className={labelStyle} text="E-mail: " />
                    <Field.Input onChange={handleChange} type="text" className="input-medium lg:input-large" variant="default" placeholder="E-mail" />
                </Field>

                <Field id="password" className="w-full">
                    <Field.Label className={labelStyle} text="Password: " />
                    <Field.Input onChange={handleChange} type="password" className="input-medium lg:input-large" variant="default" placeholder="Password" />
                </Field>


                <div className="flex flex-col gap-1">
                    <Button disabled={!isFormValid} className="w-full mt-5 btn-medium lg:btn-large" variant="brand" type="submit">
                        {"Register"}
                    </Button>

                    <span className="text-b3 lg:text-b1 font-brand">
                        <span className="text-neutral-primary/60">You already have an account?</span>
                        <a href="/login" className="text-brand-primary cursor-pointer"> Click here to Login.</a>
                    </span>
                </div>

                <span className="mx-auto text-c-medium lg:text-c-large font-brand text-neutral-primary/70">
                    Step 1 / 3
                </span>
            </form>
        </>
    )
}