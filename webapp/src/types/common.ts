export type BaseResponse<D> = {
    success: boolean
    message: string
    data: D
}

export type ListResponse<R> = {
    metadata: unknown
    rows: R[]
}

export interface TimestampResponse {
    created_at: number
    updated_at: number
}

export interface ErrorValidator {
    code: string
    field: string
    message: string
}

export interface ErrorApiResponse extends BaseResponse<null> {
    errors?: ErrorValidator[]
}
