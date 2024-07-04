import { FC } from 'react'
import Layout from '../Layout'
import { DatePicker, OrderList } from '../../components'
import { BaseResponse, ListResponse } from '../../types'
import { Order, OrderTypeEnum } from '../../types/order'
import { dateUtils } from '../../utils'
import mockOrders from '@assets/mock/orders.json'
import { DatePickerEventSelect } from '../../components/DatePicker/type'

export const OrderPage: FC = () => {
    // TODO: call service orderApi
    const raw = mockOrders as BaseResponse<ListResponse<Order>>
    console.log(raw.data);

    return (
        <Layout>
            <div className="section">
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
                                    console.debug(`onchange: `, e.target.value);
                                }}
                            >
                                <option value={0}>- select type -</option>
                                <option value={OrderTypeEnum.ORDER_IN}>order-in</option>
                                <option value={OrderTypeEnum.ORDER_OUT}>order-out</option>
                            </select>
                        </div>
                    </div>
                    <div className="column column is-10">
                        <p>Request date</p>
                        <div className="select is-primary">
                            <DatePicker
                                isRange={true}
                                onSelected={(event: DatePickerEventSelect) => {
                                    const { startDate, endDate } = event.data
                                    const from = dateUtils.dateToEpoch(startDate)
                                    const to = dateUtils.dateToEpoch(endDate)
                                    console.debug(from);
                                    console.debug(to);
                                }}
                            />
                        </div>
                    </div>
                </div>
                {raw && <OrderList
                    onPageChange={() => {

                    }}
                    data={raw.data}
                />}
            </div>
        </Layout>
    )
}
