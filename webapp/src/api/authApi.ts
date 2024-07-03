import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'
import { UserLoginRequest, UserLoginResponse } from '../types/auth'
import { EnvConfig } from '../config'

export const authApi = createApi({
    reducerPath: 'authApi',
    baseQuery: fetchBaseQuery({
        baseUrl: EnvConfig.apiBaseURL,
        prepareHeaders: (headers) => {
            headers.set('Access-Control-Allow-Origin', '*')
            return headers
        },
    }),
    endpoints: (builder) => ({
        postUserLogin: builder.mutation<UserLoginResponse, UserLoginRequest>({
            query: (body) => ({
                url: '/users/login',
                method: 'POST',
                body: body,
            }),
        }),
    }),
})

export const { usePostUserLoginMutation } = authApi
