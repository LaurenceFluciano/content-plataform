'use client'
import Image from 'next/image'

interface LogoProps {
    className?: string
}


export default function Logo({className}: LogoProps)
{
    return (
        <div className={className}>
            <Image src='/Logo.png' width={64} height={72} alt='Brand Logo' priority />
        </div>
    );
}