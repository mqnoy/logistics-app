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

export interface GoodCreateRequest {
    code: string
    name: string
    description: string
}

export type GoodCreateResponse = BaseResponse<Goods>

export type ModalActionGoods = 'create' | 'update' | 'detail'

export interface GoodUpdateRequest extends GoodCreateRequest {
    id: string
    is_active: boolean
}

export type GoodUpdateResponse = BaseResponse<Goods>
