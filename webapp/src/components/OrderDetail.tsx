import { FC } from "react"
import { Order } from "../types"
import { dateUtils } from "../utils"

type OrderDetailProps = {
    data?: Order
}

export const OrderDetail: FC<OrderDetailProps> = ({ data }) => {
    if (!data) {
        return <>no data</>
    }

    const renderRequestAt = (d: number): string => {
        return dateUtils.epochNumberToDateTimeStr(d)
    }

    return (
        <div className="">
            <div className="fixed-grid has-1-cols box">
                <div className="grid">
                    <div className="cell">{data.type.name}</div>
                    <div className="cell">Request at: {renderRequestAt(data.request_at)}</div>
                    <div className="cell">Total: {data.total}</div>
                </div>
            </div>
            {
                data.good_snapshot &&
                <div className="box content">
                    <h5>Goods</h5>
                    <div className="field">
                        <label htmlFor="" className="label">
                            Code
                        </label>
                        <div className="control">
                            <div
                                className="input"
                            >{data.good_snapshot.code}</div>
                        </div>
                    </div>
                    <div className="field">
                        <label htmlFor="" className="label">
                            Name
                        </label>
                        <div className="control">
                            <div
                                className="input"
                            >{data.good_snapshot.name}</div>
                        </div>
                    </div>
                    <div className="field">
                        <label htmlFor="" className="label">
                            Description
                        </label>
                        <div className="control">
                            <div
                                className="input"
                            >{data.good_snapshot.description}</div>
                        </div>
                    </div>
                    <div className="field">
                        <label htmlFor="" className="label">
                            Active
                        </label>
                        <div className="control">
                            <div
                                className="input"
                            >{data.good_snapshot.is_active ? 'active' : 'inactive'}</div>
                        </div>
                    </div>
                </div>
            }
        </div>

    )
}
