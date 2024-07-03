import { BaseResponse, User } from '.'

export interface UserLoginRequest {
    email: string
    password: string
}

export interface IUserLoginResponse  {
    access_token: string
    refresh_token: string
    user: User
}

export type UserLoginResponse = BaseResponse<IUserLoginResponse>