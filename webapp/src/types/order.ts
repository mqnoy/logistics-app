import { BaseResponse, ListResponse, TimestampResponse } from './common'

export interface Order extends TimestampResponse {
    id: string
    request_at: number
    type: OrderType
    total: number
}

export interface OrderType {
    id: number
    name: string
}

export type OrderListResponse = BaseResponse<ListResponse<Order>>

export enum OrderTypeEnum {
    ORDER_IN = 1,
    ORDER_OUT,
}
