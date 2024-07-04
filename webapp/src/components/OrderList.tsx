import { FC, ReactNode } from 'react'
import { Order, OrderTypeEnum } from '../types'
import { TableCustom } from './TableCustom'
import { ListResponse } from '../types'
import { dateUtils } from '../utils'

type OrderListProps = {
    data?: ListResponse<Order>
    onPageChange: (page: number) => void
}

export const OrderList: FC<OrderListProps> = ({ data, onPageChange }) => {

    const renderRequestAt = (d: number): string => {
        return dateUtils.epochNumberToDateTimeStr(d)
    }

    const renderOrderType = (t: OrderTypeEnum): ReactNode => {
        if (t === OrderTypeEnum.ORDER_IN) {
            return <span className="button is-rounded is-success is-small has-text-white">In</span>
        }

        return <span className="button is-rounded is-danger is-small has-text-white">Out</span>
    }
    return (
        <>
            {data && <TableCustom
                onPageChange={onPageChange}
                data={data}
                tableHead={
                    <>
                        <th>type</th>
                        <th>total</th>
                        <th>request at</th>
                        <th>Action</th>
                    </>
                }
                renderRow={(item: Order) => {
                    return < tr key={item.id} >
                        <td>{renderOrderType(item.type.id)}</td>
                        <td>{item.total}</td>
                        <td>{renderRequestAt(item.request_at)}</td>
                        <td>
                            <div className="field is-grouped">
                                <p className="control">
                                    <button
                                        className="button is-primary is-outlined"
                                        onClick={() => {
                                            // TODO: handle action detail
                                        }} >
                                        view
                                    </button>
                                </p>
                            </div>
                        </td>
                    </tr>
                }}
            />
            }
        </>
    )
}
