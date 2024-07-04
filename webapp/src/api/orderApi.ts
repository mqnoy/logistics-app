import { createApi } from '@reduxjs/toolkit/query/react'
import { baseQuery } from '.'
import { OrderCreateRequest, OrderCreateResponse, OrderTypeEnum } from '../types'

export const orderApi = createApi({
    reducerPath: 'orderApi',
    baseQuery: baseQuery,
    endpoints: (builder) => ({
        postOrderCreate: builder.mutation<OrderCreateResponse, OrderCreateRequest>({
            query: (body) => {
                const q = {
                    url: '',
                    method: 'POST',
                    body: body,
                }
                if (body.type === OrderTypeEnum.ORDER_IN) {
                    q.url = '/orders/goods/in'
                } else {
                    q.url = '/orders/goods/out'
                }
                return q
            },
        }),
    }),
})

export const { usePostOrderCreateMutation } = orderApi
