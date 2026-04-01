'use client'
import { useState } from "react";
import Button from "../atom/Button";
import { Field } from "../molecule/Field";
import { createClient } from "@/modules/lib/supabase/client";
import { redirect } from "next/navigation";

export default function LoginForm()
{
    const labelStyle = 
    `
        text-c-medium md:text-c-large
        text-neutral-primary/80
        group-focus-within:text-brand-primary/90
    `;

    const [formData, setFormData] = useState({email: '', password: ''});

    const supabase = createClient();

    const handleLogin = async (data: FormData) => {
        const email = data.get('email') as string;
        const password = data.get('password') as string;

        const { error } = await supabase.auth.signInWithPassword({email, password});

        if ( error )
        {
            console.error("Não foi possivel efetuar o login")
        }
        else
            redirect('/feed')
    }

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const {id, value} = e.target;
        setFormData(prev => ({...prev, [id]: value}));
    };

    const isFormValid = formData.email.length > 0 && formData.password.length > 0;

    return (
        <>
            <form action={handleLogin} className="w-[70%] md:w-[60%] xl:w-[40%] flex flex-col gap-4 lg:shadow-500 lg:px-4xl lg:py-3xl lg:rounded-xl">
                <h1 className="text-h3 text-brand-secondary mx-auto mb-5">Login</h1>

                <Field id="email" className="w-full">
                    <Field.Label htmlFor="email" className={labelStyle} text="E-mail: " />
                    <Field.Input name="email" onChange={handleChange} type="text" className="input-medium lg:input-large" variant="default" placeholder="E-mail" />
                </Field>

                <Field id="password" className="w-full">
                    <Field.Label htmlFor="password" className={labelStyle} text="Password: " />
                    <Field.Input name="password" onChange={handleChange} type="password" className="input-medium lg:input-large" variant="default" placeholder="Password" />
                </Field>


                <div className="flex flex-col gap-1">
                    <Button disabled={!isFormValid} className="w-full mt-5 btn-medium lg:btn-large" variant="brand" type="submit">
                        {"Login"}
                    </Button>

                    <span className="text-b3 lg:text-b1 font-brand">
                        <span className="text-neutral-primary/60">You don't have an account?</span>
                        <a href="/register" className="text-brand-primary cursor-pointer"> Click here to Register.</a>
                    </span>
                </div>
            </form>
        </>
    )
}