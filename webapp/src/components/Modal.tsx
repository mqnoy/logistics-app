import { FC, ReactNode } from "react";

interface ModalProps {
    title: string;
    isActive: boolean;
    onClose: () => void;
    content: ReactNode;
    buttons?: ReactNode[]
}

export const Modal: FC<ModalProps> = ({ title, isActive, onClose, content, buttons }) => {
    return (
        <div>
            <div className={`modal ${isActive ? 'is-active' : ''}`}>
                <div className="modal-background" onClick={onClose}></div>
                <div className="modal-card">
                    <header className="modal-card-head">
                        <p className="modal-card-title">{title}</p>
                        <button className="delete" aria-label="close" onClick={onClose}></button>
                    </header>
                    <section className="modal-card-body">
                        <div className="content">
                            {content}
                        </div>
                    </section>
                    <footer className="modal-card-foot">
                        <div className="buttons">
                            {
                                buttons && buttons.length > 0 &&
                                buttons.map((button) => {
                                    return button
                                })
                            }
                        </div>
                    </footer>
                </div>
            </div >
        </div >
    );
};
