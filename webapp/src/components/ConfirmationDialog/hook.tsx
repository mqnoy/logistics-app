import { ReactNode, useState } from "react";
import ConfirmationDialog from ".";
import { ConfirmationDialogProps } from "./type";

export const useConfirmationDialog = () => {
    const [isDialogOpen, setDialogOpen] = useState(false);
    const [dialogContent, setDialogContent] = useState<ReactNode>(null);
    const [onConfirmCallback, setOnConfirmCallback] = useState<() => void>(() => () => { });
    const [onCloseCallback, setOnCloseCallback] = useState<() => void>(() => () => { });

    const showDialog = (props: ConfirmationDialogProps) => {
        setDialogContent(props.content);
        setOnConfirmCallback(() => props.onConfirm);
        setOnCloseCallback(() => props.onClose || (() => setDialogOpen(false)));
        setDialogOpen(true);
    };

    const hideDialog = () => setDialogOpen(false);

    const ConfirmationDialogComponent = (
        <ConfirmationDialog
            isActive={isDialogOpen}
            content={dialogContent}
            onConfirm={() => {
                onConfirmCallback();
                hideDialog();
            }}
            onClose={() => {
                onCloseCallback();
                hideDialog();
            }}
        />
    );

    return { showDialog, ConfirmationDialogComponent };
};
