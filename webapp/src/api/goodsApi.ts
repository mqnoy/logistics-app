import { createApi } from '@reduxjs/toolkit/query/react'
import {
    BaseResponse,
    GoodCreateRequest,
    GoodCreateResponse,
    GoodUpdateRequest,
    GoodUpdateResponse,
    Goods,
    GoodsListResponse,
    ListRequest,
} from '../types'
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
        getDetailGood: builder.query<BaseResponse<Goods>, string>({
            query: (id) => ({
                url: `/goods/${id}`,
            }),
        }),
        putUpdateGood: builder.mutation<GoodUpdateResponse, GoodUpdateRequest>({
            query: (body) => ({
                url: `/goods/${body.id}`,
                method: 'PUT',
                body: body,
            }),
        }),
        deleteGood: builder.mutation<BaseResponse<null>, string>({
            query: (id) => ({
                url: `/goods/${id}`,
                method: 'DELETE',
            }),
        }),
    }),
})

export const {
    useGetListGoodsQuery,
    useLazyGetListGoodsQuery,
    usePostCreateGoodMutation,
    useLazyGetDetailGoodQuery,
    usePutUpdateGoodMutation,
    useDeleteGoodMutation,
} = goodsApi
