import { FC, useEffect, useState } from 'react'
import Layout from '../Layout'
import { DatePicker, OrderList } from '../../components'
import { OrderTypeEnum } from '../../types/order'
import { dateUtils, rtkUtils, toastUtils } from '../../utils'
import { DatePickerEventSelect } from '../../components/DatePicker/type'
import { NavLink } from 'react-router-dom'
import { useLazyGetListOrdersQuery } from '../../api'

export const OrderPage: FC = () => {
    const [page, setPage] = useState(1);
    const [query, setQuery] = useState({
        limit: 10,
        offset: 0,
        page: page,
        orders: 'id desc',
    })
    const [requestDateQuery, setRequestDateQuery] = useState({
        from: 0,
        to: 0
    })
    const [orderType, setOrderType] = useState(0)

    const [getListOrder, { data, error, isLoading }] = useLazyGetListOrdersQuery();
    useEffect(() => {
        if (error) {
            const errorApi = rtkUtils.parseErrorRtk(error);
            toastUtils.fireToastError(errorApi)

        } else if (isLoading) {
            console.log('loading...');
        }
    }, [error, isLoading])

    const handleOnPageChange = (page: number) => {
        setPage(page)
    }

    useEffect(() => {
        const dateRange = [requestDateQuery.from, requestDateQuery.to]
        setQuery({
            ...query,
            ...{
                request_at_range: `[${dateRange.join(',')}]`
            }
        })
    }, [requestDateQuery])

    useEffect(() => {
        setQuery({
            ...query,
            ...{
                type: orderType
            }
        })
    }, [orderType])


    useEffect(() => {
        setQuery({
            ...query,
            page: page
        })
    }, [page])

    useEffect(() => {
        getListOrder(query);
    }, [query])


    useEffect(() => {
        getListOrder(query);
    }, [])


    return (
        <Layout>
            <div className="py-5">
                <div className="columns">
                    <div className="column">
                        <span className="icon-text has-text-info">
                            <h5 className="title">Orders</h5>
                        </span>

                    </div>
                </div>
                <div className="columns is-vcentered">
                    <div className="column is-2">
                        <p>Order type</p>
                        <div className="select is-primary">
                            <select
                                onChange={(e) => {
                                    setOrderType(Number(e.target.value))
                                }}
                            >
                                <option value={0}>- select type -</option>
                                <option value={OrderTypeEnum.ORDER_IN}>order-in</option>
                                <option value={OrderTypeEnum.ORDER_OUT}>order-out</option>
                            </select>
                        </div>
                    </div>
                    <div className="column column is-7">
                        <p>Request date</p>
                        <div className="select is-primary">
                            <DatePicker
                                isRange={true}
                                onSelected={(event: DatePickerEventSelect) => {
                                    const { startDate, endDate } = event.data
                                    const from = dateUtils.dateToEpoch(startDate)
                                    const to = dateUtils.dateToEpoch(endDate)
                                    setRequestDateQuery({
                                        from,
                                        to
                                    })
                                }}
                            />
                        </div>
                    </div>
                    <div className="column is-3 is-flex is-justify-content-flex-end">
                        <p></p>
                        <NavLink className="button is-primary has-text-white" to={'/orders/create'}>
                            Create Order
                        </NavLink>
                    </div>
                </div>
                {data?.data && <OrderList
                    onPageChange={handleOnPageChange}
                    data={data?.data}
                />}
            </div>
        </Layout>
    )
}
