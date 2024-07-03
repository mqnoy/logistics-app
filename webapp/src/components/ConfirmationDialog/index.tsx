import { FC } from "react";
import { ConfirmationDialogProps } from "./type";

const ConfirmationDialog: FC<ConfirmationDialogProps> = ({ isActive, content, onClose, onConfirm }) => {
    if (!isActive) return null;

    return (
        <div className={`modal ${isActive ? 'is-active' : ''}`}>
            <div className="modal-background"></div>
            <div className="modal-card">
                <header className="modal-card-head">
                    <p className="modal-card-title">Confirm Action</p>
                    <button className="delete" aria-label="close" onClick={onClose}></button>
                </header>
                <section className="modal-card-body">
                    {content}
                </section>
                <footer className="modal-card-foot">
                    <div className="buttons">
                        <button className="button is-success has-text-white" onClick={onConfirm}>Yes</button>
                        <button className="button" onClick={onClose}>No</button>
                    </div>
                </footer>
            </div>
        </div>
    );
};

export default ConfirmationDialog
