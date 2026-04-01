import BrandNavBar from "@/modules/components/molecule/BrandNavBar";
import AlmostThereForm from "@/modules/components/organism/AlmostThereForm";
import ConfirmEmailForm from "@/modules/components/organism/ConfirmEmailForm";

export default function Page() {
    return (
        <>
            <header>
                <BrandNavBar />
            </header>
            <main className="w-full mt-5 flex items-center justify-center bg-surface-primary">
                <ConfirmEmailForm />
            </main>
        </>
    );
}