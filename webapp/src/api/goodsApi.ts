import { createApi } from '@reduxjs/toolkit/query/react'
import { GoodCreateRequest, GoodCreateResponse, GoodsListResponse, ListRequest } from '../types'
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
        postCreateGood: builder.mutation<GoodCreateResponse, GoodCreateRequest>({
            query: (body) => ({
                url: '/goods',
                method: 'POST',
                body: body,
            }),
        }),
    }),
})

export const { useGetListGoodsQuery, useLazyGetListGoodsQuery, usePostCreateGoodMutation } =
    goodsApi
