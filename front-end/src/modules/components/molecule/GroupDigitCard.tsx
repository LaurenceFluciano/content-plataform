'use client'
import React, { useRef, useState } from "react";
import DigitCard from "../atom/DigitCard";
import { Digit } from "@/types/digit";

interface GroupDigitCardProps {
    ammountDigitCards: number;
    className?: string;
}

export default function GroupDigitCard({ ammountDigitCards, className }: GroupDigitCardProps) {
    const [digits, setDigits] = useState<Digit[]>(Array(ammountDigitCards).fill(''));
    const digitCardRefs = useRef<(HTMLSpanElement | null)[]>([]);
    const inputRef = useRef<HTMLInputElement | null>(null);

    const isValidDigit = (char: string) => {
        return /^\d$/.test(char);
    }

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        e.preventDefault();
        const fullValue = e.target.value; 

        const newDigits = Array(ammountDigitCards).fill('') as Digit[];

        fullValue.split('').forEach((char, i) => {
            if (i < ammountDigitCards && isValidDigit(char)) {
                newDigits[i] = char as Digit;
            }
        });

        setDigits(newDigits);
        
    }

    return (
        <>
        
        <input
            ref={inputRef}
            onChange={handleChange}
            type="text"
            value={digits.join('') || ''}
            autoComplete="off"
            maxLength={ammountDigitCards}
            className="opacity-0 absolute"
            name="otp-code"
        />

        <div className="flex gap-2 justify-center ">
            {Array.from({ length: ammountDigitCards }).map((_, i) => (
                <DigitCard
                    key={i}
                    onClick={() => { 
                        inputRef.current?.focus();
                        
                        inputRef.current?.setSelectionRange(i+1, i+1);
                    }}
                    ref={(el) => { digitCardRefs.current[i] = el; }} 
                    className={className}
                    digit={digits[i]}
                />
            ))}
        </div>
        </>
    );
}