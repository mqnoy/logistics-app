import { ReactNode } from 'react'

export type ConfirmationDialogProps = {
    isActive?: boolean
    content?: ReactNode
    onClose?: () => void
    onConfirm?: () => void
}
