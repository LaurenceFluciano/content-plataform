'use client'
import { useRef } from "react";
import DigitCard from "../atom/DigitCard";


interface GroupDigitCardProps {
    ammountDigitCards: number;
    className: string;
}

export default function GroupDigitCard({ className, ammountDigitCards }: GroupDigitCardProps) {
    const digitCardRefs = useRef<(HTMLSpanElement | null)[]>([]);

    const handleNavigation = (index: number, signal: string) => {
        if (signal.length === 1 && index < ammountDigitCards - 1) {
            digitCardRefs.current[index + 1]?.focus();
        } 
        else if (signal === "backspace-empty" && index > 0) {
            digitCardRefs.current[index - 1]?.focus();
        }
    };

    return (
        <div className="flex gap-2 justify-center">
           {Array.from({ length: ammountDigitCards }).map((_, i) => (
                <DigitCard className={className}
                    key={i}
                    ref={(el) => { digitCardRefs.current[i] = el; }}
                    onDigitChange={(val) => handleNavigation(i, val)}
                />
            ))}
        </div>
    );
}