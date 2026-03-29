import "../styles/global.css"
export default function MainLayout({
    children,
}: {
    children: React.ReactNode
}) {
    return (
        <html lang="pt-br">
            <body>
                {children}
            </body>
        </html>
    )
}