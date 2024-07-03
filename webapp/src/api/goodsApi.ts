import { createApi } from '@reduxjs/toolkit/query/react'
import { GoodsListResponse, ListRequest } from '../types'
import { baseQuery } from '.'

export const goodsApi = createApi({
    reducerPath: 'goodsApi',
    baseQuery: baseQuery,
    endpoints: (builder) => ({
        getListGoods: builder.query<GoodsListResponse, ListRequest>({
            query: (params) => ({
                url: '/goods',
                params,
            }),
        }),
    }),
})

export const { useGetListGoodsQuery, useLazyGetListGoodsQuery } = goodsApi
