import { BaseResponse, ListResponse, TimestampResponse } from './common'
import { GoodSnapshot } from '.'

export interface Order extends TimestampResponse {
    id: string
    request_at: number
    type: OrderType
    total: number
    good_snapshot?: GoodSnapshot
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

interface GoodOrderCreateRequest {
    code: string
}

export interface OrderCreateRequest {
    good: GoodOrderCreateRequest
    total: number
    type?: OrderTypeEnum
}

export type OrderCreateResponse = BaseResponse<Order>

export interface MultipleOrderCreateRequest {
    items: OrderItemRequest[]
    type?: OrderTypeEnum
}

export interface OrderItemRequest {
    code: string
    total: number
}

export interface OrderItemReason {
    code: string
    total: number
    reason: string
}

export interface MultipleOrderCreateResponse {
    id: string
    success: OrderItemReason[]
    failed: OrderItemReason[]
}
