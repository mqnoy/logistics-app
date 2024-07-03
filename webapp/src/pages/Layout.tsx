import { FC, ReactNode } from "react";
import { Navbar } from "../components";

interface LayoutProps {
    children: ReactNode;
}

const Layout: FC<LayoutProps> = ({ children }) => {
    return (
        <>
            <Navbar />
            <section className="hero is-fullheight">
                <div className="container">
                    {children}
                </div>
            </section>
        </>

    );
};

export default Layout;