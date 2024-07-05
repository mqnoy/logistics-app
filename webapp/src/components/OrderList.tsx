import { FC, ReactNode, useEffect, useState } from 'react'
import { Order, OrderTypeEnum } from '../types'
import { TableCustom } from './TableCustom'
import { ListResponse } from '../types'
import { dateUtils, rtkUtils, toastUtils } from '../utils'
import { Modal } from './Modal'
import { OrderDetail } from './OrderDetail'
import { useLazyGetDetailOrderQuery } from '../api'

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

    const [modalTitle, setModalTitle] = useState('')
    const [isModalActive, setIsModalActive] = useState(false)
    const showModal = () => {
        setIsModalActive(true)
    }

    const closeModal = () => {
        setIsModalActive(false)
    }

    const [
        getDetail,
        { data: dataGetDetail, error: errorGetDetail, isLoading: isloadingGetDetail },
    ] = useLazyGetDetailOrderQuery()
    useEffect(() => {
        if (errorGetDetail) {
            const errorApi = rtkUtils.parseErrorRtk(errorGetDetail)
            toastUtils.fireToastError(errorApi)
        } else if (isloadingGetDetail) {
            console.debug('loading..')
        }
    }, [errorGetDetail, isloadingGetDetail])

    const handleActionDetail = (props: Order) => {
        getDetail(props.id)
        setModalTitle('Detail Order')
        showModal()
    }

    return (
        <>
            <Modal
                title={modalTitle}
                isActive={isModalActive}
                onClose={closeModal}
                content={<OrderDetail data={dataGetDetail?.data} />}
            />
            {data && (
                <TableCustom
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
                        return (
                            <tr key={item.id}>
                                <td>{renderOrderType(item.type.id)}</td>
                                <td>{item.count_item}</td>
                                <td>{renderRequestAt(item.request_at)}</td>
                                <td>
                                    <div className="field is-grouped">
                                        <p className="control">
                                            <button
                                                className="button is-primary is-outlined"
                                                onClick={() => {
                                                    handleActionDetail(item)
                                                }}
                                            >
                                                detail
                                            </button>
                                        </p>
                                    </div>
                                </td>
                            </tr>
                        )
                    }}
                />
            )}
        </>
    )
}
