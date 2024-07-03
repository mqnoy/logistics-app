export type BaseResponse<D> = {
    success: boolean
    message: string
    data: D
}

export interface MetaData {
    page: number
    limit: number
    total_pages: number
    total_items: number
}

export type ListResponse<R> = {
    metadata: MetaData
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

export interface ListRequest {
    page: number
    limit: number
    offset: number
    orders: string
}
