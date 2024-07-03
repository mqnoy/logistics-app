import { ToastContent, ToastOptions, toast } from 'react-toastify'
import { ErrorApiResponse } from '../types'

export const fireToastError = (errorApi: ErrorApiResponse | undefined, options?: ToastOptions) => {
    if (errorApi?.errors) {
        // error validator
        for (const err of errorApi.errors) {
            toast.error(err.message, options)
        }
    } else {
        toast.error(errorApi?.message, options)
    }
}

export const fireToastSuccess = (content: ToastContent, options?: ToastOptions) => {
    toast.success(content, options)
}
