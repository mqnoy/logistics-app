import { BaseResponse, TimestampResponse } from './common'

export interface User extends TimestampResponse {
    id: string
    full_name: string
    email: string
}

export interface UserRegisterRequest {
    full_name: string
    email: string
    password: string
}

export type UserRegisterResponse = BaseResponse<User>
