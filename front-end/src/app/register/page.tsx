import BrandNavBar from "@/modules/components/molecule/BrandNavBar";
import RegisterForm from "@/modules/components/organism/RegisterForm";

export default function Page() {
    return (
        <>
            <header>
                <BrandNavBar />
            </header>
            <main className="w-full mt-5 flex items-center justify-center bg-surface-primary">
                <RegisterForm />
            </main>
        </>
    );
}