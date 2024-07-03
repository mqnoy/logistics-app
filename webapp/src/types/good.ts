import { BaseResponse, ListResponse, TimestampResponse } from './common'

export interface Goods extends TimestampResponse {
    id: string
    code: string
    name: string
    description: string
    is_active: boolean
    stock: Stock
}

export interface Stock {
    total: number
}

export type GoodsListResponse = BaseResponse<ListResponse<Goods>>
