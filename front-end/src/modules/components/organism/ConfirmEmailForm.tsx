'use client'
import { useRef, useState } from "react";
import Button from "../atom/Button";
import GroupDigitCard from "../molecule/GroupDigitCard";

export default function ConfirmEmailForm()
{   
    const [code, setCode] = useState('');
    const inputRef = useRef<HTMLInputElement | null>(null);

    return (
        <>
            <form className="w-[70%] md:w-[60%] xl:w-[40%] flex flex-col gap-8 lg:shadow-500 lg:px-4xl lg:py-3xl lg:rounded-xl">
                <h1 className="text-h3 text-brand-secondary mx-auto mb-2">Verify Email</h1>
                
                <GroupDigitCard 
                    ammountDigitCards={5}
                    className="card-digit-mobile lg:card-digit-desktop"
                />

                <Button className="w-full mt-5 btn-medium lg:btn-large" variant="brand" type="submit">
                    {"Verify"}
                </Button>

                <span className="mx-auto text-c-medium lg:text-c-large font-brand text-neutral-primary/70">
                    Step 2 / 3
                </span>
            </form>
        </>
    )
}