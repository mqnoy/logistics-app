import { SerializedError } from '@reduxjs/toolkit'
import { FetchBaseQueryError } from '@reduxjs/toolkit/query'
import { ErrorApiResponse } from '../types'

export const parseErrorRtk = (
    error: FetchBaseQueryError | SerializedError
): ErrorApiResponse | undefined => {
    if ('status' in error) {
        const { data } = error
        return data as ErrorApiResponse
    }
}
