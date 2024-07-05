import { FC, useEffect, useState } from 'react'
import { useLazyGetListOrdersQuery } from '../../api'
import { OrderList } from '../../components'
import { rtkUtils, toastUtils } from '../../utils'
import { useParams } from 'react-router-dom'
import Layout from '../Layout'

export const GoodHistoryOrder: FC = () => {
    const { id } = useParams()

    const [page, setPage] = useState(1)
    const [query, setQuery] = useState({
        limit: 10,
        offset: 0,
        page: page,
        orders: 'id desc',
        goodId: id,
    })

    const [getListOrder, { data, error, isLoading }] = useLazyGetListOrdersQuery()
    useEffect(() => {
        if (error) {
            const errorApi = rtkUtils.parseErrorRtk(error)
            toastUtils.fireToastError(errorApi)
        } else if (isLoading) {
            console.log('loading...')
        }
    }, [error, isLoading])

    const handleOnPageChange = (page: number) => {
        setPage(page)
    }

    useEffect(() => {
        setQuery({
            ...query,
            page: page,
        })
    }, [page])

    useEffect(() => {
        getListOrder(query)
    }, [])

    return (
        <Layout>
            <div className="py-5">
                <div className="columns">
                    <div className="column">
                        <span className="icon-text has-text-info">
                            <h5 className="title">History order</h5>
                        </span>
                    </div>
                </div>
                <div>
                    {data?.data && (
                        <OrderList onPageChange={handleOnPageChange} data={data?.data} />
                    )}
                </div>
            </div>
        </Layout>
    )
}
