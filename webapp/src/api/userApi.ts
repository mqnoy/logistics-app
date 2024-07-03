import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'
import { EnvConfig } from '../config'
import { UserRegisterRequest, UserRegisterResponse } from '../types'

export const userApi = createApi({
    reducerPath: 'userApi',
    baseQuery: fetchBaseQuery({
        baseUrl: EnvConfig.apiBaseURL,
        prepareHeaders: (headers) => {
            headers.set('Access-Control-Allow-Origin', '*')
            return headers
        },
    }),
    endpoints: (builder) => ({
        postUserRegister: builder.mutation<UserRegisterResponse, UserRegisterRequest>({
            query: (body) => ({
                url: '/users/register',
                method: 'POST',
                body: body,
            }),
        }),
    }),
})

export const { usePostUserRegisterMutation } = userApi
